// Package nilrouting implements a routing client that does nothing.
package nilrouting

import (
	"context"
	"errors"

	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	pstore "gx/ipfs/QmQFFp4ntkd4C14sP3FaH9WJyBuetuGUVo6dShNHvnoEvC/go-libp2p-peerstore"
	cid "gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	routing "gx/ipfs/QmRjT8Bkut84fHf9nxMQBxGsqLAkqzMdFaemDK7e61dBNZ/go-libp2p-routing"
	ropts "gx/ipfs/QmRjT8Bkut84fHf9nxMQBxGsqLAkqzMdFaemDK7e61dBNZ/go-libp2p-routing/options"
	record "gx/ipfs/QmexPd3srWxHC76gW2p5j5tQvwpPuCoW7b9vFhJ8BRPyh9/go-libp2p-record"
	ds "gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore"
	p2phost "gx/ipfs/QmfRHxh8bt4jWLKRhNvR5fn7mFACrQBFLqV4wyoymEExKV/go-libp2p-host"
)

type nilclient struct {
}

func (c *nilclient) PutValue(_ context.Context, _ string, _ []byte, _ ...ropts.Option) error {
	return nil
}

func (c *nilclient) GetValue(_ context.Context, _ string, _ ...ropts.Option) ([]byte, error) {
	return nil, errors.New("tried GetValue from nil routing")
}

func (c *nilclient) SearchValue(_ context.Context, _ string, _ ...ropts.Option) (<-chan []byte, error) {
	return nil, errors.New("tried SearchValue from nil routing")
}

func (c *nilclient) FindPeer(_ context.Context, _ peer.ID) (pstore.PeerInfo, error) {
	return pstore.PeerInfo{}, nil
}

func (c *nilclient) FindProvidersAsync(_ context.Context, _ cid.Cid, _ int) <-chan pstore.PeerInfo {
	out := make(chan pstore.PeerInfo)
	defer close(out)
	return out
}

func (c *nilclient) Provide(_ context.Context, _ cid.Cid, _ bool) error {
	return nil
}

func (c *nilclient) Bootstrap(_ context.Context) error {
	return nil
}

// ConstructNilRouting creates an IpfsRouting client which does nothing.
func ConstructNilRouting(_ context.Context, _ p2phost.Host, _ ds.Batching, _ record.Validator) (routing.IpfsRouting, error) {
	return &nilclient{}, nil
}

//  ensure nilclient satisfies interface
var _ routing.IpfsRouting = &nilclient{}
