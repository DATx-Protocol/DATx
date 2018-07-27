package controller

import (
	"datx_chain/chainlib/types"
	"datx_chain/utils/common"
	"encoding/binary"
	"fmt"
	"testing"
)

var db *ForkDB

func NewDB() *ForkDB {
	// path := helper.MakePath("fork_db")
	path := "fork_db"

	var err error
	db, _ = NewForkDB(nil, path)
	if err != nil {
		fmt.Print(err)
		return nil
	}

	return db
}

func NewBlock(num, pre uint32, prod string) *types.Block {
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, pre)
	ha := common.BytesToHash(a)

	return &types.Block{
		BlockHeader: *types.NewBlockHeader(num, ha, prod),
	}
}

func Test_Add(t *testing.T) {
	NewDB()
	if db == nil {
		t.Error("test : new db error")
	}

	block := NewBlock(1, 0, "1")
	id := block.Hash()

	BlockState := types.NewBlockState(block)
	BlockState.DposIrreversibleBlockNum = 4
	BlockState.BftIrreversibleBlockNum = 8
	db.AddState(BlockState)

	//add third block state
	block3 := NewBlock(3, 2, "3")
	block3.Hash()
	BlockState3 := types.NewBlockState(block3)
	BlockState3.DposIrreversibleBlockNum = 4
	BlockState3.BftIrreversibleBlockNum = 7
	db.AddState(BlockState3)

	//add second block state
	block2 := NewBlock(2, 1, "2")
	block2.Hash()
	BlockState2 := types.NewBlockState(block2)
	BlockState2.DposIrreversibleBlockNum = 4
	BlockState2.BftIrreversibleBlockNum = 6
	db.AddState(BlockState2)

	//add fourth block state
	block4 := NewBlock(4, 3, "3")
	block4.Hash()
	BlockState4 := types.NewBlockState(block4)
	BlockState4.DposIrreversibleBlockNum = 5
	BlockState4.BftIrreversibleBlockNum = 4
	db.AddState(BlockState4)

	dbblock := db.GetBlock(id)
	if dbblock.GetNum() != uint32(1) {
		t.Error("add: get block num error")
	}

}

func Test_Remove(t *testing.T) {
	db.ClearAll()
	block := NewBlock(1, 0, "1")
	id := block.Hash()

	BlockState := types.NewBlockState(block)
	BlockState.DposIrreversibleBlockNum = 4
	BlockState.BftIrreversibleBlockNum = 8
	db.AddState(BlockState)

	//add second block state
	block2 := NewBlock(2, 1, "2")
	id2 := block2.Hash()
	block2.Previous = id
	BlockState2 := types.NewBlockState(block2)
	BlockState2.DposIrreversibleBlockNum = 4
	BlockState2.BftIrreversibleBlockNum = 6
	db.AddState(BlockState2)

	//add third block state
	block3 := NewBlock(3, 2, "3")
	block3.Hash()
	block3.Previous = id
	BlockState3 := types.NewBlockState(block3)
	BlockState3.DposIrreversibleBlockNum = 4
	BlockState3.BftIrreversibleBlockNum = 7
	db.AddState(BlockState3)

	//add fourth block state
	block4 := NewBlock(4, 3, "3")
	block4.Hash()
	block4.Previous = id2
	BlockState4 := types.NewBlockState(block4)
	BlockState4.DposIrreversibleBlockNum = 5
	BlockState4.BftIrreversibleBlockNum = 4
	db.AddState(BlockState4)

	//remove block id and delete all block that previous dependent on this block id
	db.Remove(id)

	delblock := db.GetBlock(id2)
	if delblock != nil {
		t.Error("delete block failed")
	}
}
