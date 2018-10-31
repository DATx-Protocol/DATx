package server

import (
	"datx/lsd/chainlib"
	"datx/lsd/delayqueue"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

//ExpirationTable
type ExpirationTable struct {
	Rows []struct {
		ID        int    `json:"id"`
		Trxid     string `json:"trxid"`
		Producer  string `json:"producer"`
		Timestamp int    `json:"timestamp"`
		Category  string `json:"category"`
	} `json:"rows"`
	More bool `json:"more"`
}

type ExtractArgs struct {
	Trxid    string `json:"trxid"`
	Producer string `json:"producer"`
	Category string `json:"category"`
}

type ExpiredArgs struct {
	Trxid     string `json:"trxid"`
	Category  string `json:"category"`
	Handler   string `json:"handler"`
	Timestamp string `json:"time"`
}

//Extract
type Extract struct {
	pos int32

	offset int32

	taskClose chan bool

	producerName string
}

func NewExtract(name string) *Extract {
	return &Extract{
		pos:          0,
		offset:       10,
		taskClose:    make(chan bool),
		producerName: name,
	}
}

//Startup
func (ext *Extract) Startup() {
	go ext.taskLoop()
}

//Close
func (ext *Extract) Close() {
	ext.taskClose <- true
}

func (ext *Extract) getExpiredTrxs() ([]chainlib.Transaction, error) {
	//get expired trx from extract smart contract table
	raw, err := GetOuterTrxTable("datxos.extra", "datxos.extra", "expiration")

	var temp ExpirationTable
	err = json.Unmarshal(raw, &temp)
	if err != nil {
		return nil, err
	}

	var result []chainlib.Transaction
	for _, v := range temp.Rows {
		var item chainlib.Transaction
		item.TransactionID = v.Trxid
		item.Category = v.Category

		result = append(result, item)
	}

	return result, nil
}

func (ext *Extract) getAllTransactions() ([]chainlib.Transaction, error) {
	var result []chainlib.Transaction

	//dbtc
	btcTrxs, err := GetExtractActions("datxos.dbtc", ext.pos, ext.offset)
	if err != nil {
		return nil, err
	}
	result = append(result, btcTrxs...)

	//deth
	ethTrxs, err := GetExtractActions("datxos.deth", ext.pos, ext.offset)
	if err != nil {
		return nil, err
	}
	result = append(result, ethTrxs...)

	//deos
	eosTrxs, err := GetExtractActions("datxos.deos", ext.pos, ext.offset)
	if err != nil {
		return nil, err
	}
	result = append(result, eosTrxs...)
	ext.pos = ext.pos + ext.offset

	expire, err := ext.getExpiredTrxs()
	if err == nil {
		result = append(result, expire...)
	}

	return result, nil
}

//PushTrxToQueue push unirreversible trx to queue
func (ext *Extract) pushTrxToQueue(trx chainlib.Transaction) {
	var job delayqueue.Job
	job.Topic = trx.Category
	job.Id = trx.Category + "_" + trx.TransactionID
	job.Delay = time.Now().Unix() + 1
	job.TTR = 60

	bytes, err := json.Marshal(trx)
	if err != nil {
		log.Printf("trx marshal failed:%v %v\n", trx, err)
		return
	}
	job.Body = string(bytes)

	if err = delayqueue.Push(job); err != nil {
		log.Printf("PushTrxToQueue Push queue failed:%v   %v\n", trx, err)
		return
	}
}

//PushExtractAction push trx to contract that check the trx is sended already
func (ext *Extract) pushExtractAction(trx chainlib.Transaction) error {
	//
	args := ExtractArgs{
		Trxid:    trx.TransactionID,
		Producer: ext.producerName,
		Category: trx.Category,
	}

	//remove
	jobid := trx.Category + "_" + trx.TransactionID
	delayqueue.Remove(jobid)

	//push action
	bytes, err := json.Marshal(args)
	if err != nil {
		log.Printf("PushExtractAction marshal failed:%v %v\n", trx, err)
		return err
	}
	_, err = chainlib.ClPushAction("datxos.extra", "recordtrx", string(bytes), ext.producerName)
	if err != nil {
		log.Printf("Extract push recordtrx failed:%v %v\n", trx, err)
		return err
	}

	var trxID string
	trxID, err = func(t chainlib.Transaction) (string, error) {
		switch t.Category {
		case "DBTC":
			return BTCMultiSig(trx)
		case "DETH":
			return ETHMultiSig(trx)
		case "DEOS":
			return EOSMultiSig(trx)
		default:
			return "", fmt.Errorf("PushExtractAction category %v not defined", trx.Category)
		}
	}(trx)

	//multisig failed
	if err != nil {
		log.Printf("MultiSig failed trxID:%v %v\n", trxID, err)
		// return err
	}

	log.Printf("Extract finished: %v\n", trx)

	return nil
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

				irreversible := CheckIrreversible(trx)
				if irreversible {
					go ext.pushExtractAction(trx)
				}
			}
		}
	}

}

//ExecExtract
func (ext *Extract) ExecExtract() {
	trxlist, err := ext.getAllTransactions()
	if err != nil || len(trxlist) == 0 {
		return
	}

	for _, v := range trxlist {
		if v.IsIrrevisible {
			go ext.pushExtractAction(v)
		} else {
			ext.pushTrxToQueue(v)
		}
	}
}
