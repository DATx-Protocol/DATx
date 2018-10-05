package server

import (
	"datx/ListenServer/chainlib"
	"datx/ListenServer/delayqueue"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Extract struct {
	Pos int32

	Offset int32

	LatestBlockNum int64

	taskClose chan bool
}

func NewExtract() *Extract {
	return &Extract{
		Pos:            0,
		Offset:         10,
		LatestBlockNum: 0,
		taskClose:      make(chan bool),
	}
}

func (ext *Extract) Startup() {
	ext.getPosAndOffset()

	go ext.taskLoop()
}

func (ext *Extract) Close() {
	ext.taskClose <- true
}

func (ext *Extract) getPosAndOffset() error {
	_, err := GetOuterTrxTable("datxio", "datxio", "status")
	if err != nil {
		return err
	}

	return nil
}

func (ext *Extract) GetAllTransactions() ([]chainlib.Transaction, error) {
	err := ext.getPosAndOffset()
	if err != nil {
		return nil, fmt.Errorf("GetAllTransactions failed: %v\n", err)
	}

	var result []chainlib.Transaction

	//dbtc
	btcTrxs, err := GetExtractActions("datxio.dbtc", ext.Pos, ext.Offset)
	if err != nil {
		return nil, err
	}
	result = append(result, btcTrxs...)

	//deth
	ethTrxs, err := GetExtractActions("datxio.deth", ext.Pos, ext.Offset)
	if err != nil {
		return nil, err
	}
	result = append(result, ethTrxs...)

	//deos
	eosTrxs, err := GetExtractActions("datxio.deos", ext.Pos, ext.Offset)
	if err != nil {
		return nil, err
	}
	result = append(result, eosTrxs...)

	ext.Pos = ext.Pos + ext.Offset
	return result, nil
}

//PushTrxToQueue push unirreversible trx to queue
func (ext *Extract) PushTrxToQueue(trx chainlib.Transaction) {
	var job delayqueue.Job
	job.Topic = trx.Category
	job.Id = trx.Category + "_" + trx.TransactionID
	job.Delay = time.Now().Unix() + 1
	job.TTR = 60

	bytes, err := json.Marshal(trx)
	if err != nil {
		fmt.Printf("trx marshal failed:%v %v\n", trx, err)
		return
	}
	job.Body = string(bytes)

	if err = delayqueue.Push(job); err != nil {
		fmt.Printf("PushTrxToQueue Push queue failed:%v   %v\n", trx, err)
		return
	}
}

//PushExtractAction push trx to contract that check the trx is sended already
func (ext *Extract) PushExtractAction(trx chainlib.Transaction) error {
	//
	bytes, err := json.Marshal(trx)
	if err != nil {
		log.Printf("PushExtractAction marshal failed:%v %v\n", trx, err)
		return err
	}
	_, err = chainlib.ClPushAction("contract_account", "check", string(bytes))
	if err != nil {
		log.Printf("Extract push check failed:%v %v\n", trx, err)
		return err
	}

	var trxID string
	trxID, err = func(t chainlib.Transaction) (string, error) {
		switch t.Category {
		case "BTC":
			return BTCMultiSig(trx)
		case "ETH":
			return ETHMultiSig(trx)
		case "EOS":
			return EOSMultiSig(trx)
		default:
			return "", fmt.Errorf("PushExtractAction category %v not defined\n", trx.Category)
		}
	}(trx)

	//multisig failed,rollback
	if err != nil {
		log.Printf("MultiSig failed trxID:%v %v\n", trxID, err)
		return err
	}

	//extract success
	log.Printf("MultiSig success trxID:%v\n", trxID)

	jobid := trx.Category + "_" + trx.TransactionID

	log.Printf("PushExtractAction success delete job id=%v %v\n", jobid, time.Now().Unix())
	delayqueue.Remove(jobid)

	return nil
}

func (ext *Extract) ExecExtract() {
	trxlist, err := ext.GetAllTransactions()
	if err != nil {
		return
	}

	for _, v := range trxlist {
		if v.IsIrrevisible {
			go ext.PushExtractAction(v)
		} else {
			ext.PushTrxToQueue(v)
		}
	}
}

func (ext *Extract) taskLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	topics := []string{"DBTC", "DETH", "DEOS"}

	for {
		select {
		case <-ext.taskClose:
			return
		case <-ticker.C:
			{
				task, err := delayqueue.Pop(topics)
				if err != nil || task == nil {
					break
				}

				var trx chainlib.Transaction
				if err = json.Unmarshal([]byte(task.Body), &trx); err != nil {
					break
				}

				fmt.Printf("Extract delay coming: %v\n", trx)
				irreversible := CheckIrreversible(trx)
				if irreversible {
					go ext.PushExtractAction(trx)
				} else {
					ext.PushTrxToQueue(trx)
				}
				fmt.Printf("Extract finished: %v\n", trx)
			}
		}
	}

}
