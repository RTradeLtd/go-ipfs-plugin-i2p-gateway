package sessionmanager

import (
	"context"
	"testing"
	"time"

	bssrs "gx/ipfs/QmYJ48z7NEzo3u2yCvUvNtBQ7wJWd5dX2nxxc7FeA6nHq1/go-bitswap/sessionrequestsplitter"

	bssession "gx/ipfs/QmYJ48z7NEzo3u2yCvUvNtBQ7wJWd5dX2nxxc7FeA6nHq1/go-bitswap/session"

	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	cid "gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	blocks "gx/ipfs/QmWoXtvgC8inqFkAATB7cp2Dax7XBi9VDvSg9RCCZufmRk/go-block-format"
)

type fakeSession struct {
	interested            bool
	receivedBlock         bool
	updateReceiveCounters bool
	id                    uint64
	pm                    *fakePeerManager
	srs                   *fakeRequestSplitter
}

func (*fakeSession) GetBlock(context.Context, cid.Cid) (blocks.Block, error) {
	return nil, nil
}
func (*fakeSession) GetBlocks(context.Context, []cid.Cid) (<-chan blocks.Block, error) {
	return nil, nil
}
func (fs *fakeSession) InterestedIn(cid.Cid) bool              { return fs.interested }
func (fs *fakeSession) ReceiveBlockFrom(peer.ID, blocks.Block) { fs.receivedBlock = true }
func (fs *fakeSession) UpdateReceiveCounters(blocks.Block)     { fs.updateReceiveCounters = true }

type fakePeerManager struct {
	id uint64
}

func (*fakePeerManager) FindMorePeers(context.Context, cid.Cid)  {}
func (*fakePeerManager) GetOptimizedPeers() []peer.ID            { return nil }
func (*fakePeerManager) RecordPeerRequests([]peer.ID, []cid.Cid) {}
func (*fakePeerManager) RecordPeerResponse(peer.ID, cid.Cid)     {}

type fakeRequestSplitter struct {
}

func (frs *fakeRequestSplitter) SplitRequest(peers []peer.ID, keys []cid.Cid) []*bssrs.PartialRequest {
	return nil
}
func (frs *fakeRequestSplitter) RecordDuplicateBlock() {}
func (frs *fakeRequestSplitter) RecordUniqueBlock()    {}

var nextInterestedIn bool

func sessionFactory(ctx context.Context, id uint64, pm bssession.PeerManager, srs bssession.RequestSplitter) Session {
	return &fakeSession{
		interested:    nextInterestedIn,
		receivedBlock: false,
		id:            id,
		pm:            pm.(*fakePeerManager),
		srs:           srs.(*fakeRequestSplitter),
	}
}

func peerManagerFactory(ctx context.Context, id uint64) bssession.PeerManager {
	return &fakePeerManager{id}
}

func requestSplitterFactory(ctx context.Context) bssession.RequestSplitter {
	return &fakeRequestSplitter{}
}

func TestAddingSessions(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	sm := New(ctx, sessionFactory, peerManagerFactory, requestSplitterFactory)

	p := peer.ID(123)
	block := blocks.NewBlock([]byte("block"))
	// we'll be interested in all blocks for this test
	nextInterestedIn = true

	currentID := sm.GetNextSessionID()
	firstSession := sm.NewSession(ctx).(*fakeSession)
	if firstSession.id != firstSession.pm.id ||
		firstSession.id != currentID+1 {
		t.Fatal("session does not have correct id set")
	}
	secondSession := sm.NewSession(ctx).(*fakeSession)
	if secondSession.id != secondSession.pm.id ||
		secondSession.id != firstSession.id+1 {
		t.Fatal("session does not have correct id set")
	}
	sm.GetNextSessionID()
	thirdSession := sm.NewSession(ctx).(*fakeSession)
	if thirdSession.id != thirdSession.pm.id ||
		thirdSession.id != secondSession.id+2 {
		t.Fatal("session does not have correct id set")
	}
	sm.ReceiveBlockFrom(p, block)
	if !firstSession.receivedBlock ||
		!secondSession.receivedBlock ||
		!thirdSession.receivedBlock {
		t.Fatal("should have received blocks but didn't")
	}
}

func TestReceivingBlocksWhenNotInterested(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	sm := New(ctx, sessionFactory, peerManagerFactory, requestSplitterFactory)

	p := peer.ID(123)
	block := blocks.NewBlock([]byte("block"))
	// we'll be interested in all blocks for this test
	nextInterestedIn = false
	firstSession := sm.NewSession(ctx).(*fakeSession)
	nextInterestedIn = true
	secondSession := sm.NewSession(ctx).(*fakeSession)
	nextInterestedIn = false
	thirdSession := sm.NewSession(ctx).(*fakeSession)

	sm.ReceiveBlockFrom(p, block)
	if firstSession.receivedBlock ||
		!secondSession.receivedBlock ||
		thirdSession.receivedBlock {
		t.Fatal("did not receive blocks only for interested sessions")
	}
}

func TestRemovingPeersWhenManagerContextCancelled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sm := New(ctx, sessionFactory, peerManagerFactory, requestSplitterFactory)

	p := peer.ID(123)
	block := blocks.NewBlock([]byte("block"))
	// we'll be interested in all blocks for this test
	nextInterestedIn = true
	firstSession := sm.NewSession(ctx).(*fakeSession)
	secondSession := sm.NewSession(ctx).(*fakeSession)
	thirdSession := sm.NewSession(ctx).(*fakeSession)

	cancel()
	// wait for sessions to get removed
	time.Sleep(10 * time.Millisecond)
	sm.ReceiveBlockFrom(p, block)
	if firstSession.receivedBlock ||
		secondSession.receivedBlock ||
		thirdSession.receivedBlock {
		t.Fatal("received blocks for sessions after manager is shutdown")
	}
}

func TestRemovingPeersWhenSessionContextCancelled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	sm := New(ctx, sessionFactory, peerManagerFactory, requestSplitterFactory)

	p := peer.ID(123)
	block := blocks.NewBlock([]byte("block"))
	// we'll be interested in all blocks for this test
	nextInterestedIn = true
	firstSession := sm.NewSession(ctx).(*fakeSession)
	sessionCtx, sessionCancel := context.WithCancel(ctx)
	secondSession := sm.NewSession(sessionCtx).(*fakeSession)
	thirdSession := sm.NewSession(ctx).(*fakeSession)

	sessionCancel()
	// wait for sessions to get removed
	time.Sleep(10 * time.Millisecond)
	sm.ReceiveBlockFrom(p, block)
	if !firstSession.receivedBlock ||
		secondSession.receivedBlock ||
		!thirdSession.receivedBlock {
		t.Fatal("received blocks for sessions that are canceled")
	}
}
