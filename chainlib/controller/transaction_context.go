package controller

import (
	"datx_chain/chainlib/chainbase"
	"datx_chain/chainlib/chainobject"
	"datx_chain/chainlib/types"
	"datx_chain/utils/common"
	"time"
)

//TransactionContext struct
type TransactionContext struct {
	Controller *Controller

	Trx *types.SignedTransaction

	ID common.Hash

	UndoSession chainbase.SessionList

	TrxTrace *types.TransactionTrace

	Start     time.Time
	Published time.Time

	Executed []types.ActionReceipt

	BillToAccounts map[string]struct{}

	Delay               uint64 //microseconds
	IsInput             bool   //default is false
	ApplyContextFree    bool   //default is true
	CanSubjectivelyFail bool   //default is true

	DeadLine        time.Time //default is maxnum time
	LeeWay          int64     //default 3000
	BilledCPUTimeUS int64     //default 0

	isInit bool //default is false
}

//NewTransactionContext New TransactionContext
func NewTransactionContext(c *Controller, db *chainbase.DataBase, t *types.SignedTransaction, trxID common.Hash) *TransactionContext {
	var res TransactionContext
	res.Controller = c
	res.Trx = t
	res.ID = trxID
	res.UndoSession = db.StartUndoSession(true)
	res.TrxTrace = &types.TransactionTrace{}
	res.Start = time.Now()

	//default
	res.DeadLine = types.MaxTime() //the seconds of max time
	res.IsInput = false
	res.ApplyContextFree = true
	res.CanSubjectivelyFail = true
	res.LeeWay = 3000
	res.isInit = false
	res.BilledCPUTimeUS = 0

	res.TrxTrace.ID = res.ID
	res.Executed = make([]types.ActionReceipt, t.TotalActions())

	return &res
}

func (tc *TransactionContext) init(initNetUsage uint64) {

	//record accounts to be billed for network and CPU usage
	for _, v := range tc.Trx.Actions {
		for _, va := range v.Authorization {
			tc.BillToAccounts[va.Actor] = struct{}{}
		}
	}

	tc.isInit = true
}

//InitForImplicitTrx init implicit trx
func (tc *TransactionContext) InitForImplicitTrx(publishTime time.Time) {
	tc.Published = publishTime
	tc.init(0)
}

//InitForInputTrx init input trx
func (tc *TransactionContext) InitForInputTrx(db *chainbase.DataBase, publishTime time.Time) {
	tc.IsInput = true
	tc.Published = publishTime

	//
	tc.init(0)

	tc.recordTransaction(db, tc.ID, tc.Trx.Expiration)
}

func (tc *TransactionContext) recordTransaction(db *chainbase.DataBase, trxID common.Hash, expire uint64) error {
	obj := chainobject.NewTrxObj(expire, trxID)

	return db.Create(chainobject.TranscationType, obj)
}

//Exec method
func (tc *TransactionContext) Exec() {
	if !tc.isInit {
		return
	}

	if tc.ApplyContextFree {
		for _, v := range tc.Trx.ContextFreeActions {

			var trace types.ActionTrace
			tc.DispatchAction(trace, v, v.Account, true, 0)

			tc.TrxTrace.ActionTraces = append(tc.TrxTrace.ActionTraces, trace)

		}
	}

	if tc.Delay == 0 {
		for _, v := range tc.Trx.Actions {
			var trace types.ActionTrace
			tc.DispatchAction(trace, v, v.Account, false, 0)

			tc.TrxTrace.ActionTraces = append(tc.TrxTrace.ActionTraces, trace)
		}
	} else {
		tc.ScheduleTransaction()
	}
}

//Finalize Finalize
func (tc *TransactionContext) Finalize() {

}

//Squash session Squash
func (tc *TransactionContext) Squash() {
	tc.UndoSession.Squash()
}

//DispatchAction  Dispatch the Action
func (tc *TransactionContext) DispatchAction(trace types.ActionTrace, act types.Action, receiver string, contxtfree bool, depth uint32) {
	applyContext := NewApplyContext(tc, act)

	applyContext.ContextFree = contxtfree
	applyContext.Receiver = receiver

	defer func() {
		if err := recover(); err != nil {
			trace = applyContext.Trace

			panic(err)
		}
	}()

	applyContext.Exec()
	trace = applyContext.Trace

}

//ScheduleTransaction schedule
func (tc *TransactionContext) ScheduleTransaction() {

}
