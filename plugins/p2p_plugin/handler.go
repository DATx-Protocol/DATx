package p2p_plugin

import (
	"errors"
	"fmt"
	slog "log"
	"reflect"
	"sync"
	"time"

	"datx_chain/chainlib/application"
	"datx_chain/chainlib/controller"
	"datx_chain/chainlib/types"
	"datx_chain/plugins/chain_plugin"
	"datx_chain/plugins/p2p_plugin/p2p"
	"datx_chain/plugins/p2p_plugin/p2p/discover"
	"datx_chain/utils/common"
	"datx_chain/utils/db"
	"datx_chain/utils/log"

	"github.com/robfig/cron"
)

var errIncompatibleConfig = errors.New("incompatible configuration")

func errResp(code errCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}

type ProtocolManager struct {
	networkId uint64
	peers     *peerSet
	maxPeers  int

	SubProtocols []p2p.Protocol
	chainDb      *datxdb.LDBDatabase
	// channels for fetcher, syncer, txsyncLoop
	newPeerCh   chan *peer
	quitSync    chan struct{}
	noMorePeers chan struct{}

	// wait group is used for graceful shutdowns during downloading
	// and processing
	wg *sync.WaitGroup

	syncMaster     *SyncManager
	dispatchMaster *DispatchManager

	localTxns map[common.Hash]nodeTranscationState
}

type SyncManager struct {
	state                int
	syncKnownlibNum      uint32
	syncLastRequestedNum uint32
	syncNextExpectedNum  uint32
	syncReqSpan          uint32

	source *peer

	chain *controller.Controller
}

type DispatchManager struct {
	justSendit uint32
	reqBlks    map[common.Hash]blockRequest
	reqTrx     map[common.Hash]struct{}

	receivedBlks map[common.Hash]blockOrigin
	receivedTrx  map[common.Hash]transactionOrigin
}

type blockRequest struct {
	id         common.Hash
	localRetry bool
}

type blockOrigin struct {
	id     common.Hash
	origin *peer
}

type transactionOrigin struct {
	id     common.Hash
	origin *peer
}

type nodeTranscationState struct {
	id        common.Hash
	expires   uint32
	packedTxn types.PackedTransaction
	blkNum    uint32
	trueBlk   uint32
	requests  uint16
}

func (nts *nodeTranscationState) updateInFlight(incr int32) {
	exp := int32(nts.expires)
	nts.expires = uint32(exp + incr*60)
	if nts.requests == 0 {
		nts.trueBlk = nts.blkNum
		nts.blkNum = 0
	}
	nts.requests = nts.requests + uint16(incr)
	if nts.requests == 0 {
		nts.blkNum = nts.trueBlk
	}
}

var P2pPlugin_impl *P2p_Plugin

// NewProtocolManager returns a new ethereum sub protocol manager. The Ethereum sub protocol manages peers capable
// with the ethereum network.
func NewProtocolManager(protocolVersions []uint, networkId uint64, peers *peerSet, chainDb *datxdb.LDBDatabase, quitSync chan struct{}, wg *sync.WaitGroup) (*ProtocolManager, error) {
	// Create the protocol manager with the base fields
	manager := &ProtocolManager{
		chainDb:     chainDb,
		networkId:   networkId,
		peers:       peers,
		newPeerCh:   make(chan *peer),
		quitSync:    quitSync,
		wg:          wg,
		noMorePeers: make(chan struct{}),

		syncMaster: &SyncManager{
			state:                SyncStageInSync,
			syncKnownlibNum:      0,
			syncLastRequestedNum: 0,
			syncNextExpectedNum:  1,
			syncReqSpan:          1,
		},
		dispatchMaster: &DispatchManager{
			reqBlks: make(map[common.Hash]blockRequest, 0),
			reqTrx:  make(map[common.Hash]struct{}, 0),

			receivedBlks: make(map[common.Hash]blockOrigin, 0),
			receivedTrx:  make(map[common.Hash]transactionOrigin, 0),
		},

		localTxns: make(map[common.Hash]nodeTranscationState, 0),
	}
	plugin, err := application.App().Find("chain")
	if err == nil {
		if chainplugin, ok := plugin.(*chainplugin.ChainPlugin); ok {
			manager.syncMaster.chain = chainplugin.Chain()
		}
	}

	// Initiate a sub-protocol for every implemented version we can handle
	manager.SubProtocols = make([]p2p.Protocol, 0, len(protocolVersions))
	for _, version := range protocolVersions {
		// Compatible, initialize the sub-protocol
		version := version // Closure for the run
		manager.SubProtocols = append(manager.SubProtocols, p2p.Protocol{
			Name:    "datx",
			Version: version,
			Length:  ProtocolLengths[version],
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				peer := manager.newPeer(int(version), networkId, p, rw)
				select {
				case manager.newPeerCh <- peer:
					manager.wg.Add(1)
					defer manager.wg.Done()
					err := manager.handle(peer)
					return err
				case <-manager.quitSync:
					return p2p.DiscQuitting
				}
			},
			NodeInfo: func() interface{} {
				return ""
			},
			PeerInfo: func(id discover.NodeID) interface{} {
				if p := manager.peers.Peer(fmt.Sprintf("%x", id[:8])); p != nil {
					return p.Info()
				}
				return nil
			},
		})
	}
	if len(manager.SubProtocols) == 0 {
		return nil, errIncompatibleConfig
	}

	return manager, nil
}

