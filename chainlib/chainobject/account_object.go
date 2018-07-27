package chainobject

import (
	"datx_chain/utils/common"
	"time"
)

type AccountObject struct {
	//the global id of account object
	ID uint64

	//account name
	Name string

	//vm type
	Vm_type uint8

	//vm version
	Vm_version uint8

	//
	Privileged bool

	//the time of updating smart contract
	Last_code_update time.Time

	//the hash code of smart contract
	Code_version common.Hash

	//the creation time of smart contract
	Creation_date time.Time

	//the code of smart contract
	Code []byte

	//the abi of smart contract
	Abi []byte
}

func NewAccount(name string) *AccountObject {
	return &AccountObject{
		ID:         GetOID(AccountType),
		Name:       name,
		Vm_type:    0,
		Vm_version: 0,
		Privileged: false,
	}
}

func (self *AccountObject) SetCode(code []byte) {
	self.Code = code
	// self.code_version = utils.RLPHash(self.code)
	self.Creation_date = time.Now()
	self.Last_code_update = time.Now()
}

func (self *AccountObject) SetAbi(abi []byte) {
	self.Abi = abi
}
