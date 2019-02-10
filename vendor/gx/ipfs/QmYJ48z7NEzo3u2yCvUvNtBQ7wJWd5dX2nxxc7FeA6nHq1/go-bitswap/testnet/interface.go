package bitswap

import (
	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	"gx/ipfs/QmVnJMgafh5MBYiyqbvDtoCL8pcQvbEGD2k9o9GFpBWPzY/go-testutil"
	bsnet "gx/ipfs/QmYJ48z7NEzo3u2yCvUvNtBQ7wJWd5dX2nxxc7FeA6nHq1/go-bitswap/network"
)

type Network interface {
	Adapter(testutil.Identity) bsnet.BitSwapNetwork

	HasPeer(peer.ID) bool
}
