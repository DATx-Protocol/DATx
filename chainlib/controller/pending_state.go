package controller

import (
	"datx_chain/chainlib/chainbase"
	"datx_chain/chainlib/types"
)

//PendingState state of pending block and transaction
type PendingState struct {
	DBSession chainbase.SessionList

	PendingBlockState *types.BlockState

	Actions []types.ActionReceipt

	BlockStatus uint16
}

//NewPendingSate create
func NewPendingSate() *PendingState {
	return &PendingState{}
}

//Push method
func (ps *PendingState) Push() {
	ps.DBSession.Push()
}

//Reset clear trx for reuse
func (ps *PendingState) Reset() {
	ps.PendingBlockState.Block.Transactions = make([]*types.TransactionReceipt, 0)
	ps.PendingBlockState.Trxs = make([]*types.TransactionMetaData, 0)
}
