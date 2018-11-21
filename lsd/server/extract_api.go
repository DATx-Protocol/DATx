package server

import (
	"datx/lsd/chainlib"
	"datx/lsd/delayqueue"
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
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
	btcpos int32

	ethpos int32

	eospos int32

	offset int32

	taskClose chan bool

	producerName string

	lastIrreversibleBlockNum int64
}

func NewExtract(name string) *Extract {
	return &Extract{
		btcpos:                   0,
		ethpos:                   0,
		eospos:                   0,
		offset:                   9,
		taskClose:                make(chan bool),
		producerName:             name,
		lastIrreversibleBlockNum: 0,
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
		item.IsIrrevisible = true

		result = append(result, item)
	}

	return result, nil
}

//GetExtractActions get transaction by escrow account(*dbtc,deth,deos) when extracting
func (ext *Extract) getExtractActions(addr string, pos, offset int32) ([]chainlib.Transaction, int32, error) {
	cmdStr := fmt.Sprintf("cldatx get actions %s -j %d %d", addr, pos, offset)
	res, err := chainlib.ExecShell(cmdStr)
	if err != nil {
		log.Printf("\nGetAccountActions get actions return: %s\n", err)
		return nil, 0, err
	}

	var resp ExtractActions
	if err := json.Unmarshal([]byte(res), &resp); err != nil {
		return nil, 0, err
	}

	ext.lastIrreversibleBlockNum = int64(resp.LastIrreversibleBlock)
	result := make([]chainlib.Transaction, 0)
	for _, v := range resp.Actions {
		if v.ActionTrace.Act.Name != "extract" {
			continue
		}

		if v.ActionTrace.Act.Data.To != addr {
			continue
		}

		var temp chainlib.Transaction
		temp.TransactionID = v.ActionTrace.TrxID

		temp.BlockNum = int64(v.BlockNum)
		temp.From = v.ActionTrace.Act.Data.From
		temp.To = v.ActionTrace.Act.Data.To
		amountpos := strings.Index(v.ActionTrace.Act.Data.Quantity, " ")
		amountstr := v.ActionTrace.Act.Data.Quantity[:amountpos]
		temp.Category = v.ActionTrace.Act.Data.Quantity[amountpos+1:]
		temp.Amount, err = strconv.ParseFloat(amountstr, 64)
		temp.Memo = v.ActionTrace.Act.Data.Memo
		if err != nil {
			log.Printf("Extract parse amount err:%v\n", err)
			return nil, 0, err
		}

		temp.Time, err = time.Parse("2006-01-02T15:04:05", v.BlockTime)
		if err != nil {
			log.Printf("Extract parse time err:%v\n", err)
			return nil, 0, err
		}
		temp.IsIrrevisible = false
		if ext.lastIrreversibleBlockNum >= temp.BlockNum {
			temp.IsIrrevisible = true
		}

		result = append(result, temp)
	}

	return result, int32(len(resp.Actions)), nil
}

func (ext *Extract) getAllTransactions() ([]chainlib.Transaction, error) {
	var result []chainlib.Transaction

	//dbtc
	btcTrxs, btcplus, err := ext.getExtractActions("datxos.dbtc", ext.btcpos, ext.offset)
	if err != nil {
		return nil, err
	}
	result = append(result, btcTrxs...)
	ext.btcpos = ext.btcpos + btcplus

	//deth
	ethTrxs, ethplus, err := ext.getExtractActions("datxos.deth", ext.ethpos, ext.offset)
	if err != nil {
		return nil, err
	}
	result = append(result, ethTrxs...)
	ext.ethpos = ext.ethpos + ethplus

	//deos
	eosTrxs, eosplus, err := ext.getExtractActions("datxos.deos", ext.eospos, ext.offset)
	if err != nil {
		return nil, err
	}
	result = append(result, eosTrxs...)
	ext.eospos = ext.eospos + eosplus

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
func (ext *Extract) pushExtractAction(trx chainlib.Transaction) (err error) {

	defer func() {
		if errs := recover(); errs != nil {
			log.Printf("[Extract] pushExtractAction panic,error: %v\n", errs)
			debug.PrintStack()
			err = fmt.Errorf("[Extract] pushExtractAction panic, %v", errs)
		}
	}()

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
		log.Printf("MultiSig failed: %v %v\n", trxID, err)
		return err
	}

	log.Printf("[Extract] multisig success: %v\n", trx)

	return nil
}

func (ext *Extract) taskLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
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

				if ext.lastIrreversibleBlockNum >= trx.BlockNum {
					ext.pushExtractAction(trx)
				} else {
					log.Printf("[Extract] trx timeout and not irreversible: %v < %v %v", ext.lastIrreversibleBlockNum, trx.BlockNum, time.Now())
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
			ext.pushExtractAction(v)
		} else {
			ext.pushTrxToQueue(v)
		}
	}
}
