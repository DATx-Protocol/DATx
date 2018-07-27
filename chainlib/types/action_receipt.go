package types

import "datx_chain/utils/common"

//ActionReceipt struct
type ActionReceipt struct {
	Receiver  string
	ActDigest common.Hash
}
