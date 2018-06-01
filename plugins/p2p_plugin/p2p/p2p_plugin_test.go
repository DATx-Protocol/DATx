package p2p

import (
	"datx_chain/plugins/p2p_plugin/p2p"
	"testing"
)

func TestP2p_Plugin_Open(t *testing.T) {
	p2pPlugin := p2p_plugin.NewP2pPlugin()
	if err := p2pPlugin.Init(); err != nil {
		t.Errorf("p2p plugin init error:%v", err)
	}
	if err := p2pPlugin.Open(); err != nil {
		t.Errorf("p2p start init error:%v", err)
	}
	p2pPlugin.Close()
	t.Log("p2p unit test")
}