// removePeer initiates disconnection from a peer by removing it from the peer set
func (pm *ProtocolManager) removePeer(id string) {
	pm.peers.Unregister(id)
}

func (pm *ProtocolManager) Start(maxPeers int) {
	pm.maxPeers = maxPeers
	go func() {
		for range pm.newPeerCh {
		}
	}()
	c := cron.New()
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		pm.expireTxns()
	})
	c.Start()
}

func (pm *ProtocolManager) Stop() {
	// Showing a log message. During download / process this could actually
	// take between 5 to 10 second3s and therefor feedback is required.

	// Quit the sync loop.
	// After this send has completed, no new peers will be accepted.
	pm.noMorePeers <- struct{}{}

	close(pm.quitSync) // quits syncer, fetcher

	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the peer set.
	// sessions which are already established but not added to pm.peers yet
	// will exit when they try to register.
	pm.peers.Close()

	// Wait for any process action
	pm.wg.Wait()

}

func (pm *ProtocolManager) newPeer(pv int, nv uint64, p *p2p.Peer, rw p2p.MsgReadWriter) *peer {
	return newPeer(pv, nv, p, rw)
}

// handle is the callback invoked to manage the life cycle of a les peer. When
// this function terminates, the peer is disconnected.
func (pm *ProtocolManager) handle(p *peer) error {
	// Ignore maxPeers if this is a trusted peer
	if pm.peers.Len() >= pm.maxPeers && !p.Peer.Info().Network.Trusted {
		return p2p.DiscTooManyPeers
	}

	p.Log().Debug("datx peer connected", "name", p.Name())

	// Register the peer locally
	if err := pm.peers.Register(p); err != nil {
		p.Log().Error("datx peer registration failed", "err", err)
		return err
	}
	defer pm.removePeer(p.id)
	p.sendHandshake()
	// main loop. handle incoming messages.
	for {
		if err := pm.handleMsg(p); err != nil {
			p.Log().Debug("Ethereum message handling failed", "err", err)
			return err
		}
	}

	return nil
}

