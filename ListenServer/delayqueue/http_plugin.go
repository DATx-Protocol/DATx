package delayqueue

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//HTTPHandler func
type HTTPHandler func(w http.ResponseWriter, r *http.Request)

//HTTPPlugin struct
type HTTPPlugin struct {
	Host      string
	Port      string
	NewRouter *mux.Router
	Server    *http.Server
}

//NewHTTPPlugin new
func NewHTTPPlugin() *HTTPPlugin {
	p := &HTTPPlugin{
		NewRouter: mux.NewRouter(),
	}
	return p
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
	p.AddHandler("/redis_request", RedisRequest, "GET", "POST")
}
