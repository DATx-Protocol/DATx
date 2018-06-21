// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package les implements the Light Ethereum Subprotocol.
package p2p_plugin

import (
	"time"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"datx_chain/utils/common"
	"datx_chain/utils/rlp"
)


// Constants to match up protocol versions and messages
const (
	lpv1 = 1
	lpv2 = 2
)

// Supported versions of the les protocol (first is primary)
var (
	ClientProtocolVersions    = []uint{lpv2, lpv1}
	ServerProtocolVersions    = []uint{lpv2, lpv1}
	AdvertiseProtocolVersions = []uint{lpv2} // clients are searching for the first advertised protocol in the list
)

// Number of implemented message corresponding to different protocol versions.
var ProtocolLengths = map[uint]uint64{lpv1: 15, lpv2: 22}

const (
	NetworkId          = 1
	ProtocolMaxMsgSize = 10 * 1024 * 1024 // Maximum cap on the size of a protocol message
)

const (
	StatusMsg         = 0x00
	AnnounceMsg       = 0x01 //获取chain_size_message
	TimeMsg           = 0x02 //
	NoticeMsg         = 0x03
	RequestMsg        = 0x04
	SyncRequestMsg    = 0x05
	SignedBlock       = 0x06
	PackedTransaction = 0x07
)

type errCode int

const (
	ErrNoReason = iota
	ErrConnectSelf
	ErrDuplicate
	ErrInvalidNetwork 
	ErrInValidCersion
	ErrForkedChain
	ErrUnLinkable
	ErrInvalidTransaction
	ErrInvalidBlock
	ErrBeningOther
	ErrOtherFatal
	ErrAuthentication
)

func (e errCode) String() string {
	return errorToString[int(e)]
}

// XXX change once legacy code is out
var errorToString = map[int]string{
	ErrNoReason:"no reason",
	ErrConnectSelf:"self connect",
	ErrDuplicate:"duplicate",
	ErrInvalidNetwork:"wrong chain",
	ErrInValidCersion:"wrong version",
	ErrForkedChain:"chain is forked",
	ErrUnLinkable:"unlinkable block received",
	ErrInvalidTransact:"bad transaction",
	ErrInvalidBlock:"invalid block",
	ErrBeningOther:"some other non-fatal condition",
	ErrOtherFatal:"some other failure",
	ErrAuthentication:"authentication failure"
}


type ChainSizeMsg struct {
	lastIrreversibleBlockNum uint32  ,
	lastIrreversibleBlockId common.Hash,
	headNum uint32,
	headId common.Hash

}

type HandShakeMsg struct {
	networkVersion uint16
	networkId	uint32
	timeStamp	time.Time
	token		common.Hash
	sigature	common.Hash
	lastIrreversibleBlockNum	uint32
	lastIrreversibleBlockId	common.Hash
	headNum	uint32
	headId	common.Hash
	os	string
	agent string
	generation  int16

}

type SelectIds struct {
	type IdListMode int16 
	const (
		None = iota
		CatchUp
		lastIrrCatchUp
		Normal
	)

	mode IdListMode,
	pending uint32,
	ids [] struct {}
}

(s *SelectIds) func () bool {
	return s.mode == SelectIds.None || len(s.ids) == 0 
}

type DisconnectMsg struct {
	reason errCode
}

type TimeMsg struct {
	orginTime time.Time,
	recvTime  time.Time,
	transmitTime  time.Time,
	destinationTime time.Time
}

type NoticeMsg struct {
	txIds SelectIds,
	blkIds selectIds,
}

type RequestMsg struct{
	txIds SelectIds,
	blkIds SelectIds
}

type SyncRequestMsg struct {
	startBlkNum uint32,
	endBlkNum uint32
}
