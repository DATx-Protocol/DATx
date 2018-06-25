package application

import (
	"datx_chain/plugins/chain_plugin"
	"datx_chain/plugins/http_plugin"
	"errors"
	"fmt"
	"sync"
)

type Application struct {
	//the version of application
	version uint64

	//the pairs of name/plugin
	plugins map[string]Pluginer
}

func (self *Application) SetVersion(v uint64) {
	self.version = v
}

func (self *Application) AddPlugin(name string) error {
	if _, ok := self.plugins[name]; ok {
		return errors.New(fmt.Sprintf("the plugin name={%s} already added.", name))
	}

	var plugin Pluginer
	switch name {
	case "chain_plugin":
		plugin = chain_plugin.NewChainPlugin()
	case "http_plugin":
		plugin = http_plugin.NewHttpPlugin()
	default:
		return errors.New(fmt.Sprintf("the plugin instance name={%s} has not defineted.", name))
	}

	self.plugins[name] = plugin

	return nil
}

func (self *Application) Find(name string) (Pluginer, error) {
	if value, ok := self.plugins[name]; ok {
		return value, nil
	}

	return nil, errors.New(fmt.Sprintf("the plugin name={%s} not found.", name))
}

//singleton pattern of thread safety
var oneinstance *Application
var once sync.Once

func App() *Application {
	//exec only once
	once.Do(func() {
		oneinstance = &Application{}
	})

	return oneinstance
}
