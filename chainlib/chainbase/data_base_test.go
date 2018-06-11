package chainbase

import (
	"datx_chain/chainlib/chainobject"
	"datx_chain/utils/helper"
	"strings"
	"testing"
)

var db *DataBase

func Test_AddIndex(t *testing.T) {
	vur := helper.MakePath("chain_test")

	db = NewDataBase(vur)

	//add index of account type
	if err := db.AddIndex(account_type); err != nil {
		t.Errorf("AddIndex obj type={%v} err={%v}", account_type, err)
	}

	//add index of unregister type
	if err := db.AddIndex(100); err != nil {
		t.Logf("AddIndex obj type={%v} err={%v}", 100, err)
	}
}

func Test_Create(t *testing.T) {
	//add index of account type
	db.AddIndex(account_type)

	//start undo session
	db.StartUndoSession(true)

	//make account object
	acc := chainobject.NewAccount("test")

	//insert account obj into db
	if err := db.Create(account_type, acc); err != nil {
		t.Errorf("insert and get obj type={%v} err={%v}", account_type, err)
	}

	//insert account obj into db whit unregistered obj type
	if err := db.Create(99, acc); err != nil {
		t.Logf("insert and get obj type={%v} err={%v}", 99, err)
	}

	//insert nil obj into db
	if err := db.Create(account_type, nil); err != nil {
		t.Logf("insert and get nil obj  err={%v}", err)
	}
}

func Test_Remove(t *testing.T) {
	//start new undo session
	db.StartUndoSession(true)

	//make account object
	acc := chainobject.NewAccount("remove")

	//insert account obj into db
	db.Create(account_type, acc)

	//remove from db
	if err := db.Remove(account_type, acc); err != nil {
		t.Errorf("remove and get obj err={%v} name={%v} failed", err, acc.Name)
	}

	//remove from db repetition
	if err := db.Remove(account_type, acc); err == nil {
		t.Errorf("remove and get obj err={%v} name={%v} failed", err, acc.Name)
	}

	// remove obj without insert
	acc_1 := chainobject.NewAccount("remove_1")
	if err := db.Remove(account_type, acc_1); err == nil {
		t.Errorf("remove and get obj err={%v} name={%v} failed", err, acc_1.Name)
	}

	//remove nil obj
	if err := db.Remove(account_type, nil); err == nil {
		t.Errorf("remove nil obj err={%v}", err)
	}
}

func Test_Modify(t *testing.T) {
	//start new undo session
	db.StartUndoSession(true)

	//make account object
	acc := chainobject.NewAccount("modify")

	//insert account obj into db
	db.Create(account_type, acc)

	//modify with new value
	acc.Name = "new_modify"
	key := acc.ID()
	if err := db.Modify(account_type, acc); err != nil {
		t.Errorf("Modify err={%v}", err)
	}

	if data, err := db.Get(account_type, key); err != nil {
		t.Errorf("Modify get err={%v}", err)
	} else {
		result := data.(*chainobject.AccountObject)
		if !strings.EqualFold(result.Name, "new_modify") {
			t.Error("Modify checked failed ")
		}
	}

	// modify nil obj
	if err := db.Modify(account_type, nil); err != nil {
		t.Logf("Modify nil obj err={%v}", err)
	}

	// modify unreg obj type
	if err := db.Modify(99, acc); err != nil {
		t.Logf("Modify unreg obj type err={%v}", err)
	}
}

//test when single session
func Test_UndoCreate(t *testing.T) {
	//start undo session
	db.StartUndoSession(true)

	//make account object
	name := "create"
	acc := chainobject.NewAccount(name)
	id := acc.ID()

	//insert account obj into db
	db.Create(account_type, acc)

	//check atate
	if data, err := db.Get(account_type, id); err != nil {
		t.Error("undo create checked failed when get obj before undo")
	} else {
		res := data.(*chainobject.AccountObject)
		if !strings.EqualFold(res.Name, name) {
			t.Error("undo create checked failed when parse data")
		}
	}

	//undo
	db.Undo()

	//check atate
	if _, err := db.Get(account_type, id); err == nil {
		t.Error("undo create checked failed when get obj after undo")
	} else {
		t.Logf("undo create checked success. err={%v}", err)
	}
}

//test when single session
func Test_UndoModify(t *testing.T) {
	//start undo session
	db.StartUndoSession(true)

	//make account object
	name := "undo"
	acc := chainobject.NewAccount(name)
	id := acc.ID()

	//insert account obj into db
	db.Create(account_type, acc)

	//start new undo session
	db.StartUndoSession(true)

	//modify new name
	acc.Name = "modify_undo"

	db.Modify(account_type, acc)

	//undo
	db.Undo()

	//check atate
	if data, err := db.Get(account_type, id); err != nil {
		t.Error("undo Modify checked failed ")
	} else {
		res := data.(*chainobject.AccountObject)
		if !strings.EqualFold(res.Name, name) {
			t.Error("undo checked failed ")
		}
	}
}

//test when single session
func Test_UndoRemove(t *testing.T) {
	//start undo session
	db.StartUndoSession(true)

	//make account object
	name := "undo"
	acc := chainobject.NewAccount(name)
	id := acc.ID()

	//insert account obj into db
	db.Create(account_type, acc)

	//start new undo session
	db.StartUndoSession(true)

	//remove the obj
	db.Remove(account_type, acc)

	//check atate
	if _, err := db.Get(account_type, id); err == nil {
		t.Error("remove checked failed")
	} else {
		t.Logf("remove checked success. err={%v} ", err)
	}

	//undo
	db.Undo()

	//check atate
	if undodata, err := db.Get(account_type, id); err != nil {
		t.Logf("undo checked failed err={%v} ", err)
	} else {
		undores := undodata.(*chainobject.AccountObject)
		if !strings.EqualFold(undores.Name, name) {
			t.Error("undo checked failed ")
		}
	}
}
