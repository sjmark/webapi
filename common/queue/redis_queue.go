package queue

import (
	"time"
	"sync/atomic"
	"webapi/common/tools/tool"
	"webapi/config"
)

const (
	PullQueueBlockTime   = 5
	PullQueueErrTime     = 2
	PullQueueBkSleepTime = 10
	BKExecDelayTime      = 1800
)

// 基于redis的队列实现
type RedisQueue struct {
	cfg         *config.Config
	queueName   string
	queueNameBK string
	running     int32
	count       uint8
}

var _ = Queue((*RedisQueue)(nil))

func NewRedisQueue(rc *config.Config) *RedisQueue { return &RedisQueue{cfg: rc, queueName: rc.QueueName, queueNameBK: rc.QueueName + "_bk"} }

func (q *RedisQueue) EnQueue(model *QueueModel) error { return q.cfg.RedisR1.Lpush(q.queueName, model) }

func (q *RedisQueue) IsRunning() bool { return atomic.LoadInt32(&q.running) != 0 }

func (q *RedisQueue) Quit() { atomic.StoreInt32(&q.running, 0) }

func (q *RedisQueue) DeQueueAll(ch chan bool) {

	defer tool.PrintPanicStack("DeQueueAll")

	defer atomic.CompareAndSwapInt32(&q.running, 1, 0)

	if !atomic.CompareAndSwapInt32(&q.running, 0, 1) {
		return
	}

	go q.DeQueueBK()

	for q.IsRunning() {
		model := &QueueModel{}
		err := q.cfg.RedisR1.BRpopLpush(q.queueName, q.queueNameBK, PullQueueBlockTime, model)

		if err != nil {
			//q.cfg.Logrus.Errorf("RedisQueue %s BRpopLpush error: %v", q.queueName, err)
			time.Sleep(PullQueueErrTime * time.Second)
		} else {

			if model.GUID != "" {
				// 每秒出列多少个
				//if q.count > q.cfg.Config.Number {
				//	<-time.Tick(q.cfg.WaitTimer)
				//	q.count = 0
				//}
				//q.count++
				if err := model.Task.Exec(); err == nil {
					q.cfg.RedisR1.Lpop(q.queueNameBK)
				} else {
					//q.cfg.Logrus.Errorf("model.Task.Exec error: %v", q.queueName, err)
				}
			} else {
				// 这里做一个预处理
				go func(ch chan bool) {
					<-time.Tick(time.Second * 10)
					if ok, _ := q.cfg.RedisR1.Exists(q.queueName); !ok {
						ch <- true
					} else {
						go q.DeQueueAll(ch)
					}
				}(ch)
			}
		}
	}
}

/*
延时执行才行 不然容易DeQueueAll执行BRpopLpush，立马DeQueueBK执行Rpop，造成重复执行model
延时BKExecDelayTime=1800s
*/
func (q RedisQueue) DeQueueBK() {
	defer tool.PrintPanicStack("DeQueueBK")
	for q.IsRunning() {
		model := &QueueModel{}
		err := q.cfg.RedisR1.Rpop(q.queueNameBK, model)
		if err != nil {
			//q.cfg.Logrus.Errorf("DeQueueBK Rpop error:%v", err)
		} else {
			modelCreateTime := time.Unix(model.Createdtime, 0)
			durationSec := int(time.Now().Sub(modelCreateTime).Seconds())
			if durationSec < BKExecDelayTime {
				time.Sleep(time.Duration(BKExecDelayTime-durationSec) * time.Second)
			}

			if model.GUID != "" {

				if model.RetryCount >= model.AllowRetryCount {
					continue
				}

				// 每秒出列多少个
				//if q.count > q.cfg.Config.Number {
				//	<-time.Tick(q.cfg.WaitTimer)
				//	q.count = 0
				//}
				//
				//q.count++

				if err := model.Task.Exec(); err != nil {
					model.RetryCount = model.RetryCount + 1
					if model.RetryCount < model.AllowRetryCount {
						q.cfg.RedisR1.Lpush(q.queueNameBK, model)
					}
				}
			}
		}
		time.Sleep(PullQueueBkSleepTime * time.Second)
	}
}
