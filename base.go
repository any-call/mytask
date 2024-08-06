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

func IsExist(id int) bool {
	if _, ok := cronMap.Value(id); ok {
		return true
	}

	return false
}

func Stop(id int) {
	if c, ok := cronMap.Value(id); ok {
		fmt.Println("2: will stop task:", id)
		c.Stop()
	}
}

func Remove(id int) {
	if _, ok := taskMap.Value(id); ok {
		fmt.Println("remove  task ID:", id)
		taskMap.Remove(id)
	}
	if c, ok := cronMap.Value(id); ok {
		fmt.Println("1: will stop task:", id)
		c.Stop()
		cronMap.Remove(id)
	}
}

func Refresh(id int, spec string) error {
	if t, ok := taskMap.Value(id); ok {
		if c, okk := cronMap.Value(id); okk {
			fmt.Println("3:will stop task:", id)
			c.Stop()
			cronMap.Remove(id)
			cc := cron.New()
			if err := cc.AddFunc(spec, t.Cmd()); err != nil {
				return err
			}
			cronMap.Insert(id, cc)
			cc.Start()
			return nil
		}
		return fmt.Errorf("incorrect cron id:%d", id)
	}

	return fmt.Errorf("incorrect task id:%d", id)
}
