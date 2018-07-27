package controller

import (
	"datx_chain/chainlib/types"
	"datx_chain/utils/common"
	"log"
	"math/big"
	"os"
	"sort"
	"time"
)

//BranchType new type
type BranchType []*types.BlockState

type byBlockNum struct {
	id common.Hash

	num            uint32 //blockstate block num
	inCurrentChain bool   //blockstate in_current_chain
}

type byPrev struct {
	id common.Hash //current block id

	prev common.Hash //previous block id
}

type byLibBlockNum struct {
	id common.Hash

	dposIrreversible uint32 // blockstate dpos_irreversible_block_num

	bftIrreversible uint32 //blockstate bfp_irreversible_block_num

	num uint32 //blockstate block num
}

//ForkDB struct
type ForkDB struct {
	//data dir
	path string

	controller *Controller

	//index map,save blockID/blockstate pairs
	index map[common.Hash]types.BlockState

	//the lastest block state
	head *types.BlockState

	lprev map[common.Hash][]*types.BlockState //list of byPrev

	lnum []byBlockNum //list of byBlockNum

	llib []byLibBlockNum //list of byLibBlockNum
}

//NewForkDB create
func NewForkDB(chain *Controller, path string) (*ForkDB, error) {
	//check for the directory's existence and create it if it doesn't exist
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	db := new(ForkDB)
	db.index = make(map[common.Hash]types.BlockState)

	db.controller = chain
	db.path = path

	db.lprev = make(map[common.Hash][]*types.BlockState)
	db.lnum = make([]byBlockNum, 0)
	db.llib = make([]byLibBlockNum, 0)

	return db, nil
}

//Close clear all data
func (fdb *ForkDB) Close() {
	fdb.ClearAll()
}

//ClearAll all data
func (fdb *ForkDB) ClearAll() {
	fdb.index = make(map[common.Hash]types.BlockState)
	fdb.lprev = make(map[common.Hash][]*types.BlockState)
	fdb.lnum = nil
	fdb.llib = nil
}

func (fdb *ForkDB) insert(b *types.BlockState) bool {
	if b == nil {
		return false
	}

	if _, ok := fdb.index[b.ID]; ok {
		return false
	}

	fdb.index[b.ID] = *b

	//update lprev
	prev := b.GetPrevious()
	if data, ok := fdb.lprev[prev]; ok {
		data = append(data, b)
		fdb.lprev[prev] = data
	} else {
		newData := make([]*types.BlockState, 0)
		newData = append(newData, b)
		fdb.lprev[prev] = newData
	}

	//update lnum
	var newlnum byBlockNum
	newlnum.id = b.ID
	newlnum.num = b.BlockNum
	newlnum.inCurrentChain = b.InChain

	fdb.lnum = append(fdb.lnum, newlnum)

	//update llib
	var newllib byLibBlockNum
	newllib.id = b.ID
	newllib.dposIrreversible = b.DposIrreversibleBlockNum
	newllib.bftIrreversible = b.BftIrreversibleBlockNum
	newllib.num = b.BlockNum

	fdb.llib = append(fdb.llib, newllib)

	return true
}

func (fdb *ForkDB) delete(id common.Hash) {

	delete(fdb.index, id)
	delete(fdb.lprev, id)

	for i, v := range fdb.llib {
		if id == v.id {
			fdb.llib = append(fdb.llib[:i], fdb.llib[i+1:]...)
		}
	}

	for i, v := range fdb.lnum {
		if id == v.id {
			fdb.lnum = append(fdb.lnum[:i], fdb.lnum[i+1:]...)
		}
	}
}

