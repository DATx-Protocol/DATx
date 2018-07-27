package controller

import (
	"log"
	"time"

	"datx_chain/chainlib/chainbase"
	"datx_chain/chainlib/chainobject"
	"datx_chain/chainlib/types"
	"datx_chain/utils/common"
	"datx_chain/utils/helper"
	"datx_chain/utils/rlp"
)

//CtlConfig config struct
type CtlConfig struct {
	//block log dir
	BlockLogDir string `yaml:"block_log_dir"`

	//chain base dir
	ChainDir string `yaml:"chain_dir"`

	//
	ForceAllChecks bool `yaml:"force_all_checks"`

	//
	ReadOnly bool `yaml:"read_only"`

	//db handles of open file capacity
	Handles int `yaml:"handles"`

	//db can cache block capacity
	Cache int `yaml:"cache"`

	//vm type
	VMType string `yaml:"vm_type"`
}

//Controller struct
type Controller struct {
	Config CtlConfig

	//fork db
	ForkDB *ForkDB

	//block log
	BlockLog *Blog

	//db used to save state of chain
	DB *chainbase.DataBase

	//default config
	Genesis *types.GenesisState

	//accept block chan
	AcceptBlockChan chan *types.BlockState

	//accept block chan
	AcceptBlockHeaderChan chan *types.BlockState

	IrreversibleBlockChan chan *types.BlockState

	//accept transcation chan
	AcceptTrxChan chan *types.TransactionMetaData

	AppliedTransactionChan chan *types.TransactionTrace

	AcceptedConfirmationChan chan *types.HeaderConfirmation

	//sync block chan
	SyncBlockChan chan *types.Block

	AsyncPackedTrxChan chan *types.AsyncTrx

	TransactionAckChan chan *types.TrxTrace

	Head *types.BlockState

	Pending *PendingState

	UnAppliedTransaction map[common.Hash]*types.TransactionMetaData

	ApplyHandlers *SystemContract

	//just for demo test
	TestAccounts map[string]*UserAccount
}

//NewController create
func NewController(cfg CtlConfig) *Controller {
	return &Controller{
		Config:               cfg,
		UnAppliedTransaction: make(map[common.Hash]*types.TransactionMetaData, 0),
		Genesis:              types.NewGenesisState(),
		Pending:              NewPendingSate(),
		ApplyHandlers:        NewSystemContract(),
		TestAccounts:         make(map[string]*UserAccount),
	}
}

//Close release resource
func (cp *Controller) Close() {
	cp.AbortBlock()
}

//StartUp start up
func (cp *Controller) StartUp() error {
	var err error
	//new fork db
	cp.ForkDB, err = NewForkDB(cp, cp.Config.BlockLogDir)
	if err != nil {
		log.Printf("chain_plugin init new fork db err={%v}", err)
		return err
	}

	//new chain base
	cp.DB, err = chainbase.NewDataBase(cp.Config.ChainDir)
	if err != nil {
		log.Printf("chain_plugin init new chain db err={%v}", err)
		return err
	}

	//new block log
	cp.BlockLog, err = NewBlog(cp.Config.BlockLogDir, cp.Config.Cache, cp.Config.Handles)
	if err != nil {
		log.Printf("chain_plugin new block log err={%v}", err)
		return err
	}
	//init pending
	cp.Pending = NewPendingSate()

	//add index
	cp.addindex()
	cp.Head = cp.ForkDB.Head()

	//init db
	cp.initdb()

	//add ayatem contract
	cp.ApplyHandlers.Add("transfer", applyTransfer)

	return nil
}

