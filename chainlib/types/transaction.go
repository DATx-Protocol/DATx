package types

import (
	"datx_chain/utils/common"
	"datx_chain/utils/helper"
)

type bdata []byte

//TransactionHeader struct
type TransactionHeader struct {
	//transaction expiration ,the seconds from 1970.
	Expiration uint64

	//reference the latest block number
	RefBlockNum uint32

	//reference block prefix
	RefBlockPerfix common.Hash

	DelaySec uint64 //number of seconds to delay this transaction for during which it may be canceled.
}

//Transaction struct
type Transaction struct {
	//transaction header, inherited from TransactionHeader
	TransactionHeader

	//action list
	Actions []Action

	ContextFreeActions []Action

	//transaction hash
	TransactionHash common.Hash
}

//SignedTransaction sign trx
type SignedTransaction struct {
	Transaction

	ContextFreeData []byte

	Signatures []bdata
}

//NewTrx transaction constructor
func NewTrx(time uint64) *Transaction {
	return &Transaction{
		TransactionHeader: TransactionHeader{
			Expiration: time,
			DelaySec:   0,
		},
	}
}

//SetReferenceBlock set reference block
func (t *Transaction) SetReferenceBlock(h BlockHeader) {
	t.RefBlockNum = h.BlockNum
	t.RefBlockPerfix = h.ID
}

//ID calculate th id
func (t *Transaction) ID() common.Hash {
	return helper.RLPHash(t)
}

//TotalActions count
func (t *Transaction) TotalActions() uint32 {
	return uint32(len(t.Actions) + len(t.ContextFreeActions))
}
