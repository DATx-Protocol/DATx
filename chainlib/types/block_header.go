package types

import (
	"datx_chain/chainlib/chainobject"
	"datx_chain/utils/common"
	"datx_chain/utils/helper"
	"time"
)

//HeaderConfirmation struct
type HeaderConfirmation struct {
	BlockID           common.Hash
	Producer          string //producer name
	ProducerSignature []byte //producer signature
}

//BlockHeader struct
type BlockHeader struct {
	//current block num
	BlockNum uint32

	//current block header hash
	ID common.Hash

	//Previous block id
	Previous common.Hash

	//producer account
	Producer string

	//producer signature
	Signature []byte

	//the root of action merkle tree
	ActionMroot common.Hash

	//the root of transaction merkle tree
	TransactionMroot common.Hash

	//time
	TimeStamp *BlockTime

	ScheduleVersion uint32
	NewProducers    chainobject.ProducerSchedule

	Confirmed uint16
}

//NewBlockHeader new
func NewBlockHeader(num uint32, pre common.Hash, prod string) *BlockHeader {
	return &BlockHeader{
		BlockNum:  num,
		Previous:  pre,
		Producer:  prod,
		TimeStamp: NewBlockTime(time.Now()),
	}
}

//Hash id
func (bh *BlockHeader) Hash() common.Hash {
	if bh.ID != common.HexToHash("") {
		bh.ID = common.HexToHash("")
	}
	bh.ID = helper.RLPHash(&bh)

	return bh.ID
}
