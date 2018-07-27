package producerplugin

import (
	"datx_chain/chainlib/controller"
	"testing"
)

func TestProduceBlock(t *testing.T) {
	cfg := controller.CtlConfig{"block_log", "chain_log", false, false, 0, 0, "binaryen"}
	pp := NewProducerPlugin()
	pp.chain = controller.NewController(cfg)
	pp.Init()
	pp.Open()
	pp.produceBlock()
}
