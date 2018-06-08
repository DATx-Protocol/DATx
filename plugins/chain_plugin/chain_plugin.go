package chain_plugin

import (
	"log"
	"sync"

	"datx_chain/chainlib/types"
	"datx_chain/utils/helper"
	"datx_chain/utils/message"

	yaml "gopkg.in/yaml.v2"
)

type Chain_Config struct {
	//block log dir
	Block_Log_Dir string `yaml:"block_log_dir"`

	// //genesis time stamp
	// Genesis_Time int64 `yaml:"genesis_time"`

	// //genesis filr path
	// Genesis_File string `yaml:"genesis_file"`

	//
	Read_Only bool `yaml:"read_only"`

	//db handles of open file capacity
	Handles int `yaml:"handles"`

	//db can cache block capacity
	Cache int `yaml:"cache"`

	//vm type
	VM_Type string `yaml:"vm_type"`
}

type chain_plugin struct {
	//default config
	Config Chain_Config

	//fork db
	Fork_DB *ForkDB

	//chain id
	Chain_ID string

	//accept block chan
	Block_Chan chan *message.Msg

	//accept transcation chan
	Trx_Chan chan *message.Msg
}

func (self *chain_plugin) Init() {
	//catch exception, do nothing if catch exception
	helper.CatchException(func() {
		return
	})

	//unmarshal yaml file
	err, data := helper.GetFileHelper("chain_config.yaml")
	if err != nil {
		log.Printf("chain_plugin init error={%v}", err)
		return
	}

	var config Chain_Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Printf("chain_plugin init unmarshal config  error={%v}", err)
		return
	}

	log.Printf("chain_plugin init config={%v}", config)
	self.Config = config

	//new fork db
	self.Fork_DB, err = NewForkDB(self.Config.Block_Log_Dir, self.Config.Cache, self.Config.Handles)
	if err != nil {
		log.Printf("chain_plugin init new fork db err={%v}", err)
	}
}

func (self *chain_plugin) Open() {

}

func (self *chain_plugin) Close() {

}

func (self *chain_plugin) AcceptBlock(block *types.Block) {
	msg := message.NewMsg(message.BlockMsg, block)

	msg.Send(self.Block_Chan)
}

func (self *chain_plugin) AccpetTranscation(packed_trx *types.Transcation) {
	msg := message.NewMsg(message.TrxMsg, packed_trx)

	msg.Send(self.Trx_Chan)
}

//singleton pattern of thread safety
var oneinstance *chain_plugin
var once sync.Once

func GetInstance() *chain_plugin {
	//exec only once
	once.Do(func() {
		oneinstance = &chain_plugin{}
	})

	return oneinstance
}
