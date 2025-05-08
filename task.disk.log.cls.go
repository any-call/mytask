package mytask

import (
	"fmt"
	"github.com/any-call/gobase/util/mycmd"
	"github.com/any-call/gobase/util/mylog"
	"sync"
)

type diskLogClear struct {
	sync.Mutex
	BaseTask
	maxLogSizeMB int
}

func NewDiskLogCls(id int, maxLogSizeMB int) ScheduleTask {
	t := &diskLogClear{maxLogSizeMB: maxLogSizeMB}
	t.SetID(int64(id))
	return t
}

func (self *diskLogClear) Cmd() func() {
	return func() {
		if !self.TryLock() {
			return
		}
		defer self.Unlock()

		mylog.Info("enter disk log ")
		outstr, err := mycmd.Exec("journalctl", nil, false, fmt.Sprintf("--vacuum-size=%dM", self.maxLogSizeMB))
		if err != nil {
			mylog.Debug("disk log err :", err)
			return
		}

		mylog.Debug("disk log ok :", outstr)
	}
}
