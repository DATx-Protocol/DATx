package types

import (
	"datx_chain/utils/common"
	"datx_chain/utils/helper"
)

type BlockHeader struct {
	//current block num
	BlockNum uint32

	//current block header hash
	ID common.Hash

	//Previous block id
	Previous uint32

	//producer account
	Producer string

	//producer signature
	Signature string
}

type Block struct {
	//block header
	BlockHeader

	//transcation pool
	Transcations []*Transcation
}

func MakeBlockHeader(num, pre uint32, prod string) *BlockHeader {
	return &BlockHeader{
		BlockNum: num,
		Previous: pre,
		Producer: prod,
	}
}

func (self *BlockHeader) Hash() common.Hash {
	self.ID = helper.RLPHash(&self)

	return self.ID
}

func MakeBlock(num, pre uint32, prod string) *Block {
	return &Block{
		BlockHeader: BlockHeader{
			BlockNum: num,
			Previous: pre,
			Producer: prod,
		},
	}
}

func (self *Block) GetID() *common.Hash {
	return &self.ID
}

func (self *Block) GetNum() uint32 {
	return self.BlockNum
}

func (self *Block) GetPrevious() uint32 {
	return self.Previous
}

func (self *Block) GetHead() *BlockHeader {
	return &self.BlockHeader
}
