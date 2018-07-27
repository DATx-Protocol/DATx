package chainplugin

import (
	"datx_chain/chainlib/application"
	"datx_chain/chainlib/controller"
	"datx_chain/utils/common"
	"datx_chain/utils/helper"
	"fmt"
	"log"

	yaml "gopkg.in/yaml.v2"
)

//ChainPlugin struct
type ChainPlugin struct {
	//default config
	Config controller.CtlConfig

	chain *controller.Controller

	//chain id
	chainID common.Hash
}

//NewChainPlugin new
func NewChainPlugin() *ChainPlugin {
	return &ChainPlugin{
		chainID: helper.RLPHash("chainplugin"),
	}
}

//Init method
func (cp *ChainPlugin) Init() error {
	//catch exception, do nothing if catch exception

	//unmarshal yaml file
	err, data := helper.GetFileHelper("chain_config.yaml", application.App().GetConfigFolder())
	if err != nil {
		log.Printf("chain_plugin init chain config error={%v}", err)
		return err
	}

	var config controller.CtlConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Printf("chain_plugin init chain config unmarshal config  error={%v}", err)
		return err
	}

	log.Printf("chain_plugin init config={%v}", config)
	cp.Config = config

	cp.chain = controller.NewController(cp.Config)
	return nil
}

//Open method
func (cp *ChainPlugin) Open() (err error) {
	defer func() {
		if nerr := recover(); nerr != nil {
			str := fmt.Sprintf("chain Plugin Open panic={%v}", nerr)
			err = nerr.(error)
			panic(str)
		}
	}()

	if err := cp.chain.StartUp(); err != nil {
		log.Printf("chain Plugin Open err={%v}", err)
		return err
	}

	cp.chainID = cp.chain.Genesis.ComputeChainID()
	log.Printf("BlockChain started;head block num is {%d}, genesis timestamp is {%v}", cp.chain.HeadBlockNum(), cp.chain.Genesis.InitTimeStamp)
	return nil
}

//Close method
func (cp *ChainPlugin) Close() {

}

//Chain get controller
func (cp *ChainPlugin) Chain() *controller.Controller {
	return cp.chain
}

//GetChainID get chain id
func (cp *ChainPlugin) GetChainID() common.Hash {
	return cp.chainID
}

//GetChainConfig return controller config
func (cp *ChainPlugin) GetChainConfig() controller.CtlConfig {
	return cp.chain.Config
}
