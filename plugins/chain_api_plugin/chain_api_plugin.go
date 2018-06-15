package chain_api_plugin

import (
	"datx_chain/plugins/chain_plugin"
	"datx_chain/plugins/http_plugin"
	"log"
	"net/http"
	"reflect"
	"sync"
)

//test chain api
type Chain_API struct {
	//chain controller db,temp init to int ,it will change next time.
	Fork_DB *ForkDB
}

func (self *Chain_API) get_chain_plugin_type(instance interface{}) interface{} {
	log.Println("get Chain_plugin_Type.")
	chain_plugin_type := reflect.TypeOf(instance)
	log.Println(chain_plugin_type)
	return chain_plugin_type
}

var chain_plugin_instance = chain_plugin.GetInstance()
var once sync.Once

func Init() {
	once.Do(func() {
		chain_plugin_instance.Init()
		log.Println(" Success,Chain plugin has been inited.")
	})
}
func Startup() {
	log.Println(chain_plugin_instance)
	chain_plugin_instance.Open()
	log.Println(" Success,Chain plugin has been opened.")
}
func Shutdown() {
	chain_plugin_instance.Close()
	log.Println(" Success,Chain plugin has been shutdown.")
}
func AddToHttpHandler(host string, port string) {
	tempServer := http_plugin.NewHttpPlugin()
	tempServer.InitWithEndpoint(host, port)
	tempServer.Open()
	var TestHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("block info is as follow：\n"))
	}
	tempServer.AddHandler("/get_block_info", TestHandler, "get")
}
