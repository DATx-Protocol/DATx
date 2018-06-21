package p2p_plugin

import (
	"crypto/ecdsa"
	"encoding/binary"
	"math"
	"sync"
	"datx_chain/utils/common"
	"datx_chain/utils/log"
	"datx_chain/utils/rlp"
)

type P2pPluginConfig struct {
	PrivateKey ecdsa.PrivateKey,
	MaxPeers   uint32,
	ProtoclVersion []string,
	NetworkId  uint32,
	ListenAddr string
}


type P2p_Plugin struct {
	p2pServer *p2p.Server,
	peers * peerSet,
	pm *ProtocolManager，
	privateKey *ecdsa.PrivateKey,
	quitSync   chan struct{},
	chainDb db.Database
}


func NewP2pPlugin() P2p_Plugin {
	return &P2p_Plugin {
		peers: newPeerSet(),
		quitSync: make(chan struct)
	}
}

func (p2p *P2p_Plugin) Init() error {

	err, data := helper.GetFileHelper("p2p_config.yaml")
	if err != nil {
		log.Printf("err %s", err)
		return
	}

	var p2pConfig P2pPluginConfig
	if err := yaml.Unmarshal(data, &p2pConfig); err != nil {
		log.Printf("p2p_plugin init unmarshal config  error={%v}", err)
		return
	}

	log.Printf("p2p_plugin init config=%v", p2pConfig)
	pm, err := NewProtocolManager(p2pConfig.protoclVersion, p2pConfig.NetworkId, p2p.peers, p2p.chainDb,  p2p.quitSync, new(sync.WaitGroup))
	if err != nil {
		log.Printf("p2p_plugin new protocol error:%v", err.Error())
		return err
	}

	// nodekey, _ := crypto.GenerateKey()
	// nodekey, _ := crypto.GenerateKey()//nodekey 要传入进来
	privateKey := crypto.HexToECDSA(p2pConfig.privateKey)
	p2pServer = &p2p.Server{
		MaxPeers:   p2pConfig.maxPeers,
		PrivateKey: privateKey,
		Name:       p2pConfig.Name,
		ListenAddr: p2pConfig.listenAddr,
		Protocols:  p2p.pm.SubProtocols,
	}
	reurn  nil
}

func (p2p *P2p_Plugin) Open() error {
	p2p.pm.Start()
	p2p.p2pServer.start()
}

func (p2p * P2p_Plugin) Close() {
	p2p.pm.Stop()
	p2p.p2pServer.Stop()
}

