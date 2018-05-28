package http_plugin

import (
	"datx_chain/utils"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	yaml "gopkg.in/yaml.v2"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request)

type Http_Config struct {
	//host address
	Host string `yaml:host`

	//listen port
	Port string `yaml:port`
}

type Http_plugin struct {
	//http server address
	Host string

	//http server listen port
	Port string

	//gorilla http router
	NewRouter *mux.Router

	//http server instance
	Server *http.Server
}

func NewHttpPlugin() *Http_plugin {
	p := &Http_plugin{
		NewRouter: mux.NewRouter(),
	}
	return p
}

func (p *Http_plugin) Init() {
	log.Println("Http server initialize")

	err, data := utils.GetFileHelper("http_config.yaml")
	if err != nil {
		log.Printf("err %s", err)
		return
	}

	var conf Http_Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		log.Printf("chain_plugin init unmarshal config  error={%v}", err)
		return
	}

	log.Printf("Http server init config=%v", conf)

	p.Host = conf.Host
	p.Port = conf.Port

	srv := &http.Server{
		Handler:      p.NewRouter,
		Addr:         p.Host + ":" + p.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	p.Server = srv
}

func (p *Http_plugin) InitWithEndpoint(host string, port string) {
	log.Println("Http server initialize")

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

func (p *Http_plugin) Open() {
	log.Println("Http server start")

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := p.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

}

func (p *Http_plugin) Close() {
	p.Server.Close()

	log.Println("Http server closed")
}

func (p *Http_plugin) AddHandler(url string, handler HttpHandler, methods ...string) {
	if len(url) == 0 || handler == nil {
		return
	}

	p.NewRouter.HandleFunc(url, handler).Methods(methods...)
}
