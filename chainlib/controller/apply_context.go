package controller

import (
	"datx_chain/chainlib/chainbase"
	"datx_chain/chainlib/types"
	"datx_chain/utils/helper"
	"log"
	"time"
)

//ApplyContent struct
type ApplyContent struct {
	Controller *Controller

	DB *chainbase.DataBase

	TrxContext *TransactionContext

	Receiver string

	ContextFree bool

	Act types.Action

	Trace types.ActionTrace
}

//NewApplyContext new ApplyContent
func NewApplyContext(trx *TransactionContext, act types.Action) *ApplyContent {
	return &ApplyContent{
		Controller:  trx.Controller,
		ContextFree: false,
		Act:         act,
		TrxContext:  trx,
	}
}

//Exec exec theh action
func (ac *ApplyContent) Exec() {
	// trace :=
	ac.execOne()

	// log.Printf("ApplyContext trace: %v", trace)
}

//
func (ac *ApplyContent) execOne() types.ActionTrace {
	start := time.Now()
	helper.CatchException(nil, func() {
		log.Print("ApplyContext execOne panic")
	})

	hand, ok := ac.Controller.ApplyHandlers.Find(ac.Act.ActionName)
	if !ok {
		return types.ActionTrace{}
	}

	hand(ac)

	var r types.ActionReceipt
	r.Receiver = ac.Receiver
	r.ActDigest = helper.RLPHash(ac.Act)

	var res types.ActionTrace
	res.Receipt = r
	res.TrxID = ac.TrxContext.ID
	res.Act = ac.Act

	ac.TrxContext.Executed = append(ac.TrxContext.Executed, r)

	res.Elapsed = int64(time.Now().Sub(start))

	return res
}