// handleMsg is invoked whenever an inbound message is received from a remote
// peer. The remote connection is torn down upon returning any error.
func (pm *ProtocolManager) handleMsg(p *peer) error {
	// Read the next message from the remote peer, and ensure it's fully consumed
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return err
	}
	p.Log().Trace("datx message arrived", "code", msg.Code, "bytes", msg.Size)

	switch {
	case msg.Code == MsgTypeStatusMsg:
		slog.Printf("received handshake_message")
		HSMsg := HandShakeMsg{}
		msg.Decode(&HSMsg)
		if !pm.isValid(p, &HSMsg) {
			p2p.Send(p.rw, MsgTypeGoAwayMsg, GoAwayMsg{Reason: ErrOtherFatal})
			return nil
		}

		if HSMsg.Generation == 1 {
			if HSMsg.PeerID == P2pPlugin_impl.Status().ID {
				p.Log().Trace("Self connection detected. Closing connection")
				p2p.Send(p.rw, MsgTypeGoAwayMsg, GoAwayMsg{Reason: ErrConnectSelf})
				return nil
			}

			if p.sentHandshakeCount == 0 {
				p.sendHandshake()
			}
		}
		p.lastHandshakeRecv = HSMsg
		pm.syncMaster.recvHandshake(p, &HSMsg)
		break
	case msg.Code == MsgTypeAnnounceMsg:
		slog.Printf("received announce_message")
		break
	case msg.Code == MsgTypeTimeMsg:
		//todo
		break
	case msg.Code == MsgTypeNoticeMsg:
		// peer tells us about one or more blocks or txns. When done syncing, forward on
		// notices of previously unknown blocks or txns,
		slog.Printf("received notice_message")
		NotcMsg := NoticeMsg{}
		msg.Decode(&NotcMsg)

		request := RequestMsg{}
		sendReq := bool(false)
		if NotcMsg.TxIds.Mode != None {
			p.Log().Trace("this is a %s notice with %d blocks", NotcMsg.TxIds.Mode, NotcMsg.TxIds.Pending)
		}
		switch NotcMsg.TxIds.Mode {
		case None:
			break
		case LastIrrCatchUp:
			p.lastHandshakeRecv.HeadNum = NotcMsg.TxIds.Pending
			request.TxIds.Mode = None
			break
		case CatchUp:
			if NotcMsg.TxIds.Pending > 0 {
				// plan to get all except what we already know about.
				request.TxIds.Mode = CatchUp
				sendReq = true
				knownSum := len(pm.localTxns)
				if knownSum != 0 {
					request.TxIds.Ids = make([]interface{}, 0)
					for key := range pm.localTxns { //syncMaster.chain.UnAppliedTransaction
						request.TxIds.Ids = append(request.TxIds.Ids, key)
					}
				}
			}
			break
		case Normal:
			pm.dispatchMaster.recvNotice(p, &NotcMsg, false)
			break
		}

		if NotcMsg.BlkIds.Mode != None {
			p.Log().Trace("this is a %s notice with %d blocks", NotcMsg.BlkIds.Mode, NotcMsg.BlkIds.Pending)
		}
		switch NotcMsg.BlkIds.Mode {
		case None:
			if NotcMsg.TxIds.Mode != Normal {
				return nil
			}
			break
		case LastIrrCatchUp:
		case CatchUp:
			pm.syncMaster.recvNotice(p, &NotcMsg)
			break
		case Normal:
			pm.dispatchMaster.recvNotice(p, &NotcMsg, false)
			break
		default:
			p.Log().Trace("bad notice_message : invalid known_blocks.mode %s", NotcMsg.BlkIds.Mode)
		}
		if sendReq {
			p2p.Send(p.rw, MsgTypeRequestMsg, request)
		}
		break
	case msg.Code == MsgTypeRequestMsg:
		slog.Printf("received request_message")
		ReqMsg := RequestMsg{}
		msg.Decode(&ReqMsg)
		switch ReqMsg.BlkIds.Mode {
		case CatchUp:
			p.Log().Trace("received request_message:catch_up")
			p.blkSendBranch()
			break
		case Normal:
			p.Log().Trace("received request_message:normal")
			p.blkSend(ReqMsg.BlkIds.Ids)
			break
		default:
			break
		}

		switch ReqMsg.TxIds.Mode {
		case CatchUp:
			p.txnSendPending(ReqMsg.TxIds.Ids)
			break
		case Normal:
			p.txnSend(ReqMsg.TxIds.Ids)
			break
		case None:
			if ReqMsg.BlkIds.Mode == None {
				p.stopSend()
			}
			break
		default:
			break
		}
		break
	case msg.Code == MsgTypeSyncRequestMsg:
		slog.Printf("received sync_request_msg")
		syncReqMsg := SyncRequestMsg{}
		msg.Decode(&syncReqMsg)
		if syncReqMsg.EndBlkNum == 0 {
			p.peerRequested = syncState{}
		} else {
			p.peerRequested = syncState{
				startBlock: syncReqMsg.StartBlkNum,
				endBlock:   syncReqMsg.EndBlkNum,
				last:       syncReqMsg.StartBlkNum - 1,
			}
			p.enqueueSyncBlock()
		}
		break
	case msg.Code == MsgTypeSignedBlock:

		blkData := types.NewBlockStateEmpty()
		msg.Decode(blkData)
		blkId := blkData.ID
		blkNum := blkData.BlockNum
		slog.Printf("received signed_block,%d", blkNum)
		if pm.syncMaster.chain.ForkDB.GetBlock(blkId) != nil {
			pm.syncMaster.recvBlock(p, blkId, blkNum)
			return nil
		}
		pm.dispatchMaster.recvBlock(p, blkId, blkNum)

		reason := ErrOtherFatal
		// panic catch?
		pm.syncMaster.chain.AcceptBlock(blkData.Block)
		reason = ErrNoReason
		if reason == ErrNoReason {
			for _, val := range blkData.Trxs {
				id := val.ID
				ltx, exist := P2pPlugin_impl.pm.localTxns[id]
				if exist {
					ltx.blkNum = blkNum
					P2pPlugin_impl.pm.localTxns[id] = ltx
				}

				ptx, exist := p.trxState[id]
				if exist {
					ptx.BlkNum = blkNum
					p.trxState[id] = ptx
				}
			}
			pm.syncMaster.recvBlock(p, blkId, blkNum)
		} else {
			pm.syncMaster.rejectedBlock(p, blkNum)
		}

		break
	case msg.Code == MsgTypePackedTransaction:
		slog.Printf("received packed_trx")
		packedTxn := types.PackedTransaction{}
		msg.Decode(&packedTxn)
		p.Log().Trace("got a packed transaction")
		if pm.syncMaster.isActive(p) {
			p.Log().Trace("got a txn during sync - dropping")
			return nil
		}
		tid := packedTxn.ID()
		if _, exist := P2pPlugin_impl.pm.localTxns[tid]; exist {
			p.Log().Trace("got a duplicate transaction - dropping")
			return nil
		}
		pm.dispatchMaster.recvTransaction(p, tid)
		pm.syncMaster.chain.AcceptTransaction(&packedTxn, func(inerr error, trace *types.TransactionTrace) {
			if inerr == nil {
				//broadcast the txn if accepted
				pm.dispatchMaster.bcastTransaction(&packedTxn)
			} else {
				// reject it
				pm.dispatchMaster.rejectedTransaction(tid)
			}
		})
		break
	case msg.Code == MsgTypeGoAwayMsg:
		slog.Printf("received go_away_message")
		pm.peers.Unregister(p.id)
		break
	default:
		slog.Printf("invalid Msg Type , code(%d), bytes(%d)", msg.Code, msg.Size)
	}

	return nil
}

