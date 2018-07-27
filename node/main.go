package main

import (
	"datx_chain/chainlib/application"
	"datx_chain/plugins/chain_plugin"
	"datx_chain/plugins/http_plugin"
	"datx_chain/plugins/p2p_plugin"
	"datx_chain/plugins/producer_plugin"
	"datx_chain/utils/helper"
	"flag"
	"log"
)

func main() {
	configFolder := "config"
	flag.StringVar(&configFolder, "c", "config", "set configuration `file`")
	flag.Parse()
	//set app version
	application.App().SetVersion(1)
	application.App().SetConfigFolder(configFolder)

	helper.CatchException(nil, func() {
		log.Print("main panic")
	})

	//chain plugin
	chainplugin := chainplugin.NewChainPlugin()
	application.App().AddPlugin("chain", chainplugin)

	//http server
	httpServer := httpplugin.NewHTTPPlugin()
	application.App().AddPlugin("http", httpServer)

	//p2p plugin
	p2p := p2p_plugin.NewP2pPlugin()
	application.App().AddPlugin("p2p", p2p)
	log.Printf("config folder:%v,%v", configFolder, application.App().GetConfigFolder())

	//producer plugin
	producer := producerplugin.NewProducerPlugin()
	application.App().AddPlugin("producer", producer)

	//start all plugins
	if err := application.App().Start(); err != nil {
		log.Printf("Start err={%v}", err)
		return
	}
	nodeInfo := p2p.Status()
	log.Printf("NodeInfo:%v", *nodeInfo)

	defer application.App().Close()

	select {}
}
