package timer

import (
	"time"
	"adwetec.com/tools/utils"
)

type cron struct {
	nt *nTimer
}

var _nt cron

func init() {
	_nt.init()
}

func (t *cron) init() {
	t.nt = NewCron()
}

func (t *cron) start() {
	t.nt.Start()
}

func (t *cron) addOnce(key string, sec time.Duration, fn func()) {
	t.nt.AddOnce(key, fn, sec)
}

func (t *cron) addForever(key string, sec time.Duration, fn func()) {
	t.nt.AddForever(key, fn, sec)
}

func (t *cron) stop(key string) {
	t.nt.Stop(key)
}

func Start() {
	defer utils.CatchError("timer start")
	_nt.start()
}

// 一次性计时器
func AddOnce(key string, sec time.Duration, fn func()) {
	_nt.addOnce(key, sec, fn)
}

// 永久任务
func AddForever(key string, sec time.Duration, fn func()) {
	_nt.addForever(key, sec, fn)
}

// 计时器停止
func StopCron(key string) {
	_nt.stop(key)
}
