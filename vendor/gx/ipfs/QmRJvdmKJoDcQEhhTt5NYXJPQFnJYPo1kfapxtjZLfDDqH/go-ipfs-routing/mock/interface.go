// Package mockrouting provides a virtual routing server. To use it,
// create a virtual routing server and use the Client() method to get a
// routing client (IpfsRouting). The server quacks like a DHT but is
// really a local in-memory hash table.
package mockrouting

import (
	"context"

	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	routing "gx/ipfs/QmRjT8Bkut84fHf9nxMQBxGsqLAkqzMdFaemDK7e61dBNZ/go-libp2p-routing"
	delay "gx/ipfs/QmUe1WCHkQaz4UeNKiHDUBV2T6i9prc3DniqyHPXyfGaUq/go-ipfs-delay"
	"gx/ipfs/QmVnJMgafh5MBYiyqbvDtoCL8pcQvbEGD2k9o9GFpBWPzY/go-testutil"
	ds "gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore"
)

// MockValidator is a record validator that always returns success.
type MockValidator struct{}

func (MockValidator) Validate(_ string, _ []byte) error        { return nil }
func (MockValidator) Select(_ string, _ [][]byte) (int, error) { return 0, nil }

// Server provides mockrouting Clients
type Server interface {
	Client(p testutil.Identity) Client
	ClientWithDatastore(context.Context, testutil.Identity, ds.Datastore) Client
}

// Client implements IpfsRouting
type Client interface {
	routing.IpfsRouting
}

// NewServer returns a mockrouting Server
func NewServer() Server {
	return NewServerWithDelay(DelayConfig{
		ValueVisibility: delay.Fixed(0),
		Query:           delay.Fixed(0),
	})
}

// NewServerWithDelay returns a mockrouting Server with a delay!
func NewServerWithDelay(conf DelayConfig) Server {
	return &s{
		providers: make(map[string]map[peer.ID]providerRecord),
		delayConf: conf,
	}
}

// DelayConfig can be used to configured the fake delays of a mock server.
// Use with NewServerWithDelay().
type DelayConfig struct {
	// ValueVisibility is the time it takes for a value to be visible in the network
	// FIXME there _must_ be a better term for this
	ValueVisibility delay.D

	// Query is the time it takes to receive a response from a routing query
	Query delay.D
}