//StartBlock method
func (cp *Controller) StartBlock(when *types.BlockTime, confirmcount uint16, status uint16) {
	if cp.Pending == nil {
		return
	}

	helper.CatchException(nil, func() {
		log.Print("Controller StartBlock panic")
	})

	// var bnum = int64(cp.Head.BlockNum)
	// revision := cp.DB.Revision()
	// if revision != bnum {
	// 	log.Printf("StartBlock db revision={%d} block num={%d}", revision, bnum)
	// 	// return
	// }

	// defer cp.Pending.Reset()

	cp.Pending.DBSession = cp.DB.StartUndoSession(true)
	cp.Pending.BlockStatus = status
	cp.Pending.PendingBlockState = types.NewBlockStateByTime(cp.Head.BlockHeaderState, when)
	cp.Pending.PendingBlockState.InChain = true

	cp.Pending.PendingBlockState.SetConfirm(confirmcount)

	wasPendingPromoted := cp.Pending.PendingBlockState.MaybePromotePending()
	raw, err := cp.DB.GetBaseValue(chainobject.GlobalPropertyType, 0)
	if err != nil {
		log.Printf("StartBlock db get global property object err={%v}", err)
		return
	}
	gpo := raw.(chainobject.GlobalPropertyObject)

	if gpo.ProposedScheduleBlockNum != 0 &&
		gpo.ProposedScheduleBlockNum <= cp.Pending.PendingBlockState.DposIrreversibleBlockNum &&
		len(cp.Pending.PendingBlockState.PendingSchedule.Producers) == 0 &&
		!wasPendingPromoted {
		cp.Pending.PendingBlockState.SetNewProducers(gpo.ProposedSchedule)
		gpo.ProposedScheduleBlockNum = 0
		gpo.ProposedSchedule.Clear()
		cp.DB.Modify(chainobject.GlobalPropertyType, gpo)
	}

	helper.CatchException(err, func() {
		log.Print("on block transaction failed, but shouldn't impact block generation ")
	})

	onTrxMeta := types.NewTrxMetaData(cp.getOnBlockTransaction())
	maxTime := types.MaxTime()
	cp.PushTransaction(onTrxMeta, maxTime, true, 0)

	cp.clearExpiredInputTransactions()
}

//SignBlock sign
func (cp *Controller) SignBlock(sig []byte, trust bool) {
	var pbs = cp.Pending.PendingBlockState
	pbs.Sign(sig, trust)
}

//ApplyBlock apply
func (cp *Controller) ApplyBlock(block *types.Block, status uint16) {
	//exception handler
	helper.CatchException(nil, func() {
		cp.AbortBlock()
		panic("ApplyBlock panic")
	})

	cp.StartBlock(block.TimeStamp, block.Confirmed, status)

	trace := &types.TransactionTrace{}

	for _, receipt := range block.Transactions {
		numReceipts := len(cp.Pending.PendingBlockState.Block.Transactions)

		if receipt.PackedTrx != nil {
			pt := receipt.PackedTrx
			mdtrx := types.NewMetaDataByPackedTrx(pt)

			trace = cp.PushTransaction(mdtrx, types.MaxTime(), false, receipt.CPUUsageUS)
		} else if len(receipt.TrxID) > 0 {
			trace = cp.PushScheduledTransaction(receipt.TrxID, types.MaxTime(), receipt.CPUUsageUS)
		}

		transactionFailed := trace != nil && trace.Except != nil
		canFailed := (receipt.Status == types.HardFail) && (len(receipt.TrxID) > 0)
		if transactionFailed && !canFailed {
			panic(trace.Except)
		}

		newsize := len(cp.Pending.PendingBlockState.Block.Transactions)

		if newsize == 0 || newsize != (numReceipts+1) {
			log.Printf("expected a receipt failed or not added. block={%v} receipt={%v}", block, receipt)
			return
		}

		receiptHead := cp.Pending.PendingBlockState.Block.Transactions[newsize-1].TransactionReceiptHeader
		if receipt.TransactionReceiptHeader != receiptHead {
			log.Printf("receipt does not match. block={%v} producer_receipt={%v} validator_receipt={%v}", block, receipt, receiptHead)
			return
		}
	}

	cp.FinalizeBlock()
	cp.SignBlock(block.Signature, false)

	cp.CommitBlock(false)

	return

}

//FinalizeBlock method
func (cp *Controller) FinalizeBlock() {
	if cp.Pending == nil {
		log.Print("it is not valid to finalize when there is no pending block")
		return
	}

	//execption handler
	helper.CatchException(nil, func() {
		panic("FinalizeBlock panic")
	})

	//start handler

	//update resource limits

	cp.setActionMerkle()
	cp.setTrxMerkle()

	p := cp.Pending.PendingBlockState
	p.ID = p.Header.Hash()
	p.Block.ID = p.ID

	cp.createBlockSummary(p.ID)

}

