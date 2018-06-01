package p2p

import (
	"datx_chain/plugins/p2p_plugin/p2p/nat"
	"datx_chain/utils/common"
	"datx_chain/utils/crypto"
	"log"
)

type P2p_Plugin struct {

	//http server instance
	p2pServer *Server
}

func NewP2pPlugin() *P2p_Plugin {
	p := &P2p_Plugin{}
	return p
}

func (p *P2p_Plugin) Init() error {
	log.Println("P2p plugin initialize")

	nodekey, _ := crypto.GenerateKey()
	p.p2pServer = &Server{
		Config: Config{
			PrivateKey: nodekey,
			MaxPeers:   10,
			Name:       common.MakeName("wnode", "6.0"),
			// Protocols:      shh.Protocols(),
			ListenAddr: ":30300",
			NAT:        nat.Any(),
			// BootstrapNodes: peers,
			// StaticNodes:    peers,
			// TrustedNodes:   peers,
		},
	}

	// err, data := utils.GetFileHelper("p2p_config.yaml")
	// if err != nil {
	// 	log.Printf("err %s", err)
	// 	return
	// }

	// var conf Http_Config
	// if err := yaml.Unmarshal(data, &conf); err != nil {
	// 	log.Printf("chain_plugin init unmarshal config  error={%v}", err)
	// 	return
	// }

	// log.Printf("P2p plugin init config=%v", conf)

	// p.Host = conf.Host
	// p.Port = conf.Port
	return nil
}

func (p *P2p_Plugin) Open() error {
	log.Println("P2p plugin start")
	return p.p2pServer.Start()

	// Run our server in a goroutine so that it doesn't block.
	// go func() {
	// 	if err := p.Server.ListenAndServe(); err != nil {
	// 		log.Println(err)
	// 	}
	// }()

}

func (p *P2p_Plugin) Close() {
	log.Println("P2p plugin closed")
	p.p2pServer.Stop()
}
