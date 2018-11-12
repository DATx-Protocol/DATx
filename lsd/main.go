package main

import (
	"datx/lsd/common"
	"datx/lsd/delayqueue"
	"datx/lsd/gatway"
	"datx/lsd/http"
	"datx/lsd/server"
	"log"
	"runtime/debug"
)

func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("main panic,error: %v\n", err)
			debug.PrintStack()
		}
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	//start delay queue
	delayqueue.InitQueue()

	//ping redis server
	ping := delayqueue.PingRedisPool()
	if !ping {
		return
	}

	//start redis server
	delayqueue.StartRedis()

	//start chain server
	chainserver := server.NewChainServer(10)
	chainserver.Start()
	defer chainserver.Close()

	ethAccount := common.GetTrusteeAccount("eth-muladdress")
	log.Printf("main get eth trustee account: %s\n", ethAccount)
	eth := gatway.NewETHBrowser(ethAccount, chainserver)
	chainserver.AddBrowser("ETH", eth)

	btcAccount := common.GetTrusteeAccount("btc-muladdress")
	log.Printf("main get btc trustee account: %s\n", btcAccount)
	btc := gatway.NewBTCBrowser(btcAccount, chainserver)
	chainserver.AddBrowser("BTC", btc)

	eosAccount := common.GetTrusteeAccount("eos-mulAccount")
	log.Printf("main get eos trustee account: %s\n", eosAccount)
	eos := gatway.NewEOSNode("http://213.239.208.37:8888", eosAccount, chainserver)
	chainserver.AddBrowser("EOS", eos)

	//start http server
	httpServer := http.NewHTTPServer()
	httpServer.InitWithEndpoint("localhost", "8880")
	if err := httpServer.Open(); err != nil {
		log.Printf("httpaccessory open err={%v}", err)
		return
	}
	defer httpServer.Close()

	select {}
}
