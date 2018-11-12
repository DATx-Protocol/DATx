package server

import (
	"datx/lsd/chainlib"
	"datx/lsd/common"
	"datx/lsd/delayqueue"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type ChainServer struct {
	trxPool    chan chainlib.Transaction
	taskClose  chan bool
	queryClose chan bool
	maxnum     int32
	browser    map[string]chainlib.CommonFunc
	count      int32

	extract *Extract
}

//NewChainServer new chain server
func NewChainServer(maxnum int32) *ChainServer {
	return &ChainServer{
		trxPool:    make(chan chainlib.Transaction, maxnum),
		taskClose:  make(chan bool),
		queryClose: make(chan bool),
		browser:    make(map[string]chainlib.CommonFunc),
		count:      1,
		extract:    NewExtract(common.GetCfgProducerName()),
	}
}

//Start method
func (tick *ChainServer) Start() {
	go tick.queryLoop()
	go tick.taskLoop() //
	go tick.execLoop()

	tick.extract.Startup()
}

//Close method
func (tick *ChainServer) Close() {
	tick.extract.Close()

	tick.taskClose <- true
	tick.queryClose <- true
	close(tick.trxPool)

	for _, v := range tick.browser {
		v.Close()
	}
}

//AddBrowser add (btc,eth,eos) browser
func (tick *ChainServer) AddBrowser(cate string, browser chainlib.CommonFunc) {
	tick.browser[cate] = browser
}

func (tick *ChainServer) queryLoop() {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-tick.queryClose:
			return
		case <-ticker.C:
			{
				isCurProducer := IsCurrentProducer()
				if !isCurProducer {
					log.Print("\nChainServer the producer-name of config.ini is not current producer.\n\n ")
					break
				}

				//charge
				for _, v := range tick.browser {
					go v.Tick()
				}

				//get expire trx from charge expiration table
				tick.pushChargeExpiredTrxs()

				//update contract table for expireation trx
				tick.count++
				if tick.count >= 15 {
					tick.count = 1
					tick.updateexpiretable()
				}

				//extract
				tick.extract.ExecExtract()
			}
		}
	}

}

func (tick *ChainServer) taskLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	topics := []string{"BTC", "ETH", "EOS"}

	for {
		select {
		case <-tick.taskClose:
			return
		case <-ticker.C:
			{
				isCurProducer := IsCurrentProducer()
				if !isCurProducer {
					break
				}

				task, err := delayqueue.Pop(topics)
				if err != nil || task == nil {
					break
				}

				var trx chainlib.Transaction
				if err = json.Unmarshal([]byte(task.Body), &trx); err != nil {
					break
				}

				// log.Printf("ChainServer delay finished: %v   %v\n", trx.TransactionID, time.Now().UnixNano())
				tick.trxPool <- trx
			}
		}
	}

}

func (tick *ChainServer) execLoop() {
	for trx := range tick.trxPool {
		//exec method
		browser, ok := tick.browser[trx.Category]
		if !ok {
			log.Printf("Trx type is not support:%v\n", trx.Category)
			continue
		}

		isSuccess := browser.ReTry(trx)
		if !isSuccess {
			continue
		}

		//exec success,delete it
		jobid := trx.Category + "_" + trx.TransactionID

		// log.Printf("ChainServer delete job id=%v %v\n", jobid, time.Now().Unix())
		delayqueue.Remove(jobid)
	}
}

//AddTask method
func (tick *ChainServer) AddTask(trx chainlib.Transaction, delay int64) {
	var job delayqueue.Job
	job.Topic = trx.Category
	job.Id = trx.Category + "_" + trx.TransactionID
	job.Delay = time.Now().Unix() + delay
	job.TTR = 60

	bytes, err := json.Marshal(trx)
	if err != nil {
		log.Printf("trx marshal failed:%v %v\n", trx, err)
		return
	}
	job.Body = string(bytes)

	if err = delayqueue.Push(job); err != nil {
		log.Printf("Push queue failed.%v\n", err)
		return
	}
}
func (tick *ChainServer) pushChargeExpiredTrxs() error {
	//get expired trx from extract smart contract table
	log.Print("[ChainServer] Get expiration transaction and push action again.\n")
	extraStr := fmt.Sprint("cldatx get table datxos.charg datxos.charg expiration")
	raw, err := chainlib.ExecShell(extraStr)
	if err != nil {
		log.Printf("[ChainServer] Get expiration transaction and push action failed: %v\n", err)
		return err
	}

	var temp ChargeExpirationTable
	err = json.Unmarshal([]byte(raw), &temp)
	if err != nil {
		log.Printf("[ChainServer] Get expiration transaction and unmarsh failed: %v\n", err)
		return err
	}

	for _, v := range temp.Rows {
		if v.Count >= 3 {
			continue
		}
		var item chainlib.Transaction
		item.TransactionID = v.Trxid
		item.Category = v.Category
		item.From = v.From
		item.To = v.To
		item.Amount, _ = strconv.ParseFloat(v.Quantity, 64)
		item.BlockNum = int64(v.Blocknum)
		item.IsIrrevisible = true
		item.Memo = v.Memo

		chainlib.PushCharge(item)
	}

	return nil
}

func (tick *ChainServer) updateexpiretable() error {
	log.Print("[ChainServer] update expiration table.\n")
	extraStr := fmt.Sprintf("cldatx push action datxos.extra updateexpire '' -p %s -f", common.GetCfgProducerName())
	_, err := chainlib.ExecShell(extraStr)
	if err != nil {
		log.Printf("[ChainServer] update datxos.extra expiration table failed: %v\n", err)
		return err
	}

	rollbackStr := fmt.Sprintf("cldatx push action datxos.extra rollbacktrx '' -p %s -f", common.GetCfgProducerName())
	_, err = chainlib.ExecShell(rollbackStr)
	if err != nil {
		log.Printf("[ChainServer] update datxos.extra rollbacktrx failed: %v\n", err)
		return err
	}

	chargStr := fmt.Sprintf("cldatx push action datxos.charg updateexptrx '' -p %s -f", common.GetCfgProducerName())
	_, err = chainlib.ExecShell(chargStr)
	if err != nil {
		log.Printf("[ChainServer] update datxos.charg expiration table failed: %v\n", err)
		return err
	}

	return nil
}
