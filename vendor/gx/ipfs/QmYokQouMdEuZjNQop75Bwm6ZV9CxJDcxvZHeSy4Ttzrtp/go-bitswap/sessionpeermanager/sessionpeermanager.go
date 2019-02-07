package sessionpeermanager

import (
	"context"
	"fmt"
	"math/rand"

	logging "gx/ipfs/QmcuXC5cxs79ro2cUuHs4HQ2bkDLJUYokwL8aivcX6HW3C/go-log"

	cid "gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	ifconnmgr "gx/ipfs/QmSFo2QrMF4M1mKdB291ZqNtsie4NfwXCRdWgDU3inw4Ff/go-libp2p-interface-connmgr"
	peer "gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
)

var log = logging.Logger("bitswap")

const (
	maxOptimizedPeers = 32
	reservePeers      = 2
)

// PeerNetwork is an interface for finding providers and managing connections
type PeerNetwork interface {
	ConnectionManager() ifconnmgr.ConnManager
	ConnectTo(context.Context, peer.ID) error
	FindProvidersAsync(context.Context, cid.Cid, int) <-chan peer.ID
}

type peerMessage interface {
	handle(spm *SessionPeerManager)
}

// SessionPeerManager tracks and manages peers for a session, and provides
// the best ones to the session
type SessionPeerManager struct {
	ctx     context.Context
	network PeerNetwork
	tag     string

	peerMessages chan peerMessage

	// do not touch outside of run loop
	activePeers         map[peer.ID]bool
	unoptimizedPeersArr []peer.ID
	optimizedPeersArr   []peer.ID
}

// New creates a new SessionPeerManager
func New(ctx context.Context, id uint64, network PeerNetwork) *SessionPeerManager {
	spm := &SessionPeerManager{
		ctx:          ctx,
		network:      network,
		peerMessages: make(chan peerMessage, 16),
		activePeers:  make(map[peer.ID]bool),
	}

	spm.tag = fmt.Sprint("bs-ses-", id)

	go spm.run(ctx)
	return spm
}

// RecordPeerResponse records that a peer received a block, and adds to it
// the list of peers if it wasn't already added
func (spm *SessionPeerManager) RecordPeerResponse(p peer.ID, k cid.Cid) {

	// at the moment, we're just adding peers here
	// in the future, we'll actually use this to record metrics
	select {
	case spm.peerMessages <- &peerResponseMessage{p}:
	case <-spm.ctx.Done():
	}
}

// RecordPeerRequests records that a given set of peers requested the given cids
func (spm *SessionPeerManager) RecordPeerRequests(p []peer.ID, ks []cid.Cid) {
	// at the moment, we're not doing anything here
	// soon we'll use this to track latency by peer
}

// GetOptimizedPeers returns the best peers available for a session
func (spm *SessionPeerManager) GetOptimizedPeers() []peer.ID {
	// right now this just returns all peers, but soon we might return peers
	// ordered by optimization, or only a subset
	resp := make(chan []peer.ID, 1)
	select {
	case spm.peerMessages <- &peerReqMessage{resp}:
	case <-spm.ctx.Done():
		return nil
	}

	select {
	case peers := <-resp:
		return peers
	case <-spm.ctx.Done():
		return nil
	}
}

// FindMorePeers attempts to find more peers for a session by searching for
// providers for the given Cid
func (spm *SessionPeerManager) FindMorePeers(ctx context.Context, c cid.Cid) {
	go func(k cid.Cid) {
		// TODO: have a task queue setup for this to:
		// - rate limit
		// - manage timeouts
		// - ensure two 'findprovs' calls for the same block don't run concurrently
		// - share peers between sessions based on interest set
		for p := range spm.network.FindProvidersAsync(ctx, k, 10) {
			go func(p peer.ID) {
				// TODO: Also use context from spm.
				err := spm.network.ConnectTo(ctx, p)
				if err != nil {
					log.Debugf("failed to connect to provider %s: %s", p, err)
				}
				select {
				case spm.peerMessages <- &peerFoundMessage{p}:
				case <-ctx.Done():
				case <-spm.ctx.Done():
				}
			}(p)
		}
	}(c)
}

