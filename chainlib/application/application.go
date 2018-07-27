package application

import (
	"errors"
	"fmt"
	"sync"
)

//Application manage the plugins
type Application struct {
	//the version of application
	version uint64

	configFolder string

	//the pairs of name/plugin
	plugins map[string]Pluginer

	sequence []string
}

//NewApp new
func NewApp() *Application {
	return &Application{
		plugins:  make(map[string]Pluginer, 1),
		sequence: make([]string, 0),
	}
}

//SetVersion set version
func (app *Application) SetVersion(v uint64) {
	app.version = v
}

func (app *Application) SetConfigFolder(folder string) {
	app.configFolder = folder
}

func (app *Application) GetConfigFolder() string {
	return app.configFolder
}

//AddPlugin add plugin to app
func (app *Application) AddPlugin(name string, plugin Pluginer) error {
	if _, ok := app.plugins[name]; ok {
		str := fmt.Sprintf("the plugin name={%s} already added.", name)
		return errors.New(str)
	}

	app.sequence = append(app.sequence, name)
	app.plugins[name] = plugin

	return nil
}

//Start init and open all plugins
func (app *Application) Start() error {
	for i := 0; i < len(app.sequence); i++ {
		name := app.sequence[i]

		if v, ok := app.plugins[name]; ok {
			if err := v.Init(); err != nil {
				return err
			}
			if err := v.Open(); err != nil {
				return err
			}
		}
	}

	return nil
}

//Close all plugins
func (app *Application) Close() {
	for _, v := range app.plugins {
		v.Close()
	}
}

//Find get pluginer by name
func (app *Application) Find(name string) (Pluginer, error) {
	if value, ok := app.plugins[name]; ok {
		return value, nil
	}

	str := fmt.Sprintf("the plugin name={%s} not found.", name)
	return nil, errors.New(str)
}

//singleton pattern of thread safety
var oneinstance *Application
var once sync.Once

//App the manager of plugins
func App() *Application {
	//exec only once
	once.Do(func() {
		oneinstance = NewApp()
	})

	return oneinstance
}
