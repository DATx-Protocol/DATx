package chainbase

import (
	"datx_chain/utils/db"
	"errors"
	"fmt"
	"log"
)

type DataBaser interface {
	Insert(v interface{})
	Modify(v interface{})
	Remove(v interface{})
	Find(key uint64) (interface{}, error)
}

type UndoStater interface {
	Undo() error
	Squash() error
	Commit(revision int64)
	UndoAll() error
}

type DataBase struct {
	//the pairs of object type and generic index instance
	index_map map[uint32]*Generic_Index

	//data dir
	path string

	db *datxdb.LDBDatabase
}

func NewDataBase(path string) *DataBase {
	var result DataBase
	result.path = path
	result.index_map = make(map[uint32]*Generic_Index)
	result.db, _ = datxdb.NewLDBDatabase(path, 0, 0)

	return &result
}

//create and insert object
func (self *DataBase) Create(objtype uint32, obj interface{}) error {
	chain, ok := self.index_map[objtype]
	if !ok {
		str := fmt.Sprintf("Create on unregistered obj type={%v}", objtype)
		log.Printf(str)
		return errors.New(str)
	}

	if !chain.IsValidate(objtype) {
		str := fmt.Sprintf("Create on missmatched obj type={%v}", objtype)
		log.Printf(str)
		return errors.New(str)
	}

	return chain.Insert(obj)
}

func (self *DataBase) Remove(objtype uint32, obj interface{}) error {
	chain, ok := self.index_map[objtype]
	if !ok {
		str := fmt.Sprintf("Remove on unregister obj type={%v}", objtype)
		log.Printf(str)
		return errors.New(str)
	}

	if !chain.IsValidate(objtype) {
		str := fmt.Sprintf("Remove on missmatched obj type={%v}", objtype)
		log.Printf(str)
		return errors.New(str)
	}

	return chain.Remove(obj)
}

//the newvalue is the pointer to raw object
func (self *DataBase) Modify(objtype uint32, newvalue interface{}) error {
	chain, ok := self.index_map[objtype]
	if !ok {
		str := fmt.Sprintf("Modify on unregister obj type={%v}", objtype)
		log.Printf(str)
		return errors.New(str)
	}

	if !chain.IsValidate(objtype) {
		str := fmt.Sprintf("Modify on missmatched obj type={%v}", objtype)
		log.Printf(str)
		return errors.New(str)
	}

	return chain.Modify(newvalue)
}

func (self *DataBase) Get(objtype uint32, key uint64) (interface{}, error) {
	chain, ok := self.index_map[objtype]
	if !ok {
		str := fmt.Sprintf("Get on unregister obj type={%v}", objtype)
		log.Printf(str)
		return nil, errors.New(str)
	}

	if !chain.IsValidate(objtype) {
		str := fmt.Sprintf("Get on missmatched obj type={%v}", objtype)
		log.Printf(str)
		return nil, errors.New(str)
	}

	return chain.Find(key)
}

func (self *DataBase) Undo() {
	for _, v := range self.index_map {
		v.Undo()
	}
}

func (self *DataBase) Squash() {
	for _, v := range self.index_map {
		v.Squash()
	}
}

func (self *DataBase) Commit(revision int64) {
	for _, v := range self.index_map {
		v.Commit(revision)
	}
}

func (self *DataBase) UndoAll() {
	for _, v := range self.index_map {
		v.UndoAll()
	}
}

func (self *DataBase) StartUndoSession(enabled bool) *SessionList {
	var list []*Session
	if enabled {
		for _, v := range self.index_map {
			ite := v.Start_Undo_Session(enabled)
			list = append(list, ite)
		}

		return NewSessionSet(list)
	}

	return NewSessionSet(list)
}

func (self *DataBase) AddIndex(objtype uint32) error {
	if objtype >= max_type {
		str := fmt.Sprintf("AddIndex: obj type={%v} has not already definted", objtype)
		log.Printf(str)
		return errors.New(str)
	}

	//check obj type is already added
	if _, ok := self.index_map[objtype]; ok {
		str := fmt.Sprintf("AddIndex: obj type={%v} already exist", objtype)
		log.Printf(str)
		return errors.New(str)
	}

	in := NewGenericIndex(objtype, self.db)
	self.index_map[objtype] = in

	return nil
}

func (self *DataBase) GetIndex(objtype uint32) *Generic_Index {
	index, ok := self.index_map[objtype]
	if !ok {
		return nil
	}

	return index
}