func (pm *ProtocolManager) sendAll(msgcode uint64, data interface{}, f func(*peer) bool) {
	for _, p := range pm.peers.peers {
		if !p.syncing && f(p) {
			p2p.Send(p.rw, msgcode, data)
		}
	}
}

func (pm *ProtocolManager) blockLoop() {
	pm.wg.Add(1)
	go func() {
		for {
			select {
			case <-pm.quitSync:
				pm.wg.Done()
				return
			}
		}
	}()
}

// Do some basic validation of an incoming handshake_message, so things
// that really aren't handshake messages can be quickly discarded without
// affecting state.
func (pm *ProtocolManager) isValid(p *peer, Msg *HandShakeMsg) bool {
	if Msg.LastIrreversibleBlockNum > Msg.HeadNum {
		p.Log().Trace("Handshake message validation: last irreversible block %d is greater than head block %d", Msg.LastIrreversibleBlockNum, Msg.HeadNum)
		return false
	}
	// if Msg.os == "" {
	// 	p.Log().Trace("Handshake message validation: os field is null string")
	// 	return false
	// }
	//other checks
	return true
}

func (pm *ProtocolManager) expireTxns() {
	nowSec := time.Now().Unix()
	bn := pm.syncMaster.chain.Head.DposIrreversibleBlockNum
	txsToDel := make([]common.Hash, 0)
	for _, val := range pm.localTxns {
		if int64(val.expires) <= nowSec {
			txsToDel = append(txsToDel, val.id)
		}

		if val.blkNum <= bn {
			txsToDel = append(txsToDel, val.id)
		}
	}

	for _, val := range txsToDel {
		delete(pm.localTxns, val)
	}

	for _, p := range pm.peers.peers {
		txsToDel = make([]common.Hash, 0)
		for _, val := range p.trxState {
			if int64(val.Expries) <= nowSec {
				txsToDel = append(txsToDel, val.Id)
			}
			if val.BlkNum < bn {
				txsToDel = append(txsToDel, val.Id)
			}
			for _, val := range txsToDel {
				delete(pm.peers.peers[p.id].trxState, val)
			}
		}

		blksToDel := make([]common.Hash, 0)
		for _, val := range p.blkState {
			if val.BlkNum < bn {
				blksToDel = append(blksToDel, val.Id)
			}
			for _, val := range blksToDel {
				delete(pm.peers.peers[p.id].blkState, val)
			}
		}
	}
}

/*-------------------SyncManager--------------------------------------------*/
func (sm *SyncManager) recvHandshake(p *peer, Msg *HandShakeMsg) {
	libNum := sm.chain.Head.DposIrreversibleBlockNum
	peerLib := Msg.LastIrreversibleBlockNum

	headNum := sm.chain.Head.BlockNum
	headID := sm.chain.Head.ID

	sm.resetLibNum(p)
	p.syncing = false

	//--------------------------------
	// sync need checks; (lib == last irreversible block)
	//
	// 0. my head block id == peer head id means we are all caugnt up block wise
	// 1. my head block num < peer lib - start sync locally
	// 2. my lib > peer head num - send an last_irr_catch_up notice if not the first generation
	//
	// 3  my head block num <= peer head block num - update sync state and send a catchup request
	// 4  my head block num > peer block num ssend a notice catchup if this is not the first generation
	//
	//-----------------------------
	if reflect.DeepEqual(headID, Msg.HeadId) {
		p.Log().Trace("sync check state 0")
		note := NoticeMsg{
			BlkIds: SelectIds{
				Mode: None,
			},
			TxIds: SelectIds{
				Mode:    CatchUp,
				Pending: uint32(len(P2pPlugin_impl.pm.localTxns)),
			},
		}
		p2p.Send(p.rw, MsgTypeNoticeMsg, note)
		return
	}
	if headNum < peerLib {
		p.Log().Trace("sync check state 1")
		sm.startSync(p, peerLib)
		return
	}
	if libNum > Msg.HeadNum {
		p.Log().Trace("sync check state 2")
		if Msg.Generation > 1 {
			note := NoticeMsg{
				BlkIds: SelectIds{
					Mode:    LastIrrCatchUp,
					Pending: headNum,
				},
				TxIds: SelectIds{
					Mode:    LastIrrCatchUp,
					Pending: libNum,
				},
			}
			p2p.Send(p.rw, MsgTypeNoticeMsg, note)
		}
		p.syncing = true
		return
	}
	if headNum < Msg.HeadNum {
		p.Log().Trace("sync check state 3")
		sm.verifyCatchup(p, Msg.HeadNum, Msg.HeadId)
		return
	} else {
		p.Log().Trace("sync check state 4")
		if Msg.Generation > 1 {
			note := NoticeMsg{
				BlkIds: SelectIds{
					Mode:    CatchUp,
					Pending: headNum,
					Ids:     []interface{}{headID},
				},
				TxIds: SelectIds{
					Mode: None,
				},
			}
			p2p.Send(p.rw, MsgTypeNoticeMsg, note)
		}
		p.syncing = true
		return
	}
	p.Log().Trace("sync check state 4")
}

