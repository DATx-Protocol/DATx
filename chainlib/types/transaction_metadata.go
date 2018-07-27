package types

import (
	"datx_chain/utils/common"
	"datx_chain/utils/helper"
)

//TransactionMetaData wrap Transaction struct
type TransactionMetaData struct {
	ID common.Hash

	SignedID common.Hash

	Trx SignedTransaction

	PackedTrx PackedTransaction

	Accepted bool
}

//NewTrxMetaData new TransactionMetaData
func NewTrxMetaData(trx *SignedTransaction) *TransactionMetaData {
	var res TransactionMetaData
	res.Trx = *trx
	res.ID = trx.ID()
	res.PackedTrx = *NewPackedTransaction(trx, None)
	res.SignedID = helper.RLPHash(res.PackedTrx)
	res.Accepted = false

	return &res
}

//NewMetaDataByPackedTrx new
func NewMetaDataByPackedTrx(pt *PackedTransaction) *TransactionMetaData {
	var res TransactionMetaData
	res.Trx = pt.GetSignTransaction()
	res.ID = pt.ID()
	res.PackedTrx = *pt
	res.SignedID = helper.RLPHash(res.PackedTrx)
	res.Accepted = false

	return &res
}
