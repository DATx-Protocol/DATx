package types

import (
	"datx_chain/utils/common"
	"datx_chain/utils/helper"
	time "time"
)

//GenesisState struct
type GenesisState struct {
	InitTimeStamp time.Time

	InitKey string
}

//NewGenesisState new
func NewGenesisState() *GenesisState {
	var res GenesisState
	res.InitTimeStamp = ToTime("2018-07-01T12:00:00.0Z")
	res.InitKey = RootKey

	return &res
}

//ComputeChainID compute chain id
func (gs *GenesisState) ComputeChainID() common.Hash {
	return helper.RLPHash(gs)
}
