package main

import (
	"datx/ListenServer/gatway"
	"datx/ListenServer/http"
	"datx/ListenServer/server"
	"log"
)

func main() {

	//start delay queue
	chainserver := server.NewChainServer(10)
	chainserver.Start()
	defer chainserver.Close()

	eth := gatway.NewETHBrowser("0x3d74f927f4a1c9d5c66acc597cb269cd31b69a89 ", chainserver)
	chainserver.AddBrowser("ETH", eth)

	btc := gatway.NewBTCBrowser("1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F", chainserver)
	chainserver.AddBrowser("BTC", btc)

	eos := gatway.NewEOSBrowser("eostea111111", chainserver)
	chainserver.AddBrowser("EOS", eos)

	server.GetOuterTrxTable("user", "user", "games")

	//自测通过
	httpServer := http.NewHTTPServer()
	httpServer.InitWithEndpoint("localhost", "8880")
	if err := httpServer.Open(); err != nil {
		log.Printf("httpaccessory open err={%v}", err)
		return
	}
	defer httpServer.Close()

	select {}
}
