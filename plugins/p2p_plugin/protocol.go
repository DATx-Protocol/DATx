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
	"datx_chain/utils/common"
	"time"
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
	MsgTypeStatusMsg         = 0x00
	MsgTypeAnnounceMsg       = 0x01 //获取chain_size_message
	MsgTypeTimeMsg           = 0x02 //
	MsgTypeNoticeMsg         = 0x03
	MsgTypeRequestMsg        = 0x04
	MsgTypeSyncRequestMsg    = 0x05
	MsgTypeSignedBlock       = 0x06
	MsgTypePackedTransaction = 0x07
	MsgTypeGoAwayMsg         = 0x08
)

type errCode uint

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

const (
	SyncStageLibCatchup = iota
	SyncStageHeadCatchup
	SyncStageInSync
)

func (e errCode) String() string {
	return errorToString[int(e)]
}

// XXX change once legacy code is out
var errorToString = map[int]string{
	ErrNoReason:           "no reason",
	ErrConnectSelf:        "self connect",
	ErrDuplicate:          "duplicate",
	ErrInvalidNetwork:     "wrong chain",
	ErrInValidCersion:     "wrong version",
	ErrForkedChain:        "chain is forked",
	ErrUnLinkable:         "unlinkable block received",
	ErrInvalidTransaction: "bad transaction",
	ErrInvalidBlock:       "invalid block",
	ErrBeningOther:        "some other non-fatal condition",
	ErrOtherFatal:         "some other failure",
	ErrAuthentication:     "authentication failure",
}

type ChainSizeMsg struct {
	LastIrreversibleBlockNum uint32
	LastIrreversibleBlockId  common.Hash
	HeadNum                  uint32
	HeadId                   common.Hash
}

type HandShakeMsg struct {
	NetworkVersion           uint16
	NetworkId                uint32
	TimeStamp                time.Time
	Token                    common.Hash
	Sigature                 common.Hash
	LastIrreversibleBlockNum uint32
	LastIrreversibleBlockId  common.Hash
	HeadNum                  uint32
	HeadId                   common.Hash
	OS                       string
	Agent                    string
	Generation               uint16
	PeerID                   string
}

type GoAwayMsg struct {
	Reason uint32
}

type IdListMode uint16

const (
	None = iota
	CatchUp
	LastIrrCatchUp
	Normal
)

type SelectIds struct {
	Mode    IdListMode
	Pending uint32
	Ids     []interface{}
}

func (s *SelectIds) Empty() bool {
	return s.Mode == None || len(s.Ids) == 0
}

type DisconnectMsg struct {
	Reason errCode
}

type TimeMsg struct {
	OrginTime       time.Time
	RecvTime        time.Time
	TransmitTime    time.Time
	DestinationTime time.Time
}

type NoticeMsg struct {
	TxIds  SelectIds
	BlkIds SelectIds
}

type RequestMsg struct {
	TxIds  SelectIds
	BlkIds SelectIds
}

type SyncRequestMsg struct {
	StartBlkNum uint32
	EndBlkNum   uint32
}