func (fdb *ForkDB) sortupdate() {
	//
	sort.Slice(fdb.lnum, func(i, j int) bool {
		//sort by block num asc
		if fdb.lnum[i].num < fdb.lnum[j].num {
			return true
		}
		if fdb.lnum[i].num > fdb.lnum[j].num {
			return false
		}
		return fdb.lnum[i].inCurrentChain
	})

	//
	sort.Slice(fdb.llib, func(i, j int) bool {
		//sort by dpos irreversible block num desc
		if fdb.llib[i].dposIrreversible > fdb.llib[j].dposIrreversible {
			return true
		}
		if fdb.llib[i].dposIrreversible < fdb.llib[j].dposIrreversible {
			return false
		}

		//sort by bft irreversible block num desc
		if fdb.llib[i].bftIrreversible > fdb.llib[j].bftIrreversible {
			return true
		}
		if fdb.llib[i].bftIrreversible < fdb.llib[j].bftIrreversible {
			return false
		}

		return fdb.llib[i].num > fdb.llib[j].num
	})
}

//Set method
func (fdb *ForkDB) Set(s *types.BlockState) {
	issuc := fdb.insert(s)
	if !issuc { //insert failed,duplicate state detected
		return
	}

	if s.ID != s.Header.ID {
		return
	}

	if fdb.head != nil {
		fdb.head = s
	} else if fdb.head.BlockNum < s.BlockNum {
		fdb.head = s
	}
}

func (fdb *ForkDB) prune(b *types.BlockState) {
	num := b.BlockNum

	fdb.sortupdate()

	for _, v := range fdb.lnum {
		if v.num < num {
			if data, ok := fdb.index[v.id]; ok {
				fdb.controller.OnIrreversible(&data)
				fdb.delete(v.id)
			}
		}
	}
}

//AddState method
func (fdb *ForkDB) AddState(b *types.BlockState) *types.BlockState {
	if !fdb.insert(b) { //duplicate block added
		return nil
	}

	//sort
	fdb.sortupdate()

	//update head
	firstid := fdb.llib[0].id
	var hbs types.BlockState
	hbs, _ = fdb.index[firstid]
	fdb.head = &hbs

	sb := fdb.head.Header.TimeStamp.Time
	log.Printf("###### AddTime: %v  %v\n", sb, big.NewInt(time.Now().Unix()))

	log.Printf("*********  AddState Head: %v  %v  %v\n", fdb.head.BlockNum, fdb.head.Block.Producer, fdb.head.Block.ID.String())

	//delete old block
	lib := fdb.head.DposIrreversibleBlockNum
	oldid := fdb.lnum[0].id
	tbs, _ := fdb.index[oldid]
	oldest := &tbs

	if oldest.BlockNum < lib {
		fdb.prune(oldest)
	}

	return b
}

//AddBlock method
func (fdb *ForkDB) AddBlock(b *types.Block, trust bool) *types.BlockState {
	if b == nil || fdb.head == nil {
		return nil
	}

	prev, ok := fdb.index[b.Previous]
	if !ok {
		return nil
	}

	log.Printf("\n****AddBlock: %v  %v  %v\n", b.BlockNum, b.Producer, b.ID.String())
	resul := types.NewBlockStateByHeadState(prev.BlockHeaderState, b, trust)
	return fdb.AddState(resul)
}

//AddConfirmation add confir to block
func (fdb *ForkDB) AddConfirmation(c *types.HeaderConfirmation) {
	b := fdb.GetBlock(c.BlockID)
	if b == nil {
		return
	}

	b.AddConfirm(c)

	suc := (len(b.ActiveSchedule.Producers) * 2) / 3
	if b.BftIrreversibleBlockNum < b.BlockNum && len(b.Confirmations) > suc {
		fdb.SetBftIrreversible(c.BlockID)
	}
}

//Head return head
func (fdb *ForkDB) Head() *types.BlockState {
	return fdb.head
}

//Index return index
func (fdb *ForkDB) GetIndex() map[common.Hash]types.BlockState {
	return fdb.index
}

//Remove delete data
func (fdb *ForkDB) Remove(id common.Hash) {
	var removeQueue []common.Hash
	removeQueue = append(removeQueue, id)

	for i := 0; i < len(removeQueue); i++ {
		v := removeQueue[i]

		data, _ := fdb.lprev[v]
		for _, vl := range data {
			removeQueue = append(removeQueue, vl.ID)
		}

		fdb.delete(v)
	}

	fdb.sortupdate()
	if len(fdb.llib) > 0 {
		headid := fdb.llib[0].id

		hbs, _ := fdb.index[headid]
		fdb.head = &hbs
	}

	fdb.head = nil
}

