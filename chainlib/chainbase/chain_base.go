package chainbase

import (
	"container/list"
	"datx_chain/chainlib/chainobject"
	"datx_chain/utils/db"
	"datx_chain/utils/helper"
	"datx_chain/utils/rlp"
	"errors"
	"fmt"
	"log"
)

const (
	account_type uint32 = 1 << iota
	transcation_type
	max_type
)

type Generic_Index struct {
	//dequene of undo_state
	stack *list.List

	//
	revision int64

	//
	next_id uint64

	//value type
	value_type uint32

	//state db
	db *datxdb.LDBDatabase
}

func NewGenericIndex(valuetype uint32, db *datxdb.LDBDatabase) *Generic_Index {
	return &Generic_Index{
		stack:      list.New(),
		revision:   0,
		next_id:    0,
		value_type: valuetype,
		db:         db,
	}
}

//get global id of this raw object
func (self *Generic_Index) getid(v interface{}) uint64 {
	var id uint64

	//return 0 if the input v is invalid
	defer func() {
		if err := recover(); err != nil {
			id = 0
		}
	}()

	switch {
	case self.value_type == account_type:
		id = v.(*chainobject.AccountObject).ID()
	default:
		log.Printf("getid on undefined value type={%v}", self.value_type)
		id = 0
	}

	return id
}

//get rlp data of this raw object
func (self *Generic_Index) getrlp(v interface{}) ([]byte, error) {
	data, err := rlp.EncodeToBytes(v)
	if err != nil {
		log.Printf("get rlp err={%v}", err)
		return []byte("0"), err
	}
	return data, nil
}

//get hash id of rlp object
func (self *Generic_Index) fromrlp(v []byte) uint64 {
	var id uint64

	switch {
	case self.value_type == account_type:
		var value *chainobject.AccountObject
		if err := rlp.DecodeBytes(v, &value); err != nil {
			log.Printf("decode rlp err={%v}", err)
			return 0
		}
		id = value.ID()
	default:
		log.Printf("getid on undefined value type={%v}", self.value_type)
	}

	return id
}

func (self *Generic_Index) IsValidate(objtype uint32) bool {
	return self.value_type == objtype
}

func (self *Generic_Index) Enabled() bool {
	return self.stack.Len() > 0
}

func (self *Generic_Index) Start_Undo_Session(enabled bool) *Session {
	if enabled {
		new_revision := self.revision + 1
		self.revision = new_revision

		state := NewUndoState(self.next_id, new_revision)
		self.stack.PushBack(state)

		return NewSession(self, self.revision)
	}

	return NewSession(self, -1)
}

func (self *Generic_Index) Undo() error {
	if !self.Enabled() {
		return errors.New("undo stack is empty")
	}

	temp := self.stack.Back()
	head := temp.Value.(*undo_state)

	//recover the old values
	for k, v := range head.old_values {
		if err := self.db.Put(helper.ToBytes(k), v); err != nil {
			return err
		}
	}

	//recover the new values
	for k, _ := range head.new_values {
		if err := self.db.Delete(helper.ToBytes(k)); err != nil {
			return err
		}
	}

	self.next_id = head.old_next_id

	//recover the removed values
	for k, v := range head.removed_values {
		if err := self.db.Put(helper.ToBytes(k), v); err != nil {
			return err
		}
	}

	self.stack.Remove(temp)
	new := self.revision - 1
	self.revision = new
	return nil
}

