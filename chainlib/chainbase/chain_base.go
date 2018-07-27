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

type Generic_Index struct {
	//dequene of undo_state
	stack *list.List

	//
	revision int64

	//
	nextid uint64

	//value type
	value_type uint32

	//state db
	db *datxdb.LDBDatabase
}

func NewGenericIndex(valuetype uint32, db *datxdb.LDBDatabase) *Generic_Index {
	return &Generic_Index{
		stack:      list.New(),
		revision:   0,
		nextid:     0,
		value_type: valuetype,
		db:         db,
	}
}

//get rlp data of this raw object
func (index *Generic_Index) getrlp(v interface{}) ([]byte, error) {
	data, err := rlp.EncodeToBytes(v)
	if err != nil {
		log.Printf("get rlp err={%v}", err)
		return []byte("0"), err
	}
	return data, nil
}

func (index *Generic_Index) IsValidate(objtype uint32) bool {
	return index.value_type == objtype
}

func (index *Generic_Index) Enabled() bool {
	return index.stack.Len() > 0
}

func (index *Generic_Index) Start_Undo_Session(enabled bool) *Session {
	if enabled {
		new_revision := index.revision + 1
		index.revision = new_revision

		state := NewUndoState(index.nextid, new_revision)
		index.stack.PushBack(state)

		return NewSession(index, index.revision)
	}

	return NewSession(index, -1)
}

func (index *Generic_Index) Undo() error {
	if !index.Enabled() {
		return errors.New("undo stack is empty")
	}

	temp := index.stack.Back()
	head := temp.Value.(*undo_state)

	//recover the old values
	for k, v := range head.old_values {
		if err := index.db.Put(helper.ToBytes(k), v); err != nil {
			return err
		}
	}

	//recover the new values
	for k, _ := range head.new_values {
		if err := index.db.Delete(helper.ToBytes(k)); err != nil {
			return err
		}
	}

	index.nextid = head.old_next_id

	//recover the removed values
	for k, v := range head.removed_values {
		if err := index.db.Put(helper.ToBytes(k), v); err != nil {
			return err
		}
	}

	index.stack.Remove(temp)
	new := index.revision - 1
	index.revision = new
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
func (index *Generic_Index) Squash() error {
	if !index.Enabled() {
		return errors.New("squash stack is empty")
	}

	//get head and prev element
	head := index.stack.Back()
	index.stack.Remove(head)
	prev := index.stack.Back()

	state := head.Value.(*undo_state)
	prev_state := prev.Value.(*undo_state)

	// We can only be outside type A/AB (the nop path) if B is not nop, so it suffices to iterate through B's three containers.
	for _, v := range state.old_values {
		id := chainobject.GetIDFromRLPData(index.value_type, v)

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
		id := chainobject.GetIDFromRLPData(index.value_type, v)

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

	new := index.revision - 1
	index.revision = new
	return nil
}

func (index *Generic_Index) Commit(revision int64) {
	if !index.Enabled() {
		return
	}

	firele := index.stack.Front()
	first := firele.Value.(*undo_state)
	for index.Enabled() && first.revision <= revision {
		index.stack.Remove(firele)
		firele = index.stack.Front()
		first = firele.Value.(*undo_state)
	}
}

func (index *Generic_Index) UndoAll() error {
	for index.Enabled() {
		err := index.Undo()
		if err != nil {
			return err
		}
	}

	return nil
}

func (index *Generic_Index) SetRevision(revision int64) {
	if !index.Enabled() {
		index.revision = revision
	}
}

func (index *Generic_Index) GetRevision() int64 {
	return index.revision
}

// the purpose of on_modify/create/remove is to save undo state for the method of Undo()
//update old values of undo_state
func (index *Generic_Index) on_modify(key uint64, value []byte) error {
	if !index.Enabled() {
		return errors.New("on_modify stack is empty.")
	}

	head := index.stack.Back().Value.(*undo_state)

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

func (index *Generic_Index) on_remove(key uint64, value []byte) error {
	if !index.Enabled() {
		return errors.New("on_remove stack is empty.")
	}

	head := index.stack.Back().Value.(*undo_state)

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

func (index *Generic_Index) on_create(key uint64, value []byte) error {
	if !index.Enabled() {
		return errors.New("on_create stack is empty.")
	}

	head := index.stack.Back().Value.(*undo_state)

	head.new_values[key] = value
	return nil
}

//Insert support Insert,Remove,Find operation for the objects
func (index *Generic_Index) Insert(v interface{}) error {
	index.nextid = index.nextid + 1

	//get global id of input object
	id := chainobject.GetIDFromRaw(index.value_type, v)

	//get rlp data
	data, err := index.getrlp(v)
	if err != nil {
		str := fmt.Sprintf("Insert obj={%v} err={%v}", v, err)
		return errors.New(str)
	}

	//insert value into levelDB
	if err := index.db.Put(helper.ToBytes(id), data); err != nil {
		str := fmt.Sprintf("Insert obj id={%v} err={%v}", id, err)
		return errors.New(str)
	}

	return index.on_create(id, data)
}

func (index *Generic_Index) Modify(v interface{}) error {
	//get global id of input object
	id := chainobject.GetIDFromRaw(index.value_type, v)
	if id == 0 {
		str := fmt.Sprintf("Modify on unreg obj value={%v}", v)
		return errors.New(str)
	}

	//get rlp data
	data, err := index.getrlp(v)
	if err != nil {
		str := fmt.Sprintf("Modify obj={%v} err={%v}", v, err)
		return errors.New(str)
	}

	//get old value
	old, err := index.db.Get(helper.ToBytes(id))
	if err != nil {
		str := fmt.Sprintf("Modify get old obj id={%v} err={%v}", id, err)
		return errors.New(str)
	}

	//insert new value into levelDB
	if err := index.db.Put(helper.ToBytes(id), data); err != nil {
		str := fmt.Sprintf("Modify obj id={%v} err={%v}", id, err)
		return errors.New(str)
	}

	return index.on_modify(id, old)
}

func (index *Generic_Index) Remove(v interface{}) error {
	//get global id of input object
	id := chainobject.GetIDFromRaw(index.value_type, v)
	if id == 0 {
		str := fmt.Sprintf("Remove on unreg obj value={%v}", v)
		return errors.New(str)
	}

	//get old value
	old, err := index.db.Get(helper.ToBytes(id))
	if err != nil {
		str := fmt.Sprintf("Remove get old obj id={%v} err={%v}", id, err)
		return errors.New(str)
	}

	//delete value from levelDB
	if err := index.db.Delete(helper.ToBytes(id)); err != nil {
		str := fmt.Sprintf("chain base remove id={%v} err={%v}", id, err)
		log.Printf(str)
		return errors.New(str)
	}

	return index.on_remove(id, old)
}

func (index *Generic_Index) Find(key uint64) (interface{}, error) {
	data, err := index.db.Get(helper.ToBytes(key))
	if err != nil {
		return nil, err
	}

	return chainobject.GetValue(index.value_type, data)
}
