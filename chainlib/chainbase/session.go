package chainbase

type Session struct {
	//Gernic_Index instace
	index *Generic_Index

	//open/close switch of operation
	apply bool

	//current revision
	revision int64
}

func NewSession(in *Generic_Index, rev int64) *Session {
	var result Session
	result.index = in
	result.revision = rev

	if rev == -1 {
		result.apply = false
	} else {
		result.apply = true
	}

	return &result
}

func (self *Session) Close() {
	if self.apply {
		self.index.Undo()
	}
}

func (self *Session) Push() {
	self.apply = false
}

func (self *Session) Squash() {
	if self.apply {
		self.index.Squash()
	}

	self.apply = false
}

func (self *Session) Undo() {
	if self.apply {
		self.index.Undo()
	}

	self.apply = false
}

func (self *Session) Revision() int64 {
	return self.revision
}

type SessionList struct {
	revision int64

	list []*Session
}

func NewSessionSet(list []*Session) SessionList {
	var result SessionList
	result.list = list

	if len(list) > 0 {
		result.revision = list[0].revision
	} else {
		result.revision = -1
	}

	return result
}

func (self *SessionList) Close() {
	self.Undo()
}

func (self *SessionList) Undo() {
	for _, v := range self.list {
		v.Undo()
	}

	self.list = make([]*Session, 0)
}

func (self *SessionList) Squash() {
	for _, v := range self.list {
		v.Squash()
	}

	self.list = make([]*Session, 0)
}

func (self *SessionList) Push() {
	for _, v := range self.list {
		v.Push()
	}

	self.list = make([]*Session, 0)
}

func (self *SessionList) Revision() int64 {
	return self.revision
}
