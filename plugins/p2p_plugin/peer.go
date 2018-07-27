package p2p_plugin

import (
	"crypto/ecdsa"
	"errors"
	"sync"
	"time"

	"datx_chain/chainlib/types"
	"datx_chain/plugins/p2p_plugin/p2p"
	"datx_chain/utils/common"
)

var (
	errClosed            = errors.New("peer set is closed")
	errAlreadyRegistered = errors.New("peer is already registered")
	errNotRegistered     = errors.New("peer is not registered")
)

type TransactionState struct {
	Id         common.Hash
	ExpireTime time.Time
	BlockNum   uint32
}

type SyncBlockState struct {
	StartBlockNum    uint32
	EndBlockNum      uint32
	LastSyncBlockNum uint32
	StartSyncTime    time.Time
	Syncing          bool
	Connecting       bool
	ForkHead         types.BlockState
	ForkHeadNum      uint32
}

type peerBlockState struct {
	Id            common.Hash
	BlkNum        uint32
	IsKnown       bool
	IsNoticed     bool
	RequestedTime time.Time
}

type transactionState struct {
	Id              common.Hash
	IsKnownByPeer   bool   ///< true if we sent or received this trx to this peer or received notice from peer
	IsNoticedToPeer bool   ///< have we sent peer notice we know it (true if we receive from this peer)
	BlkNum          uint32 ///< the block number the transaction was included in
	Expries         uint32
	RequestedTime   time.Time
}

type peer struct {
	id string
	*p2p.Peer
	pubKey *ecdsa.PublicKey

	rw p2p.MsgReadWriter

	version int    // Protocol version negotiated
	network uint64 // Network ID being on

	lastHandshakeRecv  HandShakeMsg
	lastHandshakeSent  HandShakeMsg
	sentHandshakeCount uint32

	syncing       bool
	peerRequested syncState

	forkHead    common.Hash
	forkHeadNum uint32

	blkState map[common.Hash]peerBlockState
	trxState map[common.Hash]transactionState

	lastReq RequestMsg
}

type syncState struct {
	startBlock uint32
	endBlock   uint32
	last       uint32
	timePoint  time.Time
}

func newPeer(version int, network uint64, p *p2p.Peer, rw p2p.MsgReadWriter) *peer {
	id := p.ID()
	pubKey, _ := id.Pubkey()

	return &peer{
		id:      p.ID().String(),
		Peer:    p,
		pubKey:  pubKey,
		rw:      rw,
		version: version,
		network: network,

		blkState: make(map[common.Hash]peerBlockState, 0),
		trxState: make(map[common.Hash]transactionState, 0),
	}
}

type peerSet struct {
	peers  map[string]*peer
	lock   sync.RWMutex
	closed bool
}

func newPeerSet() *peerSet {
	return &peerSet{
		peers: make(map[string]*peer),
	}
}

// Peer retrieves the registered peer with the given id.
func (ps *peerSet) Peer(id string) *peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return ps.peers[id]
}

// Len returns if the current number of peers in the set.
func (ps *peerSet) Len() int {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return len(ps.peers)
}

// Register injects a new peer into the working set, or returns an error if the
// peer is already known.
func (ps *peerSet) Register(p *peer) error {
	ps.lock.Lock()
	if ps.closed {
		ps.lock.Unlock()
		return errClosed
	}
	if _, ok := ps.peers[p.id]; ok {
		ps.lock.Unlock()
		return errAlreadyRegistered
	}
	ps.peers[p.id] = p
	ps.lock.Unlock()

	return nil
}

// Unregister removes a remote peer from the active set, disabling any further
// actions to/from that particular entity. It also initiates disconnection at the networking layer.
func (ps *peerSet) Unregister(id string) error {
	ps.lock.Lock()
	if p, ok := ps.peers[id]; !ok {
		ps.lock.Unlock()
		return errNotRegistered
	} else {
		delete(ps.peers, id)
		ps.lock.Unlock()

		p.Peer.Disconnect(p2p.DiscUselessPeer)
		return nil
	}
}

// AllPeerIDs returns a list of all registered peer IDs
func (ps *peerSet) AllPeerIDs() []string {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	res := make([]string, len(ps.peers))
	idx := 0
	for id := range ps.peers {
		res[idx] = id
		idx++
	}
	return res
}

func (ps *peerSet) Close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.peers {
		p.Disconnect(p2p.DiscQuitting)
	}
	ps.closed = true
}

func (p *peer) Reset() {
	p.peerRequested = syncState{}
	p.blkState = make(map[common.Hash]peerBlockState, 0)
	p.trxState = make(map[common.Hash]transactionState, 0)
}
