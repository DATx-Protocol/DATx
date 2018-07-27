package chainbase

import (
	"datx_chain/chainlib/chainobject"
	"datx_chain/utils/db"
	"datx_chain/utils/helper"
	"errors"
	"fmt"
	"log"
	"reflect"
)

//DataBaser interface
type DataBaser interface {
	Insert(v interface{})
	Modify(v interface{})
	Remove(v interface{})
	Find(key uint64) (interface{}, error)
}

//UndoStater interface
type UndoStater interface {
	Undo() error
	Squash() error
	Commit(revision int64)
	UndoAll() error
}

//DataBase used to state rollback
type DataBase struct {
	//the pairs of object type and generic index instance
	index_map map[uint32]*Generic_Index

	//data dir
	path string

	db *datxdb.LDBDatabase
}

//NewDataBase new
func NewDataBase(path string) (*DataBase, error) {
	var result DataBase
	result.path = path
	result.index_map = make(map[uint32]*Generic_Index)
	var err error
	result.db, err = datxdb.NewLDBDatabase(path, 0, 0)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

//Create and insert object
func (sdb *DataBase) Create(objtype uint32, obj interface{}) error {
	chain, ok := sdb.index_map[objtype]
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

	if (reflect.TypeOf(obj)).Kind() != reflect.Ptr {
		str := fmt.Sprintf("Create on pointer to obj type={%v} not on this object", objtype)
		log.Printf(str)
		return errors.New(str)
	}

	return chain.Insert(obj)
}

//Remove method
func (sdb *DataBase) Remove(objtype uint32, obj interface{}) error {
	chain, ok := sdb.index_map[objtype]
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

//Modify the newvalue is the pointer to raw object
func (sdb *DataBase) Modify(objtype uint32, newvalue interface{}) error {
	chain, ok := sdb.index_map[objtype]
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

//Get return value
func (sdb *DataBase) Get(objtype uint32, key uint64) (interface{}, error) {
	chain, ok := sdb.index_map[objtype]
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

//GetBaseValue return value witch the key is a number of uint32
func (sdb *DataBase) GetBaseValue(objtype uint32, key uint32) (interface{}, error) {
	chain, ok := sdb.index_map[objtype]
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

	//encode two uint32 to one uint64
	newid := helper.EncodeBit(objtype, key)

	return chain.Find(newid)
}

//Undo undo session
func (sdb *DataBase) Undo() {
	for _, v := range sdb.index_map {
		v.Undo()
	}
}

//Squash session
func (sdb *DataBase) Squash() {
	for _, v := range sdb.index_map {
		v.Squash()
	}
}

//Commit session
func (sdb *DataBase) Commit(revision int64) {
	for _, v := range sdb.index_map {
		v.Commit(revision)
	}
}

//UndoAll session
func (sdb *DataBase) UndoAll() {
	for _, v := range sdb.index_map {
		v.UndoAll()
	}
}

//StartUndoSession start new session
func (sdb *DataBase) StartUndoSession(enabled bool) SessionList {
	var list []*Session
	if enabled {
		for _, v := range sdb.index_map {
			ite := v.Start_Undo_Session(enabled)
			list = append(list, ite)
		}

		return NewSessionSet(list)
	}

	return NewSessionSet(list)
}

//SetRevision set block num as the revision
func (sdb *DataBase) SetRevision(revision int64) {
	for _, v := range sdb.index_map {
		v.SetRevision(revision)
	}
}

//AddIndex add index for the obj type
func (sdb *DataBase) AddIndex(objtype uint32) error {
	if objtype >= chainobject.MaxType {
		str := fmt.Sprintf("AddIndex: obj type={%v} has not already definted", objtype)
		log.Printf(str)
		return errors.New(str)
	}

	//check obj type is already added
	if _, ok := sdb.index_map[objtype]; ok {
		str := fmt.Sprintf("AddIndex: obj type={%v} already exist", objtype)
		log.Printf(str)
		return errors.New(str)
	}

	in := NewGenericIndex(objtype, sdb.db)
	sdb.index_map[objtype] = in

	return nil
}

//GetIndex get index of the obj type
func (sdb *DataBase) GetIndex(objtype uint32) *Generic_Index {
	index, ok := sdb.index_map[objtype]
	if !ok {
		return nil
	}

	return index
}

//Revision return revision
func (sdb *DataBase) Revision() int64 {
	if len(sdb.index_map) == 0 {
		return -1
	}

	k, ok := sdb.index_map[chainobject.AccountType]
	if !ok {
		return -1
	}

	return k.GetRevision()
}
