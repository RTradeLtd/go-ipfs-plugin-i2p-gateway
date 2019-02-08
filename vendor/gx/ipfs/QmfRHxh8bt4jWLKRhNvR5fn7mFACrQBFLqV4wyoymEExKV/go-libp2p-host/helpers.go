package host

import pstore "gx/ipfs/QmQFFp4ntkd4C14sP3FaH9WJyBuetuGUVo6dShNHvnoEvC/go-libp2p-peerstore"

// PeerInfoFromHost returns a PeerInfo struct with the Host's ID and all of its Addrs.
func PeerInfoFromHost(h Host) *pstore.PeerInfo {
	return &pstore.PeerInfo{
		ID:    h.ID(),
		Addrs: h.Addrs(),
	}
}