func (spm *SessionPeerManager) run(ctx context.Context) {
	for {
		select {
		case pm := <-spm.peerMessages:
			pm.handle(spm)
		case <-ctx.Done():
			spm.handleShutdown()
			return
		}
	}
}

func (spm *SessionPeerManager) tagPeer(p peer.ID) {
	cmgr := spm.network.ConnectionManager()
	cmgr.TagPeer(p, spm.tag, 10)
}

func (spm *SessionPeerManager) insertOptimizedPeer(p peer.ID) {
	if len(spm.optimizedPeersArr) >= (maxOptimizedPeers - reservePeers) {
		tailPeer := spm.optimizedPeersArr[len(spm.optimizedPeersArr)-1]
		spm.optimizedPeersArr = spm.optimizedPeersArr[:len(spm.optimizedPeersArr)-1]
		spm.unoptimizedPeersArr = append(spm.unoptimizedPeersArr, tailPeer)
	}

	spm.optimizedPeersArr = append([]peer.ID{p}, spm.optimizedPeersArr...)
}

func (spm *SessionPeerManager) removeOptimizedPeer(p peer.ID) {
	for i := 0; i < len(spm.optimizedPeersArr); i++ {
		if spm.optimizedPeersArr[i] == p {
			spm.optimizedPeersArr = append(spm.optimizedPeersArr[:i], spm.optimizedPeersArr[i+1:]...)
			return
		}
	}
}

func (spm *SessionPeerManager) removeUnoptimizedPeer(p peer.ID) {
	for i := 0; i < len(spm.unoptimizedPeersArr); i++ {
		if spm.unoptimizedPeersArr[i] == p {
			spm.unoptimizedPeersArr[i] = spm.unoptimizedPeersArr[len(spm.unoptimizedPeersArr)-1]
			spm.unoptimizedPeersArr = spm.unoptimizedPeersArr[:len(spm.unoptimizedPeersArr)-1]
			return
		}
	}
}

type peerFoundMessage struct {
	p peer.ID
}

func (pfm *peerFoundMessage) handle(spm *SessionPeerManager) {
	p := pfm.p
	if _, ok := spm.activePeers[p]; !ok {
		spm.activePeers[p] = false
		spm.unoptimizedPeersArr = append(spm.unoptimizedPeersArr, p)
		spm.tagPeer(p)
	}
}

type peerResponseMessage struct {
	p peer.ID
}

func (prm *peerResponseMessage) handle(spm *SessionPeerManager) {

	p := prm.p
	isOptimized, ok := spm.activePeers[p]
	if !ok {
		spm.activePeers[p] = true
		spm.tagPeer(p)
	} else {
		if isOptimized {
			spm.removeOptimizedPeer(p)
		} else {
			spm.activePeers[p] = true
			spm.removeUnoptimizedPeer(p)
		}
	}
	spm.insertOptimizedPeer(p)
}

type peerReqMessage struct {
	resp chan<- []peer.ID
}

func (prm *peerReqMessage) handle(spm *SessionPeerManager) {
	randomOrder := rand.Perm(len(spm.unoptimizedPeersArr))
	maxPeers := len(spm.unoptimizedPeersArr) + len(spm.optimizedPeersArr)
	if maxPeers > maxOptimizedPeers {
		maxPeers = maxOptimizedPeers
	}

	extraPeers := make([]peer.ID, maxPeers-len(spm.optimizedPeersArr))
	for i := range extraPeers {
		extraPeers[i] = spm.unoptimizedPeersArr[randomOrder[i]]
	}
	prm.resp <- append(spm.optimizedPeersArr, extraPeers...)
}

func (spm *SessionPeerManager) handleShutdown() {
	cmgr := spm.network.ConnectionManager()
	for p := range spm.activePeers {
		cmgr.UntagPeer(p, spm.tag)
	}
}