func (sm *SyncManager) startSync(p *peer, target uint32) {
	if target > sm.syncKnownlibNum {
		sm.syncKnownlibNum = target
	}
	if !sm.syncRequired() {
		bnum := sm.chain.Head.DposIrreversibleBlockNum
		hnum := sm.chain.Head.BlockNum
		p.Log().Trace("We are already caught up, my irr = %d, head = %d, target = %d", bnum, hnum, target)
		return
	}
	if sm.state == SyncStageInSync {
		sm.setState(SyncStageLibCatchup)
		sm.syncNextExpectedNum = sm.chain.Head.DposIrreversibleBlockNum + 1
		p.Log().Trace("Catching up with chain, our last req is %d, theirs is %d peer %s", sm.syncLastRequestedNum, target, p.id)
	}
	sm.requestNextChunk(p)
}

func (sm *SyncManager) syncRequired() bool {
	return (sm.syncLastRequestedNum < sm.syncKnownlibNum ||
		sm.chain.Head.BlockNum < sm.syncLastRequestedNum)
}

func (sm *SyncManager) setState(newState int) {
	if sm.state == newState {
		return
	}
	sm.state = newState
}

func (sm *SyncManager) requestNextChunk(p *peer) {
	headBlockNum := sm.chain.Head.BlockNum

	if headBlockNum < sm.syncLastRequestedNum {
		p.Log().Trace("ignoring request, head is %d last req = %d source is %s", headBlockNum, sm.syncLastRequestedNum, p.id)
		return
	}
	/* ----------
	 *Eos Logic
	 * next chunk provider selection criteria
	 * a provider is supplied and able to be used, use it.
	 * otherwise select the next available from the list, round-robin style.
	 * just use P for now
	 */
	if !p.syncing {
		sm.source = p
	} else {
		if len(P2pPlugin_impl.peers.peers) == 1 && sm.source == nil {
			for _, val := range P2pPlugin_impl.peers.peers {
				sm.source = val
			}
		} else {
			for _, val := range P2pPlugin_impl.peers.peers {
				if !val.syncing {
					sm.source = val
					break
				}
			}
			// verify there is an available source
			if sm.source == nil || sm.source.syncing {
				p.Log().Trace("Unable to continue syncing at this time")
				sm.syncKnownlibNum = sm.chain.Head.DposIrreversibleBlockNum
				sm.syncLastRequestedNum = 0
				sm.setState(SyncStageInSync) // probably not, but we can't do anything else
				return
			}
		}
	}
	if sm.syncLastRequestedNum != sm.syncKnownlibNum {
		start := sm.syncNextExpectedNum
		end := start + sm.syncReqSpan - 1
		if end > sm.syncKnownlibNum {
			end = sm.syncKnownlibNum
		}
		if end > 0 && end >= start {
			p.Log().Trace("requesting range %d to %d, from %s", start, end, p.id)
			note := SyncRequestMsg{
				StartBlkNum: start,
				EndBlkNum:   end,
			}
			p2p.Send(p.rw, MsgTypeSyncRequestMsg, note)
			sm.syncLastRequestedNum = end
		}
	}
}
func (sm *SyncManager) verifyCatchup(p *peer, num uint32, id common.Hash) {
	req := RequestMsg{}
	req.BlkIds.Mode = CatchUp
	for _, p := range P2pPlugin_impl.peers.peers {
		if reflect.DeepEqual(p.forkHead, id) || p.forkHeadNum > num {
			req.BlkIds.Mode = None
		}
		break
	}
	if req.BlkIds.Mode == CatchUp {
		p.forkHead = id
		p.forkHeadNum = num
		if sm.state == SyncStageLibCatchup {
			return
		}
		sm.setState(SyncStageHeadCatchup)
	}
	req.TxIds.Mode = None
	p2p.Send(p.rw, MsgTypeRequestMsg, req)
}
func (sm *SyncManager) resetLibNum(p *peer) {
	if sm.state == SyncStageInSync && sm.source != nil {
		sm.source.Reset()
	}
	if !p.syncing {
		if p.lastHandshakeRecv.LastIrreversibleBlockNum > sm.syncKnownlibNum {
			sm.syncKnownlibNum = p.lastHandshakeRecv.LastIrreversibleBlockNum
		}
	} else if sm.source != nil && reflect.DeepEqual(*p, *(sm.source)) {
		sm.syncLastRequestedNum = 0
		sm.requestNextChunk(p)
	}
}
func (sm *SyncManager) recvNotice(p *peer, msg *NoticeMsg) {
	p.Log().Trace("sync_manager got %s block notice", msg.BlkIds.Mode)
	if msg.BlkIds.Mode == CatchUp {
		if len(msg.BlkIds.Ids) == 0 {
			p.Log().Trace("got a catch up with ids size = 0")
		} else {
			if val, ok := msg.BlkIds.Ids[len(msg.BlkIds.Ids)-1].(common.Hash); ok {
				sm.verifyCatchup(p, msg.BlkIds.Pending, val)
			}
		}
	} else {
		p.lastHandshakeRecv.LastIrreversibleBlockNum = msg.TxIds.Pending
		sm.resetLibNum(p)
		sm.startSync(p, msg.BlkIds.Pending)
	}
}
func (sm *SyncManager) recvBlock(p *peer, blkID common.Hash, blkNum uint32) {
	if sm.state == SyncStageLibCatchup {
		if blkNum != sm.syncNextExpectedNum {
			p.Log().Trace("expected block %d but got %d", sm.syncNextExpectedNum, blkNum)
			P2pPlugin_impl.pm.peers.Unregister(p.id)
			return
		}
		sm.syncNextExpectedNum = blkNum + 1
	}
	if sm.state == SyncStageHeadCatchup {
		p.Log().Trace("sync_manager in head_catchup state")
		sm.setState(SyncStageInSync)
		p.Reset()

		nilID := common.HexToHash("")
		for _, cp := range P2pPlugin_impl.pm.peers.peers {
			if reflect.DeepEqual(cp.forkHead, nilID) {
				continue
			}
			if reflect.DeepEqual(cp.forkHead, blkID) || cp.forkHeadNum < blkNum {
				p.forkHead = nilID
				p.forkHeadNum = 0
			} else {
				sm.setState(SyncStageHeadCatchup)
			}
		}
	} else if sm.state == SyncStageLibCatchup {
		if blkNum == sm.syncKnownlibNum {
			p.Log().Trace("All caught up with last known last irreversible block resending handshake")
			sm.setState(SyncStageInSync)
			sm.sendHandshakes()
		}
	} else if blkNum == sm.syncLastRequestedNum {
		sm.requestNextChunk(p)
	}
}