//SetValidity set block validated
func (fdb *ForkDB) SetValidity(s *types.BlockState, vali bool) {
	if !vali {
		fdb.Remove(s.ID)
	} else {
		s.Validated = true
	}
}

//MarkInChain make block in chain
func (fdb *ForkDB) MarkInChain(s *types.BlockState, inChain bool) {
	if s.InChain == inChain {
		return
	}

	//modify attribute of inChain
	if data, ok := fdb.index[s.ID]; ok {
		data.InChain = inChain
	}
}

//GetBlock get block by block id
func (fdb *ForkDB) GetBlock(id common.Hash) *types.BlockState {
	if data, ok := fdb.index[id]; ok {
		return &data
	}

	return nil
}

//GetTrx get trabsaction by transaction id and timestamp
func (fdb *ForkDB) GetTrx(id common.Hash) (*types.Transaction, *types.BlockTime) {
	for _, v := range fdb.index {
		for _, j := range v.Trxs {
			if j.ID == id {
				return &j.Trx.Transaction, v.Block.TimeStamp

			}
			return nil, nil
		}
	}
	return nil, nil
}

//GetBlockInChain get block by block num in chain
func (fdb *ForkDB) GetBlockInChain(num uint32) *types.BlockState {
	for _, v := range fdb.lnum {
		if v.num == num {
			id := v.id

			if data, ok := fdb.index[id]; ok {
				return &data
			}

			return nil
		}
	}

	return nil

}

//SetBftIrreversible set bft irreversible block id
func (fdb *ForkDB) SetBftIrreversible(id common.Hash) {
	if data, ok := fdb.index[id]; ok {
		num := data.BlockNum
		data.BftIrreversibleBlockNum = num

		//update : modify
		fdb.index[id] = data
		fdb.sortupdate()

		var updateQueue []common.Hash
		updateQueue = append(updateQueue, id)

		for i := 0; i < len(updateQueue); i++ {
			v := updateQueue[i]

			if prev, ok := fdb.lprev[v]; ok {
				for j := 0; j < len(prev); j++ {
					if prev[j].BftIrreversibleBlockNum < num {
						prev[j].BftIrreversibleBlockNum = num
						updateQueue = append(updateQueue, prev[j].ID)
					}
				}
			} //lprev
		} //updateQueue
	}
}

// FetchBranch Given two head blocks, return two branches of the fork graph that end with a common ancestor
func (fdb *ForkDB) FetchBranch(first, second common.Hash) (BranchType, BranchType) {
	var one BranchType
	var two BranchType

	firstBranch := fdb.GetBlock(first)
	secondBranch := fdb.GetBlock(second)

	for firstBranch.GetNum() > secondBranch.GetNum() {
		one = append(one, firstBranch)
		firstBranch = fdb.GetBlock(firstBranch.GetPrevious())
		if firstBranch == nil {
			return nil, nil
		}
	}

	for secondBranch.GetNum() > firstBranch.GetNum() {
		two = append(two, secondBranch)
		secondBranch = fdb.GetBlock(secondBranch.GetPrevious())
		if secondBranch == nil {
			return nil, nil
		}
	}

	for firstBranch.GetPrevious() != secondBranch.GetPrevious() {
		one = append(one, firstBranch)
		two = append(two, secondBranch)

		firstBranch = fdb.GetBlock(firstBranch.GetPrevious())
		secondBranch = fdb.GetBlock(secondBranch.GetPrevious())

		if firstBranch == nil && secondBranch == nil {
			return nil, nil
		}
	}

	if firstBranch != nil && secondBranch != nil {
		one = append(one, firstBranch)
		two = append(two, secondBranch)
	}

	return one, two
}
