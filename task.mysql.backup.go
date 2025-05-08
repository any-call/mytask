package mytask

import (
	"github.com/any-call/gobase/frame/mysql"
	"github.com/any-call/gobase/util/mycmd"
	"github.com/any-call/gobase/util/mylog"
	"github.com/any-call/gobase/util/myos"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type mysqlDbBackup struct {
	BaseTask
	sync.Mutex
	dbUser     string
	dbPass     string
	dbName     string
	backupPath string
	dockerName string //标识mysql 是不是运行在docker 中，如果是请传入docker 容器名称
	isSudo     bool   //标识是否增加 sudo
}

func NewMysqlDBBackup(id int, dbUser, dbPass, dbName, backupPath, dockerName string, isSudo bool) ScheduleTask {
	t := &mysqlDbBackup{dbUser: dbUser, dbPass: dbPass, dbName: dbName, backupPath: backupPath, dockerName: dockerName, isSudo: isSudo}
	t.SetID(int64(id))
	return t
}

func (self *mysqlDbBackup) Cmd() func() {
	return func() {
		if !self.TryLock() {
			return
		}
		defer self.Unlock()

		startTime := time.Now()
		mylog.Info("enter mysqldb backup ")
		fullfile, err := mysql.BackupMySQL(self.dbUser, self.dbPass, self.dbName, self.backupPath, func() []string {
			ret := []string{}
			if self.isSudo {
				ret = []string{"sudo"}
			}
			if len(self.dockerName) > 0 {
				ret = append(ret, mysql.DockerExecCmd("mysql8")...)
			}

			return ret
		})

		if err != nil {
			mylog.Debug("mysqldb backup err :", err)
			return
		}

		mylog.Debugf("mysqldb backup ok ;duration:%.2fm ;file is :%s", time.Now().Sub(startTime).Minutes(), fullfile)
		///清除掉旧的备份文件
		list, err := myos.FindFilesWithExt(self.backupPath, "sql")
		if err != nil {
			mylog.Debug("get list file err :", err)
			return
		}

		for _, filename := range list {
			if !strings.HasSuffix(fullfile, filename) {
				if _, err := mycmd.Exec("rm", func(c *exec.Cmd) {
					c.Dir = self.backupPath
				}, false, "-f", filename); err != nil {
					mylog.Debug("remove file err:", err)
				}
			}
		}
	}
}
