package p2p_plugin

import (
	"crypto/ecdsa"
	"datx_chain/chainlib/application"
	"datx_chain/chainlib/controller"
	"datx_chain/chainlib/types"
	"datx_chain/plugins/p2p_plugin/p2p"
	"datx_chain/plugins/p2p_plugin/p2p/discover"
	"datx_chain/utils/crypto"
	"datx_chain/utils/db"
	"datx_chain/utils/helper"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

type P2pPluginConfig struct {
	PrivateKey      string   `yaml:"PrivateKey"`
	MaxPeers        int      `yaml:"MaxPeers"`
	ProtocolVersion []uint   `yaml:"ProtoclVersion"`
	NetworkId       uint64   `yaml:"NetworkId"`
	ListenAddr      string   `yaml:"listenAddr"`
	Name            string   `yaml:"name"`
	BootstrapNodes  []string `yaml:"bootstrap"`
}

type P2p_Plugin struct {
	p2pServer  *p2p.Server
	peers      *peerSet
	pm         *ProtocolManager
	privateKey *ecdsa.PrivateKey
	quitSync   chan struct{}
	chainDb    *datxdb.LDBDatabase
	p2pConfig  *P2pPluginConfig
}

func NewP2pPlugin() *P2p_Plugin {
	P2pPlugin_impl = &P2p_Plugin{
		peers:    newPeerSet(),
		quitSync: make(chan struct{}),
	}
	return P2pPlugin_impl
}

func (p *P2p_Plugin) InitWithConfig(p2pConfig *P2pPluginConfig) error {

	log.Printf("p2p_plugin init config=%v", p2pConfig)
	pm, err := NewProtocolManager(p2pConfig.ProtocolVersion, p2pConfig.NetworkId, p.peers, p.chainDb, p.quitSync, new(sync.WaitGroup))
	if err != nil {
		log.Printf("p2p_plugin new protocol error:%v", err.Error())
		return err
	}
	p.pm = pm
	privateKey, err := crypto.HexToECDSA(p2pConfig.PrivateKey)
	if err != nil {
		log.Printf("convert private key error,%v", err.Error())
		return err
	}

	nodes := make([]*discover.Node, len(p2pConfig.BootstrapNodes))
	nodes = nodes[0:0]

	for _, bootstrap := range p2pConfig.BootstrapNodes {
		node, err := discover.ParseNode(bootstrap)
		if err != nil {
			log.Printf("parse node error:%v", err)
		} else {
			nodes = append(nodes, node)
		}
	}
	log.Printf("bootstrap node len:%v", len(nodes))
	p.p2pServer = &p2p.Server{
		Config: p2p.Config{
			MaxPeers:       p2pConfig.MaxPeers,
			PrivateKey:     privateKey,
			Name:           p2pConfig.Name,
			ListenAddr:     p2pConfig.ListenAddr,
			Protocols:      p.pm.SubProtocols,
			BootstrapNodes: nodes,
		},
	}

	if len(nodes) > 0 {
		p.p2pServer.Config.BootstrapNodes = nodes
	}
	return nil
}

func (p *P2p_Plugin) Init() error {
	if p.p2pConfig == nil {
		configFileName := "p2p_config.yaml"
		err, data := helper.GetFileHelper(configFileName, application.App().GetConfigFolder())
		if err != nil {
			log.Printf("err %s", err)
			return err
		}

		var p2pConfig P2pPluginConfig
		if err := yaml.Unmarshal(data, &p2pConfig); err != nil {
			log.Printf("p2p_plugin init unmarshal config  error={%v}", err)
			return err
		}
		p.p2pConfig = &p2pConfig
		log.Printf("data:%v, p2pConfig:%v", string(data), p.p2pConfig)
	}

	return p.InitWithConfig(p.p2pConfig)
}

func (p *P2p_Plugin) Open() error {
	p.pm.Start(p.p2pServer.MaxPeers)
	p.listenAllChannels(p.pm.syncMaster.chain)
	log.Printf("p2p plugin open")
	if err := p.p2pServer.Start(); err != nil {
		log.Printf("p2p_plugin p2pserver start error,%v", err.Error())
		return err
	}
	return nil
}

func (p *P2p_Plugin) Close() {
	log.Printf("p2p plugin close")
	p.pm.Stop()
	p.p2pServer.Stop()
}

func (p *P2p_Plugin) Connect(rawNodeUrl string) error {
	node, err := discover.ParseNode(rawNodeUrl)
	if err != nil {
		log.Printf("parse node error:%v", err)
		return err
	}
	p.p2pServer.AddPeer(node)
	return nil
}

func (p *P2p_Plugin) Status() *p2p.NodeInfo {
	return p.p2pServer.NodeInfo()
}

func (p *P2p_Plugin) listenAllChannels(ct *controller.Controller) {
	if ct == nil {
		return
	}
	ct.AcceptBlockChan = make(chan *types.BlockState, 10)
	ct.AcceptBlockHeaderChan = make(chan *types.BlockState, 10)
	ct.IrreversibleBlockChan = make(chan *types.BlockState, 10)
	ct.AcceptTrxChan = make(chan *types.TransactionMetaData, 10)
	ct.AppliedTransactionChan = make(chan *types.TransactionTrace, 10)
	ct.AcceptedConfirmationChan = make(chan *types.HeaderConfirmation, 10)
	ct.TransactionAckChan = make(chan *types.TrxTrace, 10)
	go func() {
		for {
			select {
			case acceptBlockChan := <-ct.AcceptBlockChan:
				// log.Printf("accpertBlockChan is={%v}", acceptBlockChan)
				go P2pPlugin_impl.pm.dispatchMaster.bcastBlock(acceptBlockChan)
			case <-ct.AcceptBlockHeaderChan:
			case <-ct.IrreversibleBlockChan:
			case <-ct.AcceptTrxChan:
			case <-ct.AppliedTransactionChan:
			case <-ct.AcceptedConfirmationChan:
			case TranscationAckChan := <-ct.TransactionAckChan:
				go func() {
					if TranscationAckChan.Err != nil {
						P2pPlugin_impl.pm.dispatchMaster.rejectedTransaction(TranscationAckChan.Trx.ID())
					} else {
						P2pPlugin_impl.pm.dispatchMaster.bcastTransaction(TranscationAckChan.Trx)
					}
				}()
			}
		}
	}()

}
