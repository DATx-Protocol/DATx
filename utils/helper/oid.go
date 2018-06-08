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
	raw := atomic.LoadUint64(&self.id)

	id, ok := math.SafeAdd(raw, 1)
	if ok {
		return 0
	}

	return id
}

var guid *oid
var once sync.Once

//get object id that auto increment global, return 0 if occurred uint64 overflow
func GetOID() uint64 {
	once.Do(func() {
		guid = &oid{
			id: 1,
		}
	})

	return guid.add()
}
