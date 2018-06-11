package helper

import (
	"datx_chain/utils/common/math"
	"sync"
	"sync/atomic"
)

// Object ID that auto increment
type oid struct {
	//
	id uint64
}

func (self *oid) add() uint64 {
	//check whether occurred overflow
	_, ok := math.SafeAdd(self.id, 1)
	if ok {
		return 0
	}
	//add 1 to id atomic
	id := atomic.AddUint64(&self.id, 1)

	return id
}

var guid *oid
var once sync.Once

//get object id that auto increment global, return 0 if occurred uint64 overflow
func GetOID() uint64 {
	once.Do(func() {
		guid = &oid{
			id: 0,
		}
	})

	return guid.add()
}