//PopBlock pop block
func (cp *Controller) PopBlock() {
	prev := cp.ForkDB.GetBlock(cp.Head.Header.Previous)
	if prev == nil {
		log.Print("attempt to pop beyond last irreversible block")
		return
	}

	for _, v := range cp.Head.Trxs {
		cp.UnAppliedTransaction[v.SignedID] = v
	}

	cp.Head = prev
	cp.DB.Undo()
}

//PushBlock method
func (cp *Controller) PushBlock(b *types.Block, status uint16) {
	if cp.Pending == nil {
		log.Print("it is not valid to push a block when there is a pending block")
		return
	}
	if status == types.Incomplete {
		return
	}
	trust := !cp.Config.ForceAllChecks && (status == types.Irreversible || status == types.Validate)
	state := cp.ForkDB.AddBlock(b, trust)
	if state != nil {
		cp.AcceptBlockHeaderChan <- state
		cp.maybeSwitchForks(status)
	}
}

//CommitBlock method
func (cp *Controller) CommitBlock(addtoforkdb bool) {
	helper.CatchException(nil, func() {
		cp.AbortBlock()
		panic("CommitBlock panic")
	})

	if addtoforkdb {
		cp.Pending.PendingBlockState.Validated = true

		log.Print("\n**************\n")
		log.Printf("commit block : %v   %v  %v", cp.Pending.PendingBlockState.Block.Producer, cp.Pending.PendingBlockState.Block.BlockNum, cp.Pending.PendingBlockState.Trxs)
		log.Print("\n**************\n")

		newbs := cp.ForkDB.AddState(cp.Pending.PendingBlockState)
		cp.AcceptBlockHeaderChan <- cp.Pending.PendingBlockState
		cp.Head = cp.ForkDB.Head()

		if newbs != cp.Head {
			log.Print("committed block did not become the new head in fork database")
			return
		}
	}

	cp.AcceptBlockChan <- cp.Pending.PendingBlockState

	cp.Pending.Push()
	cp.Pending.Reset()
	// log.Printf("commitblock: %v", time.Now().UnixNano()/int64(time.Millisecond))
}

//AbortBlock abort this block
func (cp *Controller) AbortBlock() {
	if cp.Pending != nil && cp.Pending.PendingBlockState != nil {
		for _, v := range cp.Pending.PendingBlockState.Trxs {
			cp.UnAppliedTransaction[v.SignedID] = v
		}

		cp.Pending.Reset()
	}
}

//AcceptBlock accept a block
func (cp *Controller) AcceptBlock(block *types.Block) {
	cp.SyncBlockChan <- block
}

//AcceptTransaction accept a packedTransaction from p2p
func (cp *Controller) AcceptTransaction(trx *types.PackedTransaction, next func(inerr error, trace *types.TransactionTrace)) {
	cp.AsyncPackedTrxChan <- &types.AsyncTrx{Pack: trx, Callback: next}
}

// PushTransaction This is the entry point for new transactions to the block state. It will check authorization and
/*  determine whether to execute it now or to delay it. Lastly it inserts a transaction receipt into
 *  the pending block.
 */
