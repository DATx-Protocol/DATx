package types

import "datx_chain/utils/common"

//BaseActionTrace struct
type BaseActionTrace struct {
	Receipt ActionReceipt

	Act Action

	Elapsed int64 //microseconds

	TrxID common.Hash //the transaction id that generated this action
}

//ActionTrace struct
type ActionTrace struct {
	BaseActionTrace
}

//TransactionTrace struct
type TransactionTrace struct {
	ID common.Hash

	Receipt TransactionReceiptHeader

	Elapsed  int64 //elapse microseconds
	NetUsage uint64

	Scheduled bool

	ActionTraces []ActionTrace

	FailedTrace *TransactionTrace

	Except error //catch exception
}

//BlockTrace struct
type BlockTrace struct {
	Elapsed int64 //microseconds

	TrxTraces []TransactionTrace
}

//TrxTrace pairs
type TrxTrace struct {
	Err error
	Trx *PackedTransaction
}

//AsyncTrx struct
type AsyncTrx struct {
	Pack     *PackedTransaction
	Callback func(inerr error, trace *TransactionTrace)
}