// An object's relationship to a state can be:
// in new_ids            : new
// in old_values (was=X) : upd(was=X)
// in removed (was=X)    : del(was=X)
// not in any of above   : nop
//
// When merging A=prev_state and B=state we have a 4x4 matrix of all possibilities:
//
//                   |--------------------- B ----------------------|
//
//                +------------+------------+------------+------------+
//                | new        | upd(was=Y) | del(was=Y) | nop        |
//   +------------+------------+------------+------------+------------+
// / | new        | N/A        | new       A| nop       C| new       A|
// | +------------+------------+------------+------------+------------+
// | | upd(was=X) | N/A        | upd(was=X)A| del(was=X)C| upd(was=X)A|
// A +------------+------------+------------+------------+------------+
// | | del(was=X) | N/A        | N/A        | N/A        | del(was=X)A|
// | +------------+------------+------------+------------+------------+
// \ | nop        | new       B| upd(was=Y)B| del(was=Y)B| nop      AB|
//   +------------+------------+------------+------------+------------+
//
// Each entry was composed by labelling what should occur in the given case.
//
// Type A means the composition of states contains the same entry as the first of the two merged states for that object.
// Type B means the composition of states contains the same entry as the second of the two merged states for that object.
// Type C means the composition of states contains an entry different from either of the merged states for that object.
// Type N/A means the composition of states violates causal timing.
// Type AB means both type A and type B simultaneously.
//
// The merge() operation is defined as modifying prev_state in-place to be the state object which represents the composition of
// state A and B.
//
// Type A (and AB) can be implemented as a no-op; prev_state already contains the correct value for the merged state.
// Type B (and AB) can be implemented by copying from state to prev_state.
// Type C needs special case-by-case logic.
// Type N/A can be ignored or assert(false) as it can only occur if prev_state and state have illegal values
// (a serious logic error which should never happen).
//
func (self *Generic_Index) Squash() error {
	if !self.Enabled() {
		return errors.New("squash stack is empty")
	}

	//get head and prev element
	head := self.stack.Back()
	self.stack.Remove(head)
	prev := self.stack.Back()

	state := head.Value.(*undo_state)
	prev_state := prev.Value.(*undo_state)

	// We can only be outside type A/AB (the nop path) if B is not nop, so it suffices to iterate through B's three containers.
	for _, v := range state.old_values {
		id := self.fromrlp(v)

		// new+upd -> new, type A
		if _, ok := prev_state.new_values[id]; ok {
			continue
		}

		// upd(was=X) + upd(was=Y) -> upd(was=X), type A
		if _, ok := prev_state.old_values[id]; ok {
			continue
		}

		// del+upd -> N/A
		if _, ok := prev_state.removed_values[id]; ok {
			return errors.New("the modify operation after removed is not allowed")
		}

		prev_state.old_values[id] = v
	}

	// *+new, but we assume the N/A cases don't happen, leaving type B nop+new -> new
	for k, v := range state.new_values {
		prev_state.new_values[k] = v
	}

	// *+del
	for k, v := range state.removed_values {
		id := self.fromrlp(v)

		// new + del -> nop (type C)
		if _, ok := prev_state.new_values[id]; ok {
			delete(prev_state.new_values, id)
			continue
		}

		// upd(was=X) + del(was=Y) -> del(was=X)
		if _, ok := prev_state.old_values[id]; ok {
			prev_state.removed_values[id] = v
			delete(prev_state.old_values, id)

			continue
		}

		// del+del-> N/A
		if _, ok := prev_state.removed_values[id]; ok {
			return errors.New("the remove operation after removed is not allowed")
		}

		// nop + del(was=Y) -> del(was=Y)
		prev_state.removed_values[k] = v
	}

	new := self.revision - 1
	self.revision = new
	return nil
}

func (self *Generic_Index) Commit(revision int64) {
	if !self.Enabled() {
		return
	}

	firele := self.stack.Front()
	first := firele.Value.(*undo_state)
	for self.Enabled() && first.revision <= revision {
		self.stack.Remove(firele)
		firele = self.stack.Front()
		first = firele.Value.(*undo_state)
	}
}

