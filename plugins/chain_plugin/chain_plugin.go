package chain_plugin

import (
	"log"
	"sync"

	"datx_chain/utils"
	"datx_chain/utils/message"

	"github.com/golang/leveldb"
	yaml "gopkg.in/yaml.v2"
)

type Chain_Config struct {
	//block log dir
	Block_Log_Dir string `yaml:block_log_dir`

	// //genesis time stamp
	// Genesis_Time int64 `yaml:genesis_time`

	// //genesis filr path
	// Genesis_File string `yaml:genesis_file`

	//
	Read_Only bool `yaml:read_only`

	//vm type
	VM_Type string `yaml:vm_type`
}

type chain_plugin struct {
	//default config
	Config Chain_Config

	//fork db
	Fork_DB *leveldb.DB

	//chain id
	Chain_ID string

	//accept block chan
	Block_Chan chan *message.Msg

	//accept transcation chan
	Trx_Chan chan *message.Msg
}

func (self *chain_plugin) Init() {
	//unmarshal yaml file
	err, data := utils.GetFileHelper("chain_config.yaml")
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
}

func (self *chain_plugin) Open() {

}

func (self *chain_plugin) Close() {

}

func (self *chain_plugin) AcceptBlock(block *utils.Block) {
	msg := message.NewMsg(message.BlockMsg, block)

	msg.Send(self.Block_Chan)
}

func (self *chain_plugin) AccpetTranscation(packed_trx *utils.Transcation) {
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
