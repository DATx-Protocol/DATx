package controller

import (
	"encoding/json"
	"log"
	"strings"
)

//Transfer struct
type Transfer struct {
	From   string
	To     string
	Amount uint16
}

//SysHandle system contract handle function
type SysHandle func(a *ApplyContent)

//SystemContract manager
type SystemContract struct {
	Handler map[string]SysHandle
}

//NewSystemContract new
func NewSystemContract() *SystemContract {
	return &SystemContract{
		Handler: make(map[string]SysHandle, 0),
	}
}

//Add pairs of name/handler
func (sys *SystemContract) Add(name string, h SysHandle) {
	sys.Handler[name] = h
}

//Find get handler by name
func (sys *SystemContract) Find(name string) (SysHandle, bool) {
	var res SysHandle
	if v, ok := sys.Handler[name]; ok {
		return v, true
	}

	return res, false
}

var applyTransfer = func(a *ApplyContent) {
	act := a.Act

	if !strings.Contains(act.ActionName, "transfer") {
		return
	}

	var data Transfer
	if err := json.Unmarshal(act.Data, &data); err != nil {
		log.Printf("UNmarshal json err : %v", err)
		return
	}

	//just for test
	if _, ok := a.Controller.TestAccounts[data.From]; !ok {
		var account UserAccount
		account.Name = data.From
		account.Amount = 100

		a.Controller.TestAccounts[data.From] = &account
	}

	if _, ok := a.Controller.TestAccounts[data.To]; !ok {
		var account UserAccount
		account.Name = data.To
		account.Amount = 100

		a.Controller.TestAccounts[data.To] = &account
	}
	//just for test

	accounts := a.Controller.TestAccounts
	from := accounts[data.From]
	to := accounts[data.To]

	if err := from.SubBalance(data.Amount); err != nil {
		return
	}

	to.AddBalance(data.Amount)
}