func (sm *SyncManager) rejectedBlock(p *peer, blkNum uint32) {
	if sm.state != SyncStageInSync {
		p.Log().Trace("block %d not accepted from %s", blkNum, p.id)
		sm.syncLastRequestedNum = 0
		sm.source.Reset()
		P2pPlugin_impl.pm.peers.Unregister(p.id)
		sm.setState(SyncStageInSync)
		sm.sendHandshakes()
	}
}

func (sm *SyncManager) sendHandshakes() {
	for _, peer := range P2pPlugin_impl.pm.peers.peers {
		if !peer.syncing {
			peer.sendHandshake()
		}
	}
}

func (sm *SyncManager) isActive(p *peer) bool {
	if sm.state == SyncStageHeadCatchup {
		//fhset := reflect.DeepEqual(p.forkHead, new(common.Hash))

		return !reflect.DeepEqual(p.forkHead, common.HexToHash("")) && p.forkHeadNum < sm.chain.Head.BlockNum
	}
	return sm.state != SyncStageInSync
}

/*-------DispatchManager--------------*/
func (ds *DispatchManager) recvNotice(p *peer, msg *NoticeMsg, generated bool) {
	req := RequestMsg{}
	req.TxIds.Mode = None
	req.BlkIds.Mode = None

	sendReq := bool(false)
	if msg.TxIds.Mode == Normal {
		req.TxIds.Mode = Normal
		req.TxIds.Pending = 0
		for _, trxID := range msg.TxIds.Ids {
			if id, ok := trxID.(common.Hash); ok {
				if _, exist := P2pPlugin_impl.pm.localTxns[id]; !exist { //chain.UnAppliedTransaction[id]
					//did not find trx
					//At this point the details of the txn are not known, just its id. This
					//effectively gives 120 seconds to learn of the details of the txn which
					//will update the expiry in bcast_transaction
					p.trxState[id] = transactionState{
						Id:              id,
						IsKnownByPeer:   true,
						IsNoticedToPeer: true,
						Expries:         uint32(time.Now().Unix() + 120),
						RequestedTime:   time.Now(),
					}
					req.TxIds.Ids = append(req.TxIds.Ids, id)
					ds.reqTrx[id] = struct{}{}
				} else {
					log.Info("big msg manager found txn id in table,%s", id)
				}
			}
		}
		sendReq = len(req.TxIds.Ids) > 0
		log.Info("big msg manager send_req ids list has %d entries", len(req.TxIds.Ids))
	} else if msg.TxIds.Mode != None {
		log.Info("passed a notice_message with something other than a normal on none known_trx")
		return
	}

	if msg.BlkIds.Mode == Normal {
		req.BlkIds.Mode = Normal
		for _, blkID := range msg.BlkIds.Ids {
			if id, ok := blkID.(common.Hash); ok {
				entry := peerBlockState{
					Id:            id,
					BlkNum:        0,
					IsKnown:       true,
					IsNoticed:     true,
					RequestedTime: time.Now(),
				}

				if data := P2pPlugin_impl.pm.syncMaster.chain.ForkDB.GetBlock(id); data != nil {
					entry.BlkNum = data.BlockNum
				} else {
					sendReq = true
					req.BlkIds.Ids = append(req.BlkIds.Ids, id)
					ds.reqBlks[id] = blockRequest{
						id:         id,
						localRetry: generated,
					}
					entry.RequestedTime = time.Now()
				}
				p.addPeerBlock(&entry)
			}
		}
	} else if msg.BlkIds.Mode != None {
		log.Info("passed a notice_message with something other than a normal on none known_blocks")
		return
	}
	if sendReq {
		p2p.Send(p.rw, MsgTypeRequestMsg, req)
		p.lastReq = req
	}
}
func (ds *DispatchManager) recvBlock(p *peer, blkID common.Hash, blkNum uint32) {
	ds.receivedBlks[blkID] = blockOrigin{
		id:     blkID,
		origin: p,
	}
	if p.lastReq.BlkIds.Mode != None && reflect.DeepEqual(p.lastReq.BlkIds.Ids[len(p.lastReq.BlkIds.Ids)-1], blkID) {
		p.lastReq = RequestMsg{}
	}
	entry := peerBlockState{
		Id:            blkID,
		BlkNum:        blkNum,
		IsKnown:       false,
		IsNoticed:     true,
		RequestedTime: time.Now(),
	}
	p.addPeerBlock(&entry)
}

