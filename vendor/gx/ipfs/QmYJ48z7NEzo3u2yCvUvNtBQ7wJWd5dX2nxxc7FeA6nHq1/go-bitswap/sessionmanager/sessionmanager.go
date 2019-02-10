package sessionmanager

import (
	"context"
	"sync"

	cid "gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	blocks "gx/ipfs/QmWoXtvgC8inqFkAATB7cp2Dax7XBi9VDvSg9RCCZufmRk/go-block-format"

	exchange "gx/ipfs/QmP2g3VxmC7g7fyRJDj1VJ72KHZbJ9UW24YjSWEj1XTb4H/go-ipfs-exchange-interface"
	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	bssession "gx/ipfs/QmYJ48z7NEzo3u2yCvUvNtBQ7wJWd5dX2nxxc7FeA6nHq1/go-bitswap/session"
)

// Session is a session that is managed by the session manager
type Session interface {
	exchange.Fetcher
	InterestedIn(cid.Cid) bool
	ReceiveBlockFrom(peer.ID, blocks.Block)
	UpdateReceiveCounters(blocks.Block)
}

type sesTrk struct {
	session Session
	pm      bssession.PeerManager
	srs     bssession.RequestSplitter
}

// SessionFactory generates a new session for the SessionManager to track.
type SessionFactory func(ctx context.Context, id uint64, pm bssession.PeerManager, srs bssession.RequestSplitter) Session

// RequestSplitterFactory generates a new request splitter for a session.
type RequestSplitterFactory func(ctx context.Context) bssession.RequestSplitter

// PeerManagerFactory generates a new peer manager for a session.
type PeerManagerFactory func(ctx context.Context, id uint64) bssession.PeerManager

// SessionManager is responsible for creating, managing, and dispatching to
// sessions.
type SessionManager struct {
	ctx                    context.Context
	sessionFactory         SessionFactory
	peerManagerFactory     PeerManagerFactory
	requestSplitterFactory RequestSplitterFactory

	// Sessions
	sessLk   sync.Mutex
	sessions []sesTrk

	// Session Index
	sessIDLk sync.Mutex
	sessID   uint64
}

// New creates a new SessionManager.
func New(ctx context.Context, sessionFactory SessionFactory, peerManagerFactory PeerManagerFactory, requestSplitterFactory RequestSplitterFactory) *SessionManager {
	return &SessionManager{
		ctx:                    ctx,
		sessionFactory:         sessionFactory,
		peerManagerFactory:     peerManagerFactory,
		requestSplitterFactory: requestSplitterFactory,
	}
}

// NewSession initializes a session with the given context, and adds to the
// session manager.
func (sm *SessionManager) NewSession(ctx context.Context) exchange.Fetcher {
	id := sm.GetNextSessionID()
	sessionctx, cancel := context.WithCancel(ctx)

	pm := sm.peerManagerFactory(sessionctx, id)
	srs := sm.requestSplitterFactory(sessionctx)
	session := sm.sessionFactory(sessionctx, id, pm, srs)
	tracked := sesTrk{session, pm, srs}
	sm.sessLk.Lock()
	sm.sessions = append(sm.sessions, tracked)
	sm.sessLk.Unlock()
	go func() {
		defer cancel()
		select {
		case <-sm.ctx.Done():
			sm.removeSession(tracked)
		case <-ctx.Done():
			sm.removeSession(tracked)
		}
	}()

	return session
}

func (sm *SessionManager) removeSession(session sesTrk) {
	sm.sessLk.Lock()
	defer sm.sessLk.Unlock()
	for i := 0; i < len(sm.sessions); i++ {
		if sm.sessions[i] == session {
			sm.sessions[i] = sm.sessions[len(sm.sessions)-1]
			sm.sessions = sm.sessions[:len(sm.sessions)-1]
			return
		}
	}
}

// GetNextSessionID returns the next sequentional identifier for a session.
func (sm *SessionManager) GetNextSessionID() uint64 {
	sm.sessIDLk.Lock()
	defer sm.sessIDLk.Unlock()
	sm.sessID++
	return sm.sessID
}

// ReceiveBlockFrom receives a block from a peer and dispatches to interested
// sessions.
func (sm *SessionManager) ReceiveBlockFrom(from peer.ID, blk blocks.Block) {
	sm.sessLk.Lock()
	defer sm.sessLk.Unlock()

	k := blk.Cid()
	for _, s := range sm.sessions {
		if s.session.InterestedIn(k) {
			s.session.ReceiveBlockFrom(from, blk)
		}
	}
}

// UpdateReceiveCounters records the fact that a block was received, allowing
// sessions to track duplicates
func (sm *SessionManager) UpdateReceiveCounters(blk blocks.Block) {
	sm.sessLk.Lock()
	defer sm.sessLk.Unlock()

	for _, s := range sm.sessions {
		s.session.UpdateReceiveCounters(blk)
	}
}
