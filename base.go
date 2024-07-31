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

func add(task ScheduleTask, spec string, runImmediately bool) {
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

		if err := c.AddFunc(spec, task.Cmd()); err != nil {
			panic(err)
		}
		c.Start() // 启动 cron 调度器

		if runImmediately { //建立任务后立即运行
			go task.Cmd()()
		}
	}
}

func AddThenStart(task ScheduleTask, spec string, runImmediately bool) {
	add(task, spec, runImmediately)
}
