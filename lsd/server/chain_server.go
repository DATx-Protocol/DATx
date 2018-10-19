package server

import (
	"datx/lsd/chainlib"
	"datx/lsd/delayqueue"
	"encoding/json"
	"fmt"
	"time"
)

type ChainServer struct {
	trxPool    chan chainlib.Transaction
	taskClose  chan bool
	queryClose chan bool
	maxnum     int32
	browser    map[string]chainlib.CommonFunc

	extract *Extract
}

//NewChainServer new chain server
func NewChainServer(maxnum int32) *ChainServer {
	return &ChainServer{
		trxPool:    make(chan chainlib.Transaction, maxnum),
		taskClose:  make(chan bool),
		queryClose: make(chan bool),
		browser:    make(map[string]chainlib.CommonFunc),
		extract:    NewExtract(chainlib.GetCfgProducerName()),
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
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-tick.queryClose:
			return
		case <-ticker.C:
			{
				isCurProducer := IsCurrentProducer()
				if !isCurProducer {
					break
				}

				//charge
				for _, v := range tick.browser {
					v.Tick()
				}

				//extract
				tick.extract.ExecExtract()
			}
		}
	}

}

func (tick *ChainServer) taskLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	topics := []string{"BTC", "ETH", "EOS"}

	for {
		select {
		case <-tick.taskClose:
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

				fmt.Printf("ChainServer delay finished: %v   %v\n", trx.TransactionID, time.Now().UnixNano())
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
			fmt.Printf("Trx type is not support:%v\n", trx.Category)
			continue
		}

		isSuccess := browser.ReTry(trx)
		if !isSuccess {
			continue
		}

		//exec success,delete it
		jobid := trx.Category + "_" + trx.TransactionID

		fmt.Printf("ChainServer delete job id=%v %v\n", jobid, time.Now().Unix())
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
		fmt.Printf("trx marshal failed:%v %v\n", trx, err)
		return
	}
	job.Body = string(bytes)

	if err = delayqueue.Push(job); err != nil {
		fmt.Printf("Push queue failed.%v\n", err)
		return
	}
}
