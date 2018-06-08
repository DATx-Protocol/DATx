package chain_plugin

import (
	"datx_chain/chainlib/types"
	"datx_chain/utils/common"
	"datx_chain/utils/db"
	"datx_chain/utils/rlp"
	"fmt"
	"os"
	"sync"
)

type BranchType []*types.Block

type ForkDB struct {
	//data dir
	path string

	//levelDB instance
	db *datxdb.LDBDatabase

	//levelDB batch
	batch datxdb.Batch

	//first index map,you can get block data from db by block id
	// first map[common.Hash]interface{}

	//index map,save blockNum/blockID pairs
	index map[uint32]common.Hash

	//the lastest block state
	head *types.BlockHeader

	//read/write lock
	rdlock sync.RWMutex
}

func NewForkDB(path string, cache, handles int) (*ForkDB, error) {
	//check for the directory's existence and create it if it doesn't exist
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	db := new(ForkDB)
	db.index = make(map[uint32]common.Hash)

	db.path = path
	file := path /* + string(os.PathSeparator) + "forkdb.dat" */

	//open levelDB. the db will recover when db crashed or file is existence
	var err error
	if db.db, err = datxdb.NewLDBDatabase(file, cache, handles); err != nil {
		fmt.Printf("NewForkDB open err #{%s}", err)
		return nil, err
	}

	//create db batch for write batched
	db.batch = db.db.NewBatch()

	//unmarshal index map
	iter := db.db.NewIterator()
	defer iter.Release()

	//create index map
	for iter.First(); iter.Valid(); iter.Next() {
		val := iter.Value()
		var block types.Block
		if err := rlp.DecodeBytes(val, &block); err != nil {
			fmt.Printf("NewForkDB Decode block: #{%v}", err)
			return nil, err
		}
		num := block.GetNum()
		id := block.GetID()
		db.index[num] = *id
	}

	//unmarshal last block to head
	if iter.Last() == true && iter.Valid() {
		var block types.Block
		err := rlp.DecodeBytes(iter.Value(), &block)
		if err != nil {
			fmt.Printf("ForkDB::NewForkDB err={%v} ", err)
			return nil, err
		}
		db.head = block.GetHead()
	}

	return db, nil
}

//close the db,but not clear db data
func (self *ForkDB) Close() {
	self.index = nil
	self.head = nil
	self.db.Close()
}

//clear all leveldb data,you don'y call this method if you not need
func (self *ForkDB) ClearAll() {
	iter := self.db.NewIterator()
	defer iter.Release()

	//create index map
	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		self.db.Delete(key)
	}
}

func (self *ForkDB) Add(block *types.Block) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ForkDB::Add exception: %s\n", err)
		}
	}()

	// add lock for update fork db
	self.rdlock.Lock()
	defer self.rdlock.Unlock()

	id := block.GetID()
	key := id.Bytes()

	//if not exist,do nothing
	if _, err := self.db.Get([]byte(key)); err == nil {
		fmt.Printf("ForkDB::Add key={%s} already exist", id.Hex())
		return
	}

	num := block.GetNum()
	data, err := rlp.EncodeToBytes(block)
	if err != nil {
		fmt.Printf("ForkDB::Add encode block num={%d} err={%v}", num, err)
		return
	}

	// fmt.Printf("ForkDB::Add after encode block data={%v}}\n", data)
	// fmt.Printf("ForkDB::Add add block num={%d} id={%s} block={%v}", block.GetNum(), block.GetID().Hex(), block)

	self.index[num] = *id
	self.head = block.GetHead()
	err = self.batch.Put(key, data)
	if err != nil {
		fmt.Printf("ForkDB::Add batch.Put block num={%d} err={%v}", num, err)
		return
	}
	err = self.batch.Write()
	if err != nil {
		fmt.Printf("ForkDB::Add batch.Write block num={%d} err={%v}", num, err)
		return
	}
}

func (self *ForkDB) Delete(id common.Hash) {
	key := id.Bytes()

	//if not exist,do nothing
	data, err := self.db.Get([]byte(key))
	if err != nil {
		return
	}

	var block types.Block
	if err := rlp.DecodeBytes(data, &block); err != nil {
		fmt.Printf("ForkDB::Delete block id={%v} err={%v}\n", key, err)
		return
	}

	num := block.GetNum()
	delete(self.index, num)

	//update head.assign previous block head to the head field when you delete the head.do nothing if you don't delete the head
	if self.head.ID == id {
		pre := block.GetPrevious()
		block := self.GetBlockByNum(pre)
		if block != nil {
			self.head = block.GetHead()
		}
	}

	self.db.Delete(key)
}

func (self *ForkDB) GetBlock(id common.Hash) *types.Block {
	key := id.Bytes()

	//if not exist,return nil
	data, err := self.db.Get([]byte(key))
	if err != nil {
		fmt.Printf("ForkDB::GetBlock block id={%s} err={%v}", id.Hex(), err)
		return nil
	}

	var block types.Block
	if err := rlp.DecodeBytes(data, &block); err != nil {
		fmt.Printf("ForkDB::GetBlock decode block id={%s} err={%v}", id.Hex(), err)
		return nil
	}

	return &block
}

func (self *ForkDB) GetBlockByNum(num uint32) *types.Block {
	//find num existence,return nil if not exist or occured error
	if key, ok := self.index[num]; ok {
		id := key.Bytes()
		bytes, err := self.db.Get(id)
		if err != nil {
			fmt.Printf("ForkDB::GetBlockByNum block id={%s} err={%v}", key.Hex(), err)
			return nil
		}

		var block types.Block
		if err := rlp.DecodeBytes(bytes, block); err != nil {
			fmt.Printf("ForkDB::GetBlockByNum decode block id={%s} err={%v}", key.Hex(), err)
			return nil
		}

		return &block
	} else {
		fmt.Printf("ForkDB::GetBlockByNum index map is empty. block num={%d}", num)
		return nil
	}
}

func (self *ForkDB) GetHead() *types.BlockHeader {
	return self.head
}

/* Given two head blocks, return two branches of the fork graph that end with a common ancestor */
func (self *ForkDB) FetchBranch(first, second common.Hash) (BranchType, BranchType) {
	var one BranchType
	var two BranchType

	first_branch := self.GetBlock(first)
	second_branch := self.GetBlock(second)

	for first_branch.GetNum() > second_branch.GetNum() {
		one = append(one, first_branch)
		first_branch = self.GetBlockByNum(first_branch.GetPrevious())
		if first_branch == nil {
			return nil, nil
		}
	}

	for second_branch.GetNum() > first_branch.GetNum() {
		two = append(two, second_branch)
		second_branch = self.GetBlockByNum(second_branch.GetPrevious())
		if second_branch == nil {
			return nil, nil
		}
	}

	for first_branch.GetPrevious() != second_branch.GetPrevious() {
		one = append(one, first_branch)
		two = append(two, second_branch)

		first_branch = self.GetBlockByNum(first_branch.GetPrevious())
		second_branch = self.GetBlockByNum(second_branch.GetPrevious())

		if first_branch == nil && second_branch == nil {
			return nil, nil
		}
	}

	if first_branch != nil && second_branch != nil {
		one = append(one, first_branch)
		two = append(two, second_branch)
	}

	return one, two
}
