package httpplugin

import (
	"datx_chain/chainlib/application"
	"datx_chain/utils/helper"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	yaml "gopkg.in/yaml.v2"
)

//HTTPHandler func
type HTTPHandler func(w http.ResponseWriter, r *http.Request)

//HTTPConfig config
type HTTPConfig struct {
	//host address
	Host string `yaml:"host"`

	//listen port
	Port string `yaml:"port"`
}

//HTTPPlugin struct
type HTTPPlugin struct {
	//http server address
	Host string

	//http server listen port
	Port string

	//gorilla http router
	NewRouter *mux.Router

	//http server instance
	Server *http.Server
}

//NewHTTPPlugin new
func NewHTTPPlugin() *HTTPPlugin {
	p := &HTTPPlugin{
		NewRouter: mux.NewRouter(),
	}
	return p
}

//Init init
func (p *HTTPPlugin) Init() error {
	log.Println("HTTP server initialize")

	err, data := helper.GetFileHelper("http_config.yaml", application.App().GetConfigFolder())
	if err != nil {
		log.Printf("err %s", err)
		return err
	}

	var conf HTTPConfig
	if err := yaml.Unmarshal(data, &conf); err != nil {
		log.Printf("chain_plugin init unmarshal config  error={%v}", err)
		return err
	}

	log.Printf("HTTP server init config=%v", conf)

	p.Host = conf.Host
	p.Port = conf.Port

	srv := &http.Server{
		Handler:      p.NewRouter,
		Addr:         p.Host + ":" + p.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	p.Server = srv
	return nil
}

//InitWithEndpoint init with host and port
func (p *HTTPPlugin) InitWithEndpoint(host string, port string) {
	log.Println("HTTP server initialize")

	p.Host = host
	p.Port = port

	srv := &http.Server{
		Handler:      p.NewRouter,
		Addr:         p.Host + ":" + p.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	p.Server = srv
}

//Open open plugin
func (p *HTTPPlugin) Open() (err error) {
	log.Println("HTTP server start")

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if nerr := p.Server.ListenAndServe(); nerr != nil {
			log.Println(nerr)
			err = nerr
		}
	}()
	p.RegisterHandler()
	return nil
}

//Close close plugin
func (p *HTTPPlugin) Close() {
	p.Server.Close()

	log.Println("HTTP server closed")
}

//AddHandler add handler
func (p *HTTPPlugin) AddHandler(url string, handler HTTPHandler, methods ...string) {
	if len(url) == 0 || handler == nil {
		return
	}

	p.NewRouter.HandleFunc(url, handler).Methods(methods...)
}

//RegisterHandler register all handler
func (p *HTTPPlugin) RegisterHandler() {
	p.AddHandler("/transfer", TransferHandler, "GET", "POST")
	p.AddHandler("/transaction_list", GetTransactionListHandle, "GET", "POST")
	p.AddHandler("/transaction_query_hash", GetTransactionByHashHandle, "GET", "POST")
	p.AddHandler("/blocks_list", GetBlockListHandle, "GET", "POST")
	p.AddHandler("/general_info", GetGeneralInfo, "GET", "POST")
}
