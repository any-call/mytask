package mytask

type BaseTask struct {
	id int64
}

func (self *BaseTask) ID() int64 {
	return self.id
}

func (self *BaseTask) SetID(id int64) {
	self.id = id
}

func (self *BaseTask) Cmd() func() {
	panic("Cmd() method must be implemented by the concrete page")
}