func (ds *DispatchManager) recvTransaction(p *peer, id common.Hash) {
	ds.receivedTrx[id] = transactionOrigin{
		id:     id,
		origin: p,
	}
	if p.lastReq.TxIds.Mode != None && reflect.DeepEqual(p.lastReq.TxIds.Ids[len(p.lastReq.TxIds.Ids)-1], id) {
		p.lastReq = RequestMsg{}
	}
}
func (ds *DispatchManager) bcastTransaction(trx *types.PackedTransaction) {
	id := trx.ID()
	var skip *peer
	for _, org := range ds.receivedTrx {
		if reflect.DeepEqual(org.id, id) {
			skip = org.origin
			delete(ds.receivedTrx, id)
			break
		}
	}

	for ref := range ds.reqTrx {
		if reflect.DeepEqual(ref, id) {
			delete(ds.reqTrx, ref)
		}
	}

	if _, exist := P2pPlugin_impl.pm.localTxns[id]; exist {
		//found trxid in local_trxs
		return
	}
	trxExpiration := trx.Expiration()

	nts := nodeTranscationState{
		id:        id,
		expires:   uint32(trxExpiration),
		packedTxn: *trx,
		blkNum:    0,
		trueBlk:   0,
		requests:  0,
	}
	P2pPlugin_impl.pm.localTxns[id] = nts
	P2pPlugin_impl.pm.sendAll(MsgTypePackedTransaction, *trx, func(p *peer) bool {
		if p.syncing || reflect.DeepEqual(p, skip) {
			return false
		}
		bs, exist := p.trxState[id]
		if !exist {
			p.trxState[id] = transactionState{
				Id:              id,
				IsKnownByPeer:   false,
				IsNoticedToPeer: true,
				BlkNum:          0,
				Expries:         uint32(trxExpiration),
				RequestedTime:   time.Now(),
			}
		} else {
			bs.Expries = uint32(trxExpiration)
			p.trxState[id] = bs
		}
		return !exist
	})
}

func (ds *DispatchManager) bcastBlock(bsum *types.BlockState) {
	var skip *peer
	for _, org := range ds.receivedBlks {
		if reflect.DeepEqual(org.id, bsum.ID) {
			skip = org.origin
			delete(ds.receivedBlks, org.id)
			break
		}
	}

	pbstate := peerBlockState{
		Id:            bsum.ID,
		BlkNum:        bsum.BlockNum,
		IsKnown:       true,
		IsNoticed:     true,
		RequestedTime: time.Now(),
	}

	for _, p := range P2pPlugin_impl.peers.peers {
		if reflect.DeepEqual(p, skip) || p.syncing {
			continue
		}
		p.addPeerBlock(&pbstate)
		p2p.Send(p.rw, MsgTypeSignedBlock, bsum)
	}
}

func (ds *DispatchManager) rejectedTransaction(id common.Hash) {
	for _, org := range ds.receivedTrx {
		if org.id == id {
			delete(ds.receivedTrx, id)
			break
		}
	}
}

/*----------------peer------------------*/
func (p *peer) blkSendBranch() {
	headNum := P2pPlugin_impl.pm.syncMaster.chain.Head.BlockNum
	note := NoticeMsg{}

	note.BlkIds.Mode = Normal
	note.BlkIds.Pending = 0
	if headNum == 0 {
		p2p.Send(p.rw, MsgTypeNoticeMsg, note)
		return
	}
	var headId common.Hash
	var LibId common.Hash
	var LibNum uint32

	LibNum = P2pPlugin_impl.pm.syncMaster.chain.Head.DposIrreversibleBlockNum
	LibId = P2pPlugin_impl.pm.syncMaster.chain.ForkDB.GetBlockInChain(LibNum).ID
	headId = P2pPlugin_impl.pm.syncMaster.chain.Head.ID
	bstack := make([](*types.BlockState), 0)
	for bid := headId; !common.EmptyHash(bid) && bid != LibId; {
		blk := P2pPlugin_impl.pm.syncMaster.chain.ForkDB.GetBlock(bid)
		if blk != nil {
			bid = blk.Header.Previous
			bstack = append(bstack, blk)
		}
	}
	count := int(0)
	if len(bstack) > 0 {
		if reflect.DeepEqual(bstack[len(bstack)-1].Header.Previous, LibId) {
			count = len(bstack)
			for i := len(bstack) - 1; i >= 0; i-- {
				p2p.Send(p.rw, MsgTypeSignedBlock, *(bstack[i]))
			}
		}
		p.Log().Trace("Sent %d blocks on my fork", count)
	} else {
		p.Log().Trace("Nothing to send on fork request")
	}
	p.syncing = false
}

