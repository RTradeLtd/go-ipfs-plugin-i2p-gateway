package libp2pquic

import (
	ma "gx/ipfs/QmNTCey11oxhb1AxDnQBRHtdhap6Ctud872NjAYPYYXPuc/go-multiaddr"
	ic "gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	quic "gx/ipfs/QmU44KWVkSHno7sNDTeUcL4FBgxgoidkFuTUyTXWJPXXFJ/quic-go"
	tpt "gx/ipfs/QmUDtgnEr7FFrtK2LQM2dFzTNWghnrApBDcU3iHEJz8eQS/go-libp2p-transport"
	smux "gx/ipfs/QmY9JXR3FupnYAYJWK9aMr9bCpqWKcToQ1tz8DVGTrHpHw/go-stream-muxer"
)

type conn struct {
	sess      quic.Session
	transport tpt.Transport

	localPeer      peer.ID
	privKey        ic.PrivKey
	localMultiaddr ma.Multiaddr

	remotePeerID    peer.ID
	remotePubKey    ic.PubKey
	remoteMultiaddr ma.Multiaddr
}

var _ tpt.Conn = &conn{}

func (c *conn) Close() error {
	return c.sess.Close()
}

// IsClosed returns whether a connection is fully closed.
func (c *conn) IsClosed() bool {
	return c.sess.Context().Err() != nil
}

// OpenStream creates a new stream.
func (c *conn) OpenStream() (smux.Stream, error) {
	qstr, err := c.sess.OpenStreamSync()
	return &stream{Stream: qstr}, err
}

// AcceptStream accepts a stream opened by the other side.
func (c *conn) AcceptStream() (smux.Stream, error) {
	qstr, err := c.sess.AcceptStream()
	return &stream{Stream: qstr}, err
}

// LocalPeer returns our peer ID
func (c *conn) LocalPeer() peer.ID {
	return c.localPeer
}

// LocalPrivateKey returns our private key
func (c *conn) LocalPrivateKey() ic.PrivKey {
	return c.privKey
}

// RemotePeer returns the peer ID of the remote peer.
func (c *conn) RemotePeer() peer.ID {
	return c.remotePeerID
}

// RemotePublicKey returns the public key of the remote peer.
func (c *conn) RemotePublicKey() ic.PubKey {
	return c.remotePubKey
}

// LocalMultiaddr returns the local Multiaddr associated
func (c *conn) LocalMultiaddr() ma.Multiaddr {
	return c.localMultiaddr
}

// RemoteMultiaddr returns the remote Multiaddr associated
func (c *conn) RemoteMultiaddr() ma.Multiaddr {
	return c.remoteMultiaddr
}

func (c *conn) Transport() tpt.Transport {
	return c.transport
}
