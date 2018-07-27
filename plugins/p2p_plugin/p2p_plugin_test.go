package p2p_plugin

import (
	"datx_chain/plugins/p2p_plugin/p2p"
	"datx_chain/utils/rlp"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/robfig/cron"
)

func initConfigFile(configFileName string) error {
	content :=
		`key : 112350326689685529176416839667165445850930628506430118257577941014994870680344
maxpeers : 256
versions : [ 1, 2] 
network : 12
listen : 0.0.0.0:4345
name : Test`
	// if _, err := os.Stat("config"); os.IsNotExist(err) {
	// 	os.Mkdir("config", os.ModePerm)
	// }

	data := []byte(content)
	return ioutil.WriteFile(configFileName, data, 0644)
}

func TestP2p_Plugin_Open(t *testing.T) {
	versions := [...]uint{1, 2, 3}
	p2pConfig := P2pPluginConfig{
		PrivateKey:      "87754059b19801941602f125ee1a747825e99568b9832ab2b3a1f47cdd6da6a8",
		MaxPeers:        25,
		ProtocolVersion: versions[0:],
		NetworkId:       24,
		ListenAddr:      "0.0.0.0:2342",
		Name:            "test",
		BootstrapNodes:  []string{"enode://0x49bf39a33057c4f09952245f9ae8db947c32a44d01bc2ec556ef2417c418a3748146bf545c45220ca9320a05fe039b4c9c8a8231151ef2da5e09a6dfd1199312@10.3.58.6:30303?discport=30301"},
	}
	p := NewP2pPlugin()
	if err := p.InitWithConfig(&p2pConfig); err != nil {
		t.Errorf("init with config file error:%v", err)
	}

	if err := p.Open(); err != nil {
		t.Errorf("p2p_plugin start error:%v", err)
	}

}

func TestP2p_Plugin_Connect(t *testing.T) {
	versions := [...]uint{1, 2, 3, 63}
	p2pConfig := P2pPluginConfig{
		PrivateKey:      "87754059b19801941602f125ee1a747825e99568b9832ab2b3a1f47cdd6da6a7",
		MaxPeers:        256,
		ProtocolVersion: versions[0:],
		NetworkId:       12,
		ListenAddr:      "0.0.0.0:34235",
		Name:            "Test",
	}
	p := NewP2pPlugin()
	if err := p.InitWithConfig(&p2pConfig); err != nil {
		t.Errorf("init with config file error:%v", err)
	}

	if err := p.Open(); err != nil {
		t.Errorf("p2p_plugin start error:%v", err)
	}

	nodeAddr := "enode://10683b69c543586411f08da56bfa7fff17eefe9fb31c3bc3e853421dd19c230c2641f980b60bbb28042bf87eea6ccfebf29e860bcbb36eaf4350bfeb7144551d@172.31.1.91:34235"
	if err := p.Connect(nodeAddr); err != nil {
		t.Errorf("connect error:%v", err)
	}
	for _, peer := range p.pm.peers.peers {
		p2p.Send(peer.rw, MsgTypeGoAwayMsg, GoAwayMsg{Reason: ErrOtherFatal})
	}
}

func TestRlpEncode(t *testing.T) {
	HandShake := HandShakeMsg{
		NetworkVersion:           1,
		NetworkId:                12,
		LastIrreversibleBlockNum: 8,
		PeerID:     "asr134fwe1454yy6",
		Generation: 1,
		TimeStamp:  time.Now(),
	}

	size, r, _ := rlp.EncodeToReader(HandShake)
	Msg := p2p.Msg{Code: 0, Size: uint32(size), Payload: r}

	HSMsg := HandShakeMsg{}
	Msg.Decode(&HSMsg)

}

func TestScheduleTask(t *testing.T) {
	i := 0
	c := cron.New()
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})
	c.Start()
	select {}
}
