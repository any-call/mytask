package mytask

import (
	"github.com/any-call/gobase/util/mylog"
	"testing"
	"time"
)

type MyTask struct {
}

func (self *MyTask) ID() int64 {
	return 1
}

func (self *MyTask) Cmd() func() {
	return func() {
		mylog.Debug("run here ")
	}
}

func TestNewTimerTask(t *testing.T) {
	task := NewTimerTask(time.Second, &MyTask{})
	task.Start(false)
	time.Sleep(time.Second * 10)
}
