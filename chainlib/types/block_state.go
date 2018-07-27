package types

import (
	"datx_chain/chainlib/chainobject"
	"datx_chain/utils/common"
	"math/big"
)

//BlockState block state
type BlockState struct {
	BlockHeaderState
	Block     *Block
	Validated bool
	InChain   bool

	Trxs []*TransactionMetaData
}

func NewBlockStateEmpty() *BlockState {
	var res BlockState

	res.BlockHeaderState = BlockHeaderState{
		Header: BlockHeader{
			Signature:    make([]byte, 0),
			TimeStamp:    &BlockTime{Time: big.NewInt(0)},
			NewProducers: chainobject.ProducerSchedule{Producers: make([]chainobject.ProducerKey, 0)},
		},
		PendingSchedule:          chainobject.ProducerSchedule{Producers: make([]chainobject.ProducerKey, 0)},
		ActiveSchedule:           chainobject.ProducerSchedule{Producers: make([]chainobject.ProducerKey, 0)},
		ProducerToLastImpliedIrb: make(map[string]interface{}, 0),
		ProducerToLastProduced:   make(map[string]interface{}, 0),
		Count:         make([]uint8, 0),
		Confirmations: make([]HeaderConfirmation, 0),
	}

	res.Block = &Block{
		BlockHeader: BlockHeader{
			Signature:    make([]byte, 0),
			TimeStamp:    &BlockTime{Time: big.NewInt(0)},
			NewProducers: chainobject.ProducerSchedule{Producers: make([]chainobject.ProducerKey, 0)},
		},
		Transactions: make([]*TransactionReceipt, 0),
	}
	res.Trxs = make([]*TransactionMetaData, 0)

	return &res
}

//NewBlockState new
func NewBlockState(b *Block) *BlockState {
	return &BlockState{
		BlockHeaderState: NewBlockHeaderState(*b),
		Block:            b,
		Validated:        false,
		InChain:          false,
	}
}

//NewBlockStateByHeader new
func NewBlockStateByHeader(s BlockHeaderState) *BlockState {
	return &BlockState{
		BlockHeaderState: s,
		Block:            nil,
		Validated:        false,
		InChain:          false,
	}
}

//NewBlockStateByTime new
func NewBlockStateByTime(prev BlockHeaderState, when *BlockTime) *BlockState {
	var result BlockState
	next := prev.GenerateNext(when)

	if next != nil {
		result.BlockHeaderState = *next
		result.Block = &Block{BlockHeader: result.BlockHeaderState.Header}
	} else {
		result.BlockHeaderState = BlockHeaderState{}
		result.Block = &Block{}
	}

	result.Validated = false
	result.InChain = false

	return &result
}

//NewBlockStateByHeadState new
func NewBlockStateByHeadState(prev BlockHeaderState, b *Block, trust bool) *BlockState {
	var result BlockState
	next := prev.Next(&b.BlockHeader, trust)

	if next != nil {
		result.BlockHeaderState = *next
		result.Block = b
	} else {
		result.BlockHeaderState = BlockHeaderState{}
		result.Block = b
	}

	result.Validated = false
	result.InChain = false

	return &result
}

//GetID id
func (bs *BlockState) GetID() *common.Hash {
	return &bs.Block.ID
}

//GetNum block num
func (bs *BlockState) GetNum() uint32 {
	return bs.Block.BlockNum
}

//GetHead block head
func (bs *BlockState) GetHead() *BlockHeader {
	return &bs.Block.BlockHeader
}

//GetPrevious Previous
func (bs *BlockState) GetPrevious() common.Hash {
	return bs.Block.Previous
}
