package types

//PermissionLevel struct
type PermissionLevel struct {
	Actor      string
	Permission string
}

//Action struct
type Action struct {
	//account name
	Account string

	//action name
	ActionName string

	//Authorization
	Authorization []PermissionLevel

	//action dat
	Data []byte
}

//NewAction Action constructor
func NewAction(account, name string, data []byte) *Action {
	return &Action{
		Account:    account,
		ActionName: name,
		Data:       data,
	}
}
