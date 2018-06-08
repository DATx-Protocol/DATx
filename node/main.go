package main

import (
	"datx_chain/plugins/chain_plugin"
	"datx_chain/plugins/http_plugin"
	"log"
)

func main() {
	// nodekey, _ := crypto.GenerateKey()
	// srv := p2p.Server{
	// 	Config: p2p.Config{
	// 		MaxPeers:   10,
	// 		PrivateKey: nodekey,
	// 		Name:       "datx",
	// 		ListenAddr: ":30300",
	// 		Protocols:  []p2p.Protocol{MyProtocol()},
	// 	},
	// }

	// if err := srv.Start(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	//http server
	httpServer := http_plugin.NewHttpPlugin()
	httpServer.Init()
	httpServer.Open()
	defer httpServer.Close()

	//test chain_plugin
	sd := chain_plugin.GetInstance()
	sd.Init()
	log.Printf("chain %v", sd.Config)

	select {}
}