func (self *Generic_Index) UndoAll() error {
	for self.Enabled() {
		err := self.Undo()
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Generic_Index) SetRevision(revision int64) {
	if !self.Enabled() {
		self.revision = revision
	}
}

// the purpose of on_modify/create/remove is to save undo state for the method of Undo()
//update old values of undo_state
func (self *Generic_Index) on_modify(key uint64, value []byte) error {
	if !self.Enabled() {
		return errors.New("on_modify stack is empty.")
	}

	head := self.stack.Back().Value.(*undo_state)

	//if found on new valus,do nothing
	if _, ok := head.new_values[key]; ok {
		return nil
	}

	//if found on old values,do nothing
	if _, ok := head.old_values[key]; ok {
		return nil
	}

	//you should update old values when there is no record on the last back of stack
	head.old_values[key] = value
	return nil
}

func (self *Generic_Index) on_remove(key uint64, value []byte) error {
	if !self.Enabled() {
		return errors.New("on_remove stack is empty.")
	}

	head := self.stack.Back().Value.(*undo_state)

	//if found, remove from the new values
	if _, ok := head.new_values[key]; ok {
		delete(head.new_values, key)
		return nil
	}

	//if found on old values, insert and remove
	if _, ok := head.old_values[key]; ok {
		head.removed_values[key] = value
		delete(head.old_values, key)
		return nil
	}

	//if found on removed values, do nothing
	if _, ok := head.removed_values[key]; ok {
		return errors.New("on_remove remove repetition.")
	}

	head.removed_values[key] = value
	return nil
}

func (self *Generic_Index) on_create(key uint64, value []byte) error {
	if !self.Enabled() {
		return errors.New("on_create stack is empty.")
	}

	head := self.stack.Back().Value.(*undo_state)

	head.new_values[key] = value
	return nil
}

//support Insert,Remove,Find operation for the objects
func (self *Generic_Index) Insert(v interface{}) error {
	self.next_id = self.next_id + 1

	//get global id of input object
	id := self.getid(v)
	if id == 0 {
		str := fmt.Sprintf("Insert on unreg obj value={%v}", v)
		return errors.New(str)
	}

	//get rlp data
	data, err := self.getrlp(v)
	if err != nil {
		str := fmt.Sprintf("Insert obj={%v} err={%v}", v, err)
		return errors.New(str)
	}

	//insert value into levelDB
	if err := self.db.Put(helper.ToBytes(id), data); err != nil {
		str := fmt.Sprintf("Insert obj id={%v} err={%v}", id, err)
		return errors.New(str)
	}

	return self.on_create(id, data)
}

func (self *Generic_Index) Modify(v interface{}) error {
	//get global id of input object
	id := self.getid(v)
	if id == 0 {
		str := fmt.Sprintf("Modify on unreg obj value={%v}", v)
		return errors.New(str)
	}

	//get rlp data
	data, err := self.getrlp(v)
	if err != nil {
		str := fmt.Sprintf("Modify obj={%v} err={%v}", v, err)
		return errors.New(str)
	}

	//get old value
	old, err := self.db.Get(helper.ToBytes(id))
	if err != nil {
		str := fmt.Sprintf("Modify get old obj id={%v} err={%v}", id, err)
		return errors.New(str)
	}

	//insert new value into levelDB
	if err := self.db.Put(helper.ToBytes(id), data); err != nil {
		str := fmt.Sprintf("Modify obj id={%v} err={%v}", id, err)
		return errors.New(str)
	}

	return self.on_modify(id, old)
}

func (self *Generic_Index) Remove(v interface{}) error {
	//get global id of input object
	id := self.getid(v)
	if id == 0 {
		str := fmt.Sprintf("Remove on unreg obj value={%v}", v)
		return errors.New(str)
	}

	//get old value
	old, err := self.db.Get(helper.ToBytes(id))
	if err != nil {
		str := fmt.Sprintf("Remove get old obj id={%v} err={%v}", id, err)
		return errors.New(str)
	}

	//delete value from levelDB
	if err := self.db.Delete(helper.ToBytes(id)); err != nil {
		str := fmt.Sprintf("chain base remove id={%v} err={%v}", id, err)
		log.Printf(str)
		return errors.New(str)
	}

	return self.on_remove(id, old)
}

func (self *Generic_Index) Find(key uint64) (interface{}, error) {
	data, err := self.db.Get(helper.ToBytes(key))
	if err != nil {
		return nil, err
	}

	switch {
	case self.value_type == account_type:
		var value *chainobject.AccountObject
		if err := rlp.DecodeBytes(data, &value); err != nil {
			return nil, err
		}
		return value, nil
	default:
		str := fmt.Sprintf("Find on unreg value type={%v}", self.value_type)
		log.Printf(str)
		return nil, errors.New(str)
	}
}
