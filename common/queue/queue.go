package queue

import (
	"time"
	"encoding/gob"

	"github.com/google/uuid"
)

func init() {
	Register(&QueueModel{})
}

// 注册对象用于序列化和反序列化
func Register(obj interface{}) {
	gob.Register(obj)
}

// 队列任务接口
type Task interface {
	Exec() error
}

// 队列状态
type QueueStatus int

const (
	QueueStatusInit        QueueStatus = iota
	QueueStatusExecSuccess
	QueueStatusExecFail
)

// 队列任务模型
type QueueModel struct {
	GUID            string
	AllowRetryCount int
	RetryCount      int
	Status          QueueStatus
	Createdtime     int64
	Extra           string
	Task            Task
}

// 队列接口
type Queue interface {
	EnQueue(model *QueueModel) error
	DeQueueAll(ch chan bool)
	Quit()
	IsRunning() bool
}

// 入队任务queuebiz
func EnQueueTask(q Queue, t Task) (string, error) {

	//新建队列任务模型,会分配guid
	guid, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}
	var m = &QueueModel{
		GUID:            guid.String(),
		AllowRetryCount: 3,
		Status:          QueueStatusInit,
		Createdtime:     time.Now().Unix(),
	}

	m.Task = t
	err = q.EnQueue(m)
	return m.GUID, err
}