func (cp *Controller) PushTransaction(trx *types.TransactionMetaData, deadline time.Time, implicit bool, cputime uint32) *types.TransactionTrace {
	var result *types.TransactionTrace

	var err error
	helper.CatchException(err, func() {
		result.Except = err.(error)
	})

	trxContext := NewTransactionContext(cp, cp.DB, &trx.Trx, trx.ID)
	trxContext.DeadLine = deadline
	trxContext.BilledCPUTimeUS = int64(cputime)
	result = trxContext.TrxTrace

	if implicit {
		trxContext.InitForImplicitTrx(cp.pendingBlockTime())
		trxContext.CanSubjectivelyFail = false
	} else {
		trxContext.InitForInputTrx(cp.DB, cp.pendingBlockTime())
	}

	// if trxContext.CanSubjectivelyFail && cp.Pending.BlockStatus == types.Incomplete {
	// 	cp.checkActorList(trxContext.BillToAccounts)
	// }

	trxContext.Delay = trx.Trx.DelaySec

	//check authorization

	trxContext.Exec()
	trxContext.Finalize()

	if !implicit {
		var status uint8
		if trxContext.Delay == 0 {
			status = types.Executed
		} else {
			status = types.Delayed
		}

		result.Receipt = cp.PushReceipt(trx.PackedTrx, status, 0, 0).TransactionReceiptHeader
		cp.Pending.PendingBlockState.Trxs = append(cp.Pending.PendingBlockState.Trxs, trx)
		log.Printf("pushtrx Trxs: %v\n", cp.Pending.PendingBlockState.Trxs)
	} else {
		var r types.TransactionReceiptHeader
		r.Status = types.Executed
		r.CPUUsageUS = uint32(trxContext.BilledCPUTimeUS)
		r.NetUsageWords = uint(result.NetUsage) / 8
		result.Receipt = r
	}

	cp.Pending.Actions = append(cp.Pending.Actions, trxContext.Executed...)

	if !trx.Accepted {
		cp.AcceptTrxChan <- trx
		trx.Accepted = true
	}

	cp.AppliedTransactionChan <- result

	trxContext.Squash()

	if !implicit {
		delete(cp.UnAppliedTransaction, trx.SignedID)
	}

	return result
}

//PushScheduledTransaction handle some defered trx
func (cp *Controller) PushScheduledTransaction(trxid common.Hash, deadline time.Time, cputime uint32) *types.TransactionTrace {
	var result types.TransactionTrace

	return &result
}

//PushReceipt push receipt
func (cp *Controller) PushReceipt(trx interface{}, status uint8, cpuusage uint32, newusage uint) *types.TransactionReceipt {
	networds := newusage / 8

	var tr *types.TransactionReceipt

	switch trx.(type) {
	case types.PackedTransaction:
		pt := trx.(types.PackedTransaction)
		tr = types.NewTrxReceiptPacked(&pt)
	case common.Hash:
		id := trx.(common.Hash)
		tr = types.NewTrxReceiptID(id)
	default:
		return nil
	}

	tr.CPUUsageUS = cpuusage
	tr.NetUsageWords = networds
	tr.Status = status

	cp.Pending.PendingBlockState.Block.Transactions = append(cp.Pending.PendingBlockState.Block.Transactions, tr)
	return tr
}

//PushConfirmation push confirmation
func (cp *Controller) PushConfirmation(c *types.HeaderConfirmation) {
	if cp.Pending == nil {
		log.Print("it is not valid to push a confirmation when there is a pending block")
		return
	}

	cp.ForkDB.AddConfirmation(c)
	cp.AcceptedConfirmationChan <- c

	cp.maybeSwitchForks(types.Complete)
}

//HeadBlockNum get head block num
func (cp *Controller) HeadBlockNum() uint32 {
	return cp.Head.BlockNum
}

//HeadBlockID get head block id
func (cp *Controller) HeadBlockID() common.Hash {
	return cp.Head.ID
}

//HeadBlockTime return head time
func (cp *Controller) HeadBlockTime() time.Time {
	return time.Unix(cp.Head.Header.TimeStamp.Time.Int64(), 0)
}

//HeadBlockState return block state
func (cp *Controller) HeadBlockState() *types.BlockState {
	return cp.Head
}

//PendingBlockState get pending block state
func (cp *Controller) PendingBlockState() *types.BlockState {
	if cp.Pending != nil {
		return cp.Pending.PendingBlockState
	}

	return nil
}

//PendingBlockTime get pending block time
func (cp *Controller) PendingBlockTime() time.Time {
	return time.Unix(cp.Pending.PendingBlockState.Header.TimeStamp.Time.Int64(), 0)
}

//LastIrreversibleBlockNum return LastIrreversibleBlockNum
func (cp *Controller) LastIrreversibleBlockNum() uint32 {
	if cp.Head.BftIrreversibleBlockNum > cp.Head.DposIrreversibleBlockNum {
		return cp.Head.BftIrreversibleBlockNum
	}
	return cp.Head.DposIrreversibleBlockNum
}

//LastIrreversibleBlockID retur nLastIrreversibleBlockID
func (cp *Controller) LastIrreversibleBlockID() common.Hash {
	return cp.Head.ID
}

