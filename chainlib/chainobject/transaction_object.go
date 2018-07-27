package chainobject

import "datx_chain/utils/common"

//TransactionObject object
type TransactionObject struct {
	ID uint64

	Expiration uint64

	TrxID common.Hash
}

//NewTrxObj new trx obj
func NewTrxObj(expir uint64, id common.Hash) *TransactionObject {
	return &TransactionObject{
		ID:         GetOID(TranscationType),
		Expiration: expir,
		TrxID:      id,
	}
}
