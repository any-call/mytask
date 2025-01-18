package mytask

import (
	"sync"
	"time"
)

type TimerTask struct {
	sync.Mutex
	duration   time.Duration
	task       ScheduleTask
	stop       chan struct{}
	running    bool
	cancelOnce sync.Once
}

func NewTimerTask(duration time.Duration, t ScheduleTask) *TimerTask {
	return &TimerTask{
		Mutex:      sync.Mutex{},
		duration:   duration,
		task:       t,
		stop:       make(chan struct{}, 1),
		running:    false,
		cancelOnce: sync.Once{},
	}
}

// 启动定时任务
func (self *TimerTask) Start(asyncTask bool) {
	self.Lock()
	defer self.Unlock()

	if self.running {
		return
	}

	self.running = true
	go func() {
		timer := time.NewTimer(self.duration)
		for {
			select {
			case <-timer.C:
				if self.task != nil {
					if asyncTask { //异步任务，协程运行
						go self.task.Cmd()()
					} else {
						self.task.Cmd()()
					}
				}
				timer.Reset(self.duration)
				break
			case <-self.stop:
				return
			}
		}
	}()
}

// 停止任务，重置定时器
func (self *TimerTask) Stop() {
	self.Lock()
	defer self.Unlock()

	if !self.running {
		return
	}

	self.running = false
	self.stop <- struct{}{} //退出协程
}

// Cancel ：取消停任，释放资源
func (self *TimerTask) Cancel() {
	self.cancelOnce.Do(func() {
		self.Stop()
		close(self.stop)
	})
}
