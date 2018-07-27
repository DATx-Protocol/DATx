package application

//Pluginer plugin interface
type Pluginer interface {
	Init() error
	Open() error
	Close()
}
