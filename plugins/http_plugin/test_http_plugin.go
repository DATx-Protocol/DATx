package httpplugin

import (
	"log"
	"net/http"
)

//TestDefaultHTTPServer test
func TestDefaultHTTPServer() {
	testServer := NewHTTPPlugin()

	//start your http server based on http_plugin.yaml
	testServer.Init()
	testServer.Open()

	//test handler
	var TestHelloHandler = func(w http.ResponseWriter, r *http.Request) {
		log.Println("[TestHelloHandler] start")

		w.Write([]byte("hello, welcome to http_plugin!\n"))
	}

	//register your url and handler
	testServer.AddHandler("/get_info", TestHelloHandler, "get")
}

//TestHTTPServer test
func TestHTTPServer(host string, port string) {
	testServer := NewHTTPPlugin()

	//start your http server based on passed parameter
	testServer.InitWithEndpoint(host, port)
	testServer.Open()

	//test handler
	var TestHelloHandler = func(w http.ResponseWriter, r *http.Request) {
		log.Println("[TestHelloHandler] start")

		w.Write([]byte("hello, welcome to http_plugin!\n"))
	}

	//register your url and handler
	testServer.AddHandler("/get_info", TestHelloHandler, "get")
}
