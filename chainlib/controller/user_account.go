package controller

import "errors"

//UserAccount user
type UserAccount struct {
	Name   string
	Amount uint16
}

//AddBalance add
func (user *UserAccount) AddBalance(am uint16) {
	user.Amount = user.Amount + am
}

//SubBalance sub
func (user *UserAccount) SubBalance(am uint16) error {
	if user.Amount < am {
		return errors.New("the token is insufficient")
	}

	user.Amount = user.Amount - am
	return nil
}