//GetUnappliedTransaction get unapplied trx
func (cp *Controller) GetUnappliedTransaction() []*types.TransactionMetaData {
	var result []*types.TransactionMetaData
	for _, v := range cp.UnAppliedTransaction {
		result = append(result, v)
	}

	return result
}

//DropUnappliedTransaction delete expired trx
func (cp *Controller) DropUnappliedTransaction(trx *types.TransactionMetaData) {
	delete(cp.UnAppliedTransaction, trx.SignedID)
}

func (cp *Controller) maybeSwitchForks(status uint16) {

	//start
	newHead := cp.ForkDB.Head()
	if newHead.Header.Previous == cp.Head.ID {
		//exception handle
		helper.CatchException(nil, func() {
			cp.ForkDB.SetValidity(newHead, false)
			panic("maybeSwitchForks on new head")
		})

		cp.ApplyBlock(newHead.Block, status)
		cp.ForkDB.MarkInChain(newHead, true)
		cp.ForkDB.SetValidity(newHead, true)
		cp.Head = newHead
	} else if newHead.ID != cp.Head.ID {
		first, second := cp.ForkDB.FetchBranch(newHead.ID, cp.Head.ID)

		for _, v := range second {
			cp.ForkDB.MarkInChain(v, false)
			cp.PopBlock()
		}

		prv := second[len(second)-1]
		if cp.HeadBlockID() != prv.Header.Previous {
			log.Print("loss of sync between fork_db and chainbase during fork switch")
			return
		}

		var recordApplied []*types.BlockState
		for i := len(first); i > 0; i-- {
			bs := first[i]

			//exception handle
			helper.CatchException(nil, func() {
				cp.ForkDB.SetValidity(bs, false)

				// pop all blocks from the bad fork
				// ritr base is a forward itr to the last block successfully applied
				for _, d := range recordApplied {
					cp.ForkDB.MarkInChain(d, false)
					cp.PopBlock()
				}

				prv := second[len(second)-1]
				if cp.HeadBlockID() != prv.Header.Previous {
					log.Print("loss of sync between fork_db and chainbase during fork switch")
					return
				}

				//// re-apply good blocks
				for j := len(second); j > 0; j-- {
					reb := second[j]

					cp.ApplyBlock(reb.Block, types.Validate)
					cp.Head = reb
					cp.ForkDB.MarkInChain(reb, true)
				}

				panic("maybeSwitchForks on fork")
			})

			var st uint16
			if bs.Validated {
				st = types.Validate
			} else {
				st = types.Complete
			}

			cp.ApplyBlock(bs.Block, st)
			cp.Head = bs
			cp.ForkDB.MarkInChain(bs, true)
			bs.Validated = true

			recordApplied = append(recordApplied, bs)

		}
		log.Printf("successfully switched fork to new head={%v}", newHead.ID)
	}

}

//OnIrreversible method
func (cp *Controller) OnIrreversible(s *types.BlockState) {
	if cp.BlockLog.Head() == nil {
		cp.BlockLog.ReadHead()
	}

	logHead := cp.BlockLog.Head()
	blockNum := logHead.BlockNum

	cp.IrreversibleBlockChan <- s
	cp.DB.Commit(int64(blockNum))

	if s.BlockNum <= blockNum {
		return
	}

	if (s.BlockNum-1) == blockNum || s.Block.Previous == logHead.ID {
		log.Printf("unlinkable block or irreversible doesn't link to block log head on blocknum={%v}", s.BlockNum)
		return
	}

	cp.BlockLog.Write(s.Block)
}

//SetProposedProducers set proposed producers
func (cp *Controller) SetProposedProducers(producers []chainobject.ProducerKey) int64 {
	return 0
}

//GetGlobalProperties get global property object
func (cp *Controller) GetGlobalProperties() chainobject.GlobalPropertyObject {
	raw, err := cp.DB.Get(chainobject.GlobalPropertyType, 0)
	if err != nil {
		log.Printf("GetGlobalProperties db get global property object err={%v}", err)
		return chainobject.GlobalPropertyObject{}
	}
	gpo := raw.(chainobject.GlobalPropertyObject)
	return gpo
}

