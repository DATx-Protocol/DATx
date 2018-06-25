package application

type Pluginer interface {
	Init()
	Open()
	Close()
}
