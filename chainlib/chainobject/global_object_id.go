package chainobject

import (
	"datx_chain/utils/helper"
	"sync"
	"sync/atomic"
)

// Object ID that auto increment
type oid struct {
	//
	id map[uint32]uint32
}

func (self *oid) add(objtyp uint32) uint64 {
	var id uint32

	if data, ok := self.id[objtyp]; ok {
		//add 1 to id atomic
		id = atomic.AddUint32(&data, 1)
	} else {
		id = 0
	}
	self.id[objtyp] = id

	//encode two uint32 to one uint64
	new_id := helper.EncodeBit(objtyp, id)

	return new_id
}

var guid *oid
var once sync.Once

//get object id that auto increment global, return 0 if occurred uint64 overflow
func GetOID(objtyp uint32) uint64 {
	once.Do(func() {
		guid = &oid{
			id: make(map[uint32]uint32, 0),
		}
	})

	return guid.add(objtyp)
}