//FetchBlockByID get blockby id
func (cp *Controller) FetchBlockByID(id common.Hash) *types.Block {
	state := cp.ForkDB.GetBlock(id)
	if state != nil {
		return state.Block
	}

	return nil
}

//FetchBlockByNum get block by block num
func (cp *Controller) FetchBlockByNum(num uint32) *types.Block {
	state := cp.ForkDB.GetBlockInChain(num)
	if state != nil {
		return state.Block
	}

	return cp.BlockLog.ReadByBlockNum(num)
}

//GetAccount get AccountObject
func (cp *Controller) GetAccount(name string) *chainobject.AccountObject {
	return nil
}

func (cp *Controller) getOnBlockTransaction() *types.SignedTransaction {
	var onBlockAct types.Action
	onBlockAct.Account = types.SystemAccountName
	onBlockAct.ActionName = "onblock"
	onBlockAct.Data, _ = rlp.EncodeToBytes(cp.Head.Header)

	var trx types.SignedTransaction
	trx.Actions = append(trx.Actions, onBlockAct)
	trx.SetReferenceBlock(cp.Head.Header)

	//add 1 second to the time
	trx.Expiration = cp.Pending.PendingBlockState.Header.TimeStamp.AddPlusTime(1)

	return &trx
}

func (cp *Controller) pendingBlockTime() time.Time {
	if cp.Pending != nil && cp.Pending.PendingBlockState != nil {
		return time.Unix(cp.Pending.PendingBlockState.Header.TimeStamp.Time.Int64(), 0)
	}

	return time.Now()
}

func (cp *Controller) checkActorList(actors map[string]struct{}) {

}

func (cp *Controller) addindex() {
	//add index for chain db
	cp.DB.AddIndex(chainobject.AccountType)
	cp.DB.AddIndex(chainobject.GlobalPropertyType)
	cp.DB.AddIndex(chainobject.TranscationType)
}

func (cp *Controller) initdb() {
	if cp.Head == nil {
		//init schedule
		var prodkey chainobject.ProducerKey
		prodkey.ProducerName = "datx" //system account
		prodkey.SigningKey = helper.RLPHash(cp.Genesis.InitKey)
		initschedule := chainobject.NewInitSchedule(prodkey)

		//set fork db head to the genesis state
		geneheader := types.MakeBlockHerderState()
		geneheader.ActiveSchedule = initschedule
		geneheader.PendingSchedule = initschedule
		geneheader.BlockNum = 0
		geneheader.Header.TimeStamp = types.NewBlockTime(cp.Genesis.InitTimeStamp)
		geneheader.ID = geneheader.Header.Hash()

		cp.Head = types.NewBlockStateByHeader(geneheader)
		cp.Head.Block = types.NewBlock(geneheader.Header)

		cp.ForkDB.AddState(cp.Head)
		var revision = int64(cp.Head.GetNum())
		cp.DB.SetRevision(revision)

		//create system acount,init global property object and dynamic global property object
		gpo := chainobject.NewGlobalPropertyObject()
		cp.DB.Create(chainobject.GlobalPropertyType, gpo)

		//push block
		bloghead := cp.BlockLog.ReadHead()
		if bloghead != nil && bloghead.BlockNum > 1 {
			next := cp.BlockLog.ReadByBlockNum(cp.Head.BlockNum + 1)
			for next != nil {
				cp.PushBlock(next, types.Irreversible)
				next = cp.BlockLog.ReadByBlockNum(cp.Head.BlockNum + 1)
			}
		} else {
			cp.BlockLog.ResetGenesis(cp.Head.Block)
		}
	}

	//undo pending changes when db revision is greater than head block number
	var undonum = int64(cp.Head.BlockNum)
	for cp.DB.Revision() > undonum {
		cp.DB.Undo()
	}
}

func (cp *Controller) clearExpiredInputTransactions() {

}

func (cp *Controller) setActionMerkle() {

}

func (cp *Controller) setTrxMerkle() {

}

func (cp *Controller) createBlockSummary(id common.Hash) {

}

//IsKnownUnexpiredTransaction check the trx is knoen
func (cp *Controller) IsKnownUnexpiredTransaction(id common.Hash) bool {
	return false
}
