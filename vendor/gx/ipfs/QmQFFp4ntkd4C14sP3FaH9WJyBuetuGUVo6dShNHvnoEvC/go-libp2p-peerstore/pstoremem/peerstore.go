package pstoremem

import pstore "gx/ipfs/QmQFFp4ntkd4C14sP3FaH9WJyBuetuGUVo6dShNHvnoEvC/go-libp2p-peerstore"

// NewPeerstore creates an in-memory threadsafe collection of peers.
func NewPeerstore() pstore.Peerstore {
	return pstore.NewPeerstore(
		NewKeyBook(),
		NewAddrBook(),
		NewPeerMetadata())
}
