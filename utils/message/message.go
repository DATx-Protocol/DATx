package message

import (
	"datx_chain/utils/rlp"
	"fmt"
	"io"
)

//defines msg code
const (
	BlockMsg uint16 = 1 << iota
	BlockHeaderMsg
	TrxMsg
)

//Msg defines the message structure of chain_plugin and p2p_plugin
type Msg struct {
	//msg code
	Code uint16

	Size uint32 //size of the payload

	//byte array of struct
	Payload io.Reader
}

//New Msg with given MsgCode and Msg Data
func NewMsg(code uint16, data interface{}) *Msg {
	size, r, err := rlp.EncodeToReader(data)
	if err != nil {
		return nil
	}

	return &Msg{Code: code, Size: uint32(size), Payload: r}
}

//Send msg to given chan
func (self *Msg) Send(c chan<- *Msg) {
	c <- self
}

// Decode parses the RLP content of a message into
// the given value, which must be a pointer.
//
// For the decoding rules, please see package rlp.
func (self *Msg) Decode(val interface{}) error {
	s := rlp.NewStream(self.Payload, uint64(self.Size))
	if err := s.Decode(val); err != nil {
		return fmt.Errorf("Invalid massage (code %x) (size %d) %v", self.Code, self.Size, err)
	}
	return nil
}

//
func (self *Msg) String() string {
	return fmt.Sprintf("msg #%v (%v bytes)", self.Code, self.Size)
}
