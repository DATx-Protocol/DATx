package p2p_plugin

import (
	"crypto/ecdsa"
	"datx_chain/plugins/p2p_plugin/p2p"
	"datx_chain/utils/crypto"
	"datx_chain/utils/db"
	"datx_chain/utils/helper"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

type P2pPluginConfig struct {
	PrivateKey      string
	MaxPeers        int
	ProtocolVersion []uint
	NetworkId       uint64
	ListenAddr      string
	Name            string
}

type P2p_Plugin struct {
	p2pServer  *p2p.Server
	peers      *peerSet
	pm         *ProtocolManager
	privateKey *ecdsa.PrivateKey
	quitSync   chan struct{}
	chainDb    *datxdb.LDBDatabase
}

func NewP2pPlugin() *P2p_Plugin {
	return &P2p_Plugin{
		peers:    newPeerSet(),
		quitSync: make(chan struct{}),
	}
}

func (p *P2p_Plugin) Init() error {

	err, data := helper.GetFileHelper("p2p_config.yaml")
	if err != nil {
		log.Printf("err %s", err)
		return nil
	}

	var p2pConfig P2pPluginConfig
	if err := yaml.Unmarshal(data, &p2pConfig); err != nil {
		log.Printf("p2p_plugin init unmarshal config  error={%v}", err)
		return err
	}

	log.Printf("p2p_plugin init config=%v", p2pConfig)
	pm, err := NewProtocolManager(p2pConfig.ProtocolVersion, p2pConfig.NetworkId, p.peers, p.chainDb, p.quitSync, new(sync.WaitGroup))
	if err != nil {
		log.Printf("p2p_plugin new protocol error:%v", err.Error())
		return err
	}

	// nodekey, _ := crypto.GenerateKey()
	// nodekey, _ := crypto.GenerateKey()//nodekey 要传入进来
	p.pm = pm
	privateKey, err := crypto.HexToECDSA(p2pConfig.PrivateKey)
	if err != nil {
		log.Printf("convert private key error,%v", err.Error())
		return err
	}
	p.p2pServer = &p2p.Server{
		Config: p2p.Config{
			MaxPeers:   p2pConfig.MaxPeers,
			PrivateKey: privateKey,
			Name:       p2pConfig.Name,
			ListenAddr: p2pConfig.ListenAddr,
			Protocols:  p.pm.SubProtocols,
		},
	}
	return nil
}

func (p *P2p_Plugin) Open() error {
	p.pm.Start(p.p2pServer.MaxPeers)
	if err := p.p2pServer.Start(); err != nil {
		log.Printf("p2p_plugin p2pserver start error,%v", err.Error())
		return err
	}
	return nil
}

func (p *P2p_Plugin) Close() {
	p.pm.Stop()
	p.p2pServer.Stop()
}
