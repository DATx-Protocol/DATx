package chainobject

import (
	"datx_chain/utils/common"
	"datx_chain/utils/helper"
	"time"
)

type AccountObject struct {
	//the global id of account object
	id uint64

	//account name
	name string

	//vm type
	vm_type uint8

	//vm version
	vm_version uint8

	//
	privileged bool

	//the time of updating smart contract
	last_code_update time.Time

	//the hash code of smart contract
	code_version common.Hash

	//the creation time of smart contract
	creation_date time.Time

	//the code of smart contract
	code []byte

	//the abi of smart contract
	abi []byte
}

func NewAccount(name string) *AccountObject {
	return &AccountObject{
		id:         helper.GetOID(),
		name:       name,
		vm_type:    0,
		vm_version: 0,
		privileged: false,
	}
}

func (self *AccountObject) ID() uint64 {
	return self.id
}

func (self *AccountObject) SetCode(code []byte) {
	self.code = code
	// self.code_version = utils.RLPHash(self.code)
	self.creation_date = time.Now()
	self.last_code_update = time.Now()
}

func (self *AccountObject) SetAbi(abi []byte) {
	self.abi = abi
}
