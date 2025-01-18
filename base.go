package mytask

import (
	"fmt"
	"github.com/any-call/gobase/util/mylog"
	"github.com/any-call/gobase/util/mymap"
	"github.com/robfig/cron/v3"
	"time"
)

type ScheduleTask interface {
	ID() int64
	Cmd() func()
}

var (
	taskMap         = mymap.NewMap[int64, ScheduleTask]() //map[int]ScheduleTask{} =
	cronAndTimerMap = mymap.NewMap[int64, any]()          //map[int]*cron.Cron{}
)

func add(task ScheduleTask, spec string, runImmediately bool) {
	if _, ok := taskMap.Value(task.ID()); ok {
		mylog.Debug(fmt.Errorf("add  task ID %d exist ", task.ID()))
		return
	}

	if _, ok := cronAndTimerMap.Value(task.ID()); ok {
		mylog.Debug(fmt.Errorf("add  task ID %d exist ", task.ID()))
		return
	}

	taskMap.Insert(task.ID(), task)
	{
		c := cron.New()
		cronAndTimerMap.Insert(task.ID(), c)

		if _, err := c.AddFunc(spec, task.Cmd()); err != nil {
			panic(err)
		}
		c.Start() // 启动 cron 调度器

		if runImmediately { //建立任务后立即运行
			go task.Cmd()()
		}
	}
}

func addTimeTask(task ScheduleTask, t time.Duration, asyncTask bool, runImmediately bool) {
	if _, ok := taskMap.Value(task.ID()); ok {
		mylog.Debug(fmt.Errorf("add  task ID %d exist ", task.ID()))
		return
	}

	if _, ok := cronAndTimerMap.Value(task.ID()); ok {
		mylog.Debug(fmt.Errorf("add  task ID %d exist ", task.ID()))
		return
	}

	taskMap.Insert(task.ID(), task)
	{
		c := NewTimerTask(t, task)
		cronAndTimerMap.Insert(task.ID(), c)
		c.Start(asyncTask) // 启动 cron 调度器

		if runImmediately { //建立任务后立即运行
			go task.Cmd()()
		}
	}
}

func AddThenStart(task ScheduleTask, spec string, runImmediately bool) {
	add(task, spec, runImmediately)
}

func AddTimerThenStart(task ScheduleTask, t time.Duration, asyncTask, runImmediately bool) {
	addTimeTask(task, t, asyncTask, runImmediately)
}

func GetTaskModel(taskId int64) (any, bool) {
	return taskMap.Value(taskId)
}

func IsExist(id int64) bool {
	if _, ok := cronAndTimerMap.Value(id); ok {
		return true
	}

	return false
}

func Stop(id int64) {
	if c, ok := cronAndTimerMap.Value(id); ok {
		fmt.Println("2: will stop task:", id)
		if c1, ok := c.(*cron.Cron); ok {
			c1.Stop()
		} else {
			if c2, ok := c.(*TimerTask); ok {
				c2.Stop()
			}
		}
	}
}

func Remove(id int64) {
	if _, ok := taskMap.Value(id); ok {
		fmt.Println("remove  task ID:", id)
		taskMap.Remove(id)
	}

	if c, ok := cronAndTimerMap.Value(id); ok {
		fmt.Println("1: will stop task:", id)
		if c1, ok := c.(*cron.Cron); ok {
			c1.Stop()
		} else {
			if c2, ok := c.(*TimerTask); ok {
				c2.Cancel()
			}
		}
		cronAndTimerMap.Remove(id)
	}
}

func ResetCron(id int64, spec string) error {
	if t, ok := taskMap.Value(id); ok {
		if c, okk := cronAndTimerMap.Value(id); okk {
			fmt.Println("3:will stop task:", id)
			if c1, ok := c.(*cron.Cron); ok {
				listEntry := c1.Entries()
				if listEntry != nil {
					for i, _ := range listEntry {
						c1.Remove(listEntry[i].ID)
					}
				}

				c1.Stop()
				cronAndTimerMap.Remove(id)
				cc := cron.New()
				if _, err := cc.AddFunc(spec, t.Cmd()); err != nil {
					return err
				}
				cronAndTimerMap.Insert(id, cc)
				cc.Start()
				return nil
			}
		}
		return fmt.Errorf("incorrect cron id:%d", id)
	}

	return fmt.Errorf("incorrect task id:%d", id)
}

func ResetTimer(id int64, tm time.Duration, asyncTask bool) error {
	if t, ok := taskMap.Value(id); ok {
		if c, okk := cronAndTimerMap.Value(id); okk {
			fmt.Println("3:will stop task:", id)
			if c1, ok := c.(*TimerTask); ok {
				c1.Cancel()
				cronAndTimerMap.Remove(id)
				cc := NewTimerTask(tm, t)
				cronAndTimerMap.Insert(id, cc)
				cc.Start(asyncTask)
				return nil
			}
		}
		return fmt.Errorf("incorrect cron id:%d", id)
	}

	return fmt.Errorf("incorrect task id:%d", id)
}
