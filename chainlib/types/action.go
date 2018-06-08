package types

type Action struct {
	//account name
	Account string

	//action name
	ActionName string

	//Authorization

	//action dat
	Data []byte
}

//Action constructor
func NewAction(account, name string, data []byte) *Action {
	return &Action{
		Account:    account,
		ActionName: name,
		Data:       data,
	}
}
