package mytask

import (
	"fmt"
	"github.com/any-call/gobase/util/mylog"
	"github.com/any-call/gobase/util/mymap"
	"github.com/robfig/cron"
)

type ScheduleTask interface {
	ID() int
	Cmd() func()
}

var (
	taskMap = mymap.NewMap[int, ScheduleTask]() //map[int]ScheduleTask{} =
	cronMap = mymap.NewMap[int, *cron.Cron]()   //map[int]*cron.Cron{}
)

func Add(task ScheduleTask, spec string, start bool) {
	if _, ok := taskMap.Value(task.ID()); ok {
		mylog.Debug(fmt.Errorf("add  task ID %d exist ", task.ID()))
		return
	}

	if _, ok := cronMap.Value(task.ID()); ok {
		mylog.Debug(fmt.Errorf("add  task ID %d exist ", task.ID()))
		return
	}

	taskMap.Insert(task.ID(), task)
	{
		c := cron.New()
		cronMap.Insert(task.ID(), c)

		if start {
			if err := c.AddFunc(spec, task.Cmd()); err != nil {
				panic(err)
			}
			c.Start()
		}
	}
}

func AddThenStart(task ScheduleTask, spec string) {
	Add(task, spec, true)
}
