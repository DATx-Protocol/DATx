package server

import (
	"datx/ListenServer/chainlib"
	"datx/ListenServer/delayqueue"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type ChainServer struct {
	trxPool    chan chainlib.Transaction
	hashPool   chan chainlib.Transaction
	clLimit    chan bool // 限制push transfer并发数
	taskClose  chan bool
	queryClose chan bool
	maxnum     int32
	browser    map[string]chainlib.CommonFunc

	extract *Extract
}

func NewChainServer(maxnum int32) *ChainServer {
	return &ChainServer{
		trxPool:    make(chan chainlib.Transaction, maxnum),
		hashPool:   make(chan chainlib.Transaction, maxnum),
		clLimit:    make(chan bool, 100),
		taskClose:  make(chan bool),
		queryClose: make(chan bool),
		browser:    make(map[string]chainlib.CommonFunc),
		extract:    NewExtract(),
	}
}

func (tick *ChainServer) Start() {
	//start delay queue
	delayqueue.InitQueue()

	go tick.queryLoop()
	go tick.taskLoop() //
	go tick.execLoop()
	go tick.datxLoop()

	tick.extract.Startup()
}

func (tick *ChainServer) Close() {
	tick.extract.Close()

	tick.taskClose <- true
	tick.queryClose <- true
	close(tick.trxPool)
	close(tick.hashPool)

	for _, v := range tick.browser {
		v.Close()
	}
}

func (tick *ChainServer) AddBrowser(cate string, browser chainlib.CommonFunc) {
	tick.browser[cate] = browser
}

func (tick *ChainServer) AddHash(trx chainlib.Transaction) {
	tick.hashPool <- trx
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

				fmt.Printf("delay finished: %v\n", trx.TransactionID)
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

		fmt.Printf("delete job id=%v %v\n", jobid, time.Now().Unix())
		delayqueue.Remove(jobid)
	}
}

func (tick *ChainServer) datxLoop() {
	for trx := range tick.hashPool {
		tick.clLimit <- true
		go func(trx chainlib.Transaction) {
			err := chainlib.ClWaitIrreversible(trx.BlockNum)
			if err != nil {
				log.Println("wait irreversible failed\t", err)
				<-tick.clLimit
				return
			}
			var trans chainlib.TransferInfo
			trans.Hash = trx.TransactionID
			trans.From = "datxio.d" + strings.ToLower(trx.Category)
			trans.To = trx.To
			trans.Quantity = strconv.FormatFloat(trx.Amount, 'f', 4, 64) + " D" + trx.Category
			trans.Memo = trx.Memo

			// 需要钱包密码
			_, err = chainlib.ClWalletUnlock("PW5JHPpaGrS7bKhmQJ5Rb7rNSXhp3S3sXN2fGWaqQNzQufQaWrkUJ")
			// 需要合约权限
			_, err = chainlib.ClPushTransfer("datxio.charg", "transtoken", trans)
			if err != nil {
				log.Println("push transfer failed\t", err)
				<-tick.clLimit
				return
			}
			<-tick.clLimit
		}(trx)
	}
}

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