func (p *peer) blkSend(ids []interface{}) {
	chain := P2pPlugin_impl.pm.syncMaster.chain
	count := int(0)
	for _, id := range ids {
		count = count + 1
		if val, ok := id.(common.Hash); ok {
			blkPtr := chain.ForkDB.GetBlock(val)
			if blkPtr != nil {
				p2p.Send(p.rw, MsgTypeSignedBlock, *blkPtr)
			}
		}
	}
}
func (p *peer) txnSendPending(ids []interface{}) {
	for _, tx := range P2pPlugin_impl.pm.localTxns {
		if tx.blkNum == 0 {
			found := bool(false)
			for _, know := range ids {
				if id, ok := know.(common.Hash); ok {
					if reflect.DeepEqual(id, tx.id) {
						found = true
						break
					}
				}
			}
			if !found {
				tx.updateInFlight(1)
				P2pPlugin_impl.pm.localTxns[tx.id] = tx
				p2p.Send(p.rw, MsgTypePackedTransaction, tx.packedTxn)
				if val, exist := P2pPlugin_impl.pm.localTxns[tx.id]; exist {
					val.updateInFlight(-1)
					P2pPlugin_impl.pm.localTxns[tx.id] = val
				}
			}
		}
	}
}
func (p *peer) txnSend(ids []interface{}) {
	for _, val := range ids {
		if id, ok := val.(common.Hash); ok {
			if tx, exist := P2pPlugin_impl.pm.localTxns[id]; exist {
				tx.updateInFlight(1)
				P2pPlugin_impl.pm.localTxns[id] = tx
				p2p.Send(p.rw, MsgTypePackedTransaction, tx.packedTxn)
				if val, exist := P2pPlugin_impl.pm.localTxns[tx.id]; exist {
					val.updateInFlight(-1)
					P2pPlugin_impl.pm.localTxns[tx.id] = val
				}
			}
		}
	}
}
func (p *peer) stopSend() {
	p.syncing = false
}
func (p *peer) enqueueSyncBlock() bool {
	chain := P2pPlugin_impl.pm.syncMaster.chain
	if (reflect.DeepEqual(p.peerRequested, syncState{})) {
		return false
	}
	p.peerRequested.last++
	num := p.peerRequested.last
	if num == p.peerRequested.endBlock {
		p.peerRequested = syncState{}
	}
	blkPtr := chain.ForkDB.GetBlockInChain(num)
	if blkPtr != nil {
		p2p.Send(p.rw, MsgTypeSignedBlock, *blkPtr)
		return true
	}
	return false
}
func (p *peer) sendHandshake() {
	p.populate(&p.lastHandshakeSent)
	p.sentHandshakeCount++
	p.lastHandshakeSent.Generation++
	p.Log().Trace("Sending handshake generation %d to %s", p.lastHandshakeSent.Generation, p.id)
	p2p.Send(p.rw, MsgTypeStatusMsg, p.lastHandshakeSent)
}

func (p *peer) populate(hello *HandShakeMsg) {
	hello.NetworkId = uint32(P2pPlugin_impl.pm.networkId)
	hello.TimeStamp = time.Now()
	//token key sig ,wait to do
	hello.PeerID = P2pPlugin_impl.Status().ID
	cc := P2pPlugin_impl.pm.syncMaster.chain
	hello.HeadNum = cc.Head.BlockNum
	hello.LastIrreversibleBlockNum = cc.Head.DposIrreversibleBlockNum

	if hello.LastIrreversibleBlockNum != 0 {
		if data := cc.ForkDB.GetBlockInChain(hello.LastIrreversibleBlockNum); data != nil {
			hello.LastIrreversibleBlockId = data.ID
		}
	}

	//HeadNum maybe 0
	if 1 == 1 || hello.HeadNum != 0 {
		if data := cc.ForkDB.GetBlockInChain(hello.HeadNum); data != nil {
			hello.HeadId = data.ID
		}
	}
}
func (p *peer) addPeerBlock(entry *peerBlockState) bool {
	pbs, exist := p.blkState[entry.Id]
	if !exist {
		p.blkState[entry.Id] = *entry
	} else {
		pbs.IsKnown = true
		if entry.BlkNum == 0 {
			pbs.BlkNum = entry.BlkNum
		} else {
			pbs.RequestedTime = time.Now()
		}
	}
	return !exist
}
