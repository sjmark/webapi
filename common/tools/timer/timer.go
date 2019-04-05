package timer

import (
	"sort"
	"time"
	"adwetec.com/tools/utils"
)

type nTimer struct {
	timers   []*timer
	stop     chan string
	location *time.Location
	close    chan *time.Timer
	begin    chan struct{}
}

type timer struct {
	cType   clockType // 定时器类型 1 为一次性 2 为永久
	running bool
	next    time.Time
	oType   string
	sec     time.Duration
	fn      func()
}

type clockType uint8

const (
	clockOnce    clockType = 1 // 一次性定时器 只执行一次
	clockForever clockType = 2 // 循环执行定时器
)

func NewCron() *nTimer {
	return &nTimer{
		timers:   nil,
		stop:     make(chan string),
		close:    make(chan *time.Timer),
		location: time.Now().Location(),
		begin:    make(chan struct{}),
	}
}

func (c *nTimer) AddOnce(ot string, fn func(), sec time.Duration) { c.add(ot, clockOnce, fn, sec) }

func (c *nTimer) AddForever(ot string, fn func(), sec time.Duration) { c.add(ot, clockForever, fn, sec) }

func (c *nTimer) Start() { c.run() }

func (c *nTimer) Stop(ot string) {

	for index, v := range c.timers {

		if index >= len(c.timers) {
			break
		}

		if v.oType == ot {
			c.timers[index].running = false
			continue
		}
	}
}

func (c *nTimer) runWorking(j func()) {
	defer utils.CatchError(j)
	j()
}

func (c *nTimer) run() {
	t := time.NewTicker(time.Second)
	n := 0

	for {
		select {
		case <-c.begin:
			n++
		case <-func() <-chan time.Time {
			if n > 0 {
				return t.C
			}
			return nil
		}():
			if len(c.timers) > 0 {
				cro := c.timers[0]

				if !cro.running {
					c.timers = append(c.timers[:0], c.timers[1:]...)
					n--
					continue
				}

				if time.Now().After(cro.next) {
					n--
					if cro.running {
						go c.runWorking(cro.fn)

						if cro.cType == clockOnce {
							c.timers = append(c.timers[:0], c.timers[1:]...)
						}

						if cro.cType == clockForever {
							c.timers[0].next = time.Now().Add(cro.sec)
							sort.Sort(clockTime(c.timers))
						}
					} else {
						c.timers = append(c.timers[:0], c.timers[1:]...)
						continue
					}
				}
			}
		}
	}
}

func (c *nTimer) add(ot string, cType clockType, fn func(), sec time.Duration) {

	c.timers = append(
		c.timers,
		&timer{
			oType:   ot,
			cType:   cType,
			fn:      fn,
			sec:     sec,
			next:    time.Now().Add(sec),
			running: true,
		},
	)

	sort.Sort(clockTime(c.timers))

	if !c.timers[0].running {
		c.timers = append(c.timers[:0], c.timers[1:]...)
	}

	c.begin <- struct{}{}
}

func (c *nTimer) now() time.Time { return time.Now().In(c.location) }

type clockTime []*timer

func (s clockTime) Len() int { return len(s) }

func (s clockTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s clockTime) Less(i, j int) bool {
	if s[i].next.IsZero() {
		return false
	}
	if s[j].next.IsZero() {
		return true
	}
	return s[i].next.Before(s[j].next)
}
