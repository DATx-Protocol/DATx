package p2p_plugin

import (
	"datx_chain/plugins/p2p_plugin"
	"testing"
)

func TestP2p_Plugin_Init(t *testing.T) {
	tests := []struct {
		name    string
		p       *p2p_plugin.P2p_Plugin
		wantErr bool
	}{
		{
			name:    "123",
			p:       p2p_plugin.NewP2pPlugin(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Init(); (err != nil) != tt.wantErr {
				t.Errorf("P2p_Plugin.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
