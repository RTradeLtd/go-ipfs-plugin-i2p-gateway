package bitswap

import (
	"context"

	bsnet "gx/ipfs/QmYJ48z7NEzo3u2yCvUvNtBQ7wJWd5dX2nxxc7FeA6nHq1/go-bitswap/network"

	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	mockrouting "gx/ipfs/QmRJvdmKJoDcQEhhTt5NYXJPQFnJYPo1kfapxtjZLfDDqH/go-ipfs-routing/mock"
	mockpeernet "gx/ipfs/QmSgtf5vHyugoxcwMbyNy6bZ9qPDDTJSYEED2GkWjLwitZ/go-libp2p/p2p/net/mock"
	testutil "gx/ipfs/QmVnJMgafh5MBYiyqbvDtoCL8pcQvbEGD2k9o9GFpBWPzY/go-testutil"
	ds "gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore"
)

type peernet struct {
	mockpeernet.Mocknet
	routingserver mockrouting.Server
}

func StreamNet(ctx context.Context, net mockpeernet.Mocknet, rs mockrouting.Server) (Network, error) {
	return &peernet{net, rs}, nil
}

func (pn *peernet) Adapter(p testutil.Identity) bsnet.BitSwapNetwork {
	client, err := pn.Mocknet.AddPeer(p.PrivateKey(), p.Address())
	if err != nil {
		panic(err.Error())
	}
	routing := pn.routingserver.ClientWithDatastore(context.TODO(), p, ds.NewMapDatastore())
	return bsnet.NewFromIpfsHost(client, routing)
}

func (pn *peernet) HasPeer(p peer.ID) bool {
	for _, member := range pn.Mocknet.Peers() {
		if p == member {
			return true
		}
	}
	return false
}

var _ Network = (*peernet)(nil)
