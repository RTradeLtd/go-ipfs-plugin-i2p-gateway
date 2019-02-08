package decision

import (
	"fmt"
	"math"
	"testing"

	u "gx/ipfs/QmNohiVssaPw3KVLZik59DBVGTSm2dGvYT9eoXt5DQ36Yz/go-ipfs-util"
	"gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	cid "gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	"gx/ipfs/QmVnJMgafh5MBYiyqbvDtoCL8pcQvbEGD2k9o9GFpBWPzY/go-testutil"
	"gx/ipfs/QmYJ48z7NEzo3u2yCvUvNtBQ7wJWd5dX2nxxc7FeA6nHq1/go-bitswap/wantlist"
)

// FWIW: At the time of this commit, including a timestamp in task increases
// time cost of Push by 3%.
func BenchmarkTaskQueuePush(b *testing.B) {
	q := newPRQ()
	peers := []peer.ID{
		testutil.RandPeerIDFatal(b),
		testutil.RandPeerIDFatal(b),
		testutil.RandPeerIDFatal(b),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := cid.NewCidV0(u.Hash([]byte(fmt.Sprint(i))))

		q.Push(peers[i%len(peers)], &wantlist.Entry{Cid: c, Priority: math.MaxInt32})
	}
}
