package queuebiz

import (
	"webapi/common/queue"
	"webapi/config"
)

var (
	cfgtools      *config.Config
	businessQueue queue.Queue
	cBusiness BusinessInterface
)

func init() {
	queue.Register(&Task{})
}

func InitConfig(conf *config.Config) error {
	cfgtools = conf
	businessQueue = queue.NewRedisQueue(cfgtools)
	return nil
}

// 4) bizutils - 需要执行则初始化,如果仅仅是放入执行队列则不需要
func InitBusiness(bizutils ...interface{}) {
	if len(bizutils) == 0 {
		return
	}

	cBusiness = bizutils[0].(BusinessInterface)
}
