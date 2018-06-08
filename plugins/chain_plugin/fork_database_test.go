package chain_plugin

import (
	"datx_chain/chainlib/types"
	"fmt"
	"os"
	"strings"
	"testing"
)

func NewDB(dir string) *ForkDB {
	// file := utils.GetCurrentPath()
	// ins := strings.Split(file, string(os.PathSeparator))

	// ps := append(ins[:len(ins)-2], "db_test")

	ps := []string{"home", "simon", dir, "db_test"}

	path := strings.Join(ps, string(os.PathSeparator))

	db, err := NewForkDB(path, 0, 0)
	if err != nil {
		fmt.Print(err)
		return nil
	}

	return db
}

func Test_NewForkDB(t *testing.T) {
	db := NewDB("test")
	if db == nil {
		t.Error("test : new db error")
	}
	defer db.Close()
}

func NewBlock(num, pre uint32, prod string) *types.Block {
	return &types.Block{
		BlockHeader: *types.MakeBlockHeader(num, pre, prod),
	}
}

func Test_Add(t *testing.T) {
	block := NewBlock(1, 0, "1")
	id := block.Hash()

	db := NewDB("add")
	if db == nil {
		t.Error("add: new db error")
	}
	defer db.Close()

	db.Add(block)

	db_block := db.GetBlock(id)
	if db_block.GetNum() != uint32(1) {
		t.Error("add: get block num error")
	}

}

func Test_Delete(t *testing.T) {
	block := NewBlock(100, 99, "delete")
	id := block.Hash()

	db := NewDB("add")
	if db == nil {
		t.Error("add: new db error")
	}
	defer db.Close()
	db.Add(block)

	//delete
	db.Delete(id)

	del_block := db.GetBlock(id)
	if del_block != nil {
		t.Error("delete block failed")
	}

	//delete id where not exist
	block_new := NewBlock(1001, 1000, "delete")
	id_new := block_new.Hash()
	db.Delete(id_new)

	del_block_new := db.GetBlock(id_new)
	if del_block_new != nil {
		t.Error("delete block where not exist occured failed")
	}

}
