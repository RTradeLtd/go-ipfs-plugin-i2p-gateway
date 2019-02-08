package routinghelpers

import (
	"context"

	routing "gx/ipfs/QmRjT8Bkut84fHf9nxMQBxGsqLAkqzMdFaemDK7e61dBNZ/go-libp2p-routing"
	ropts "gx/ipfs/QmRjT8Bkut84fHf9nxMQBxGsqLAkqzMdFaemDK7e61dBNZ/go-libp2p-routing/options"

	ci "gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	pstore "gx/ipfs/QmQFFp4ntkd4C14sP3FaH9WJyBuetuGUVo6dShNHvnoEvC/go-libp2p-peerstore"
	cid "gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	record "gx/ipfs/QmexPd3srWxHC76gW2p5j5tQvwpPuCoW7b9vFhJ8BRPyh9/go-libp2p-record"
	multierror "gx/ipfs/QmfGQp6VVqdPCDyzEM6EGwMY74YPabTSEoQWHUxZuCSWj3/go-multierror"
)

// Tiered is like the Parallel except that GetValue and FindPeer
// are called in series.
type Tiered struct {
	Routers   []routing.IpfsRouting
	Validator record.Validator
}

func (r Tiered) PutValue(ctx context.Context, key string, value []byte, opts ...ropts.Option) error {
	return Parallel{Routers: r.Routers}.PutValue(ctx, key, value, opts...)
}

func (r Tiered) get(ctx context.Context, do func(routing.IpfsRouting) (interface{}, error)) (interface{}, error) {
	var errs []error
	for _, ri := range r.Routers {
		val, err := do(ri)
		switch err {
		case nil:
			return val, nil
		case routing.ErrNotFound, routing.ErrNotSupported:
			continue
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		errs = append(errs, err)
	}
	switch len(errs) {
	case 0:
		return nil, routing.ErrNotFound
	case 1:
		return nil, errs[0]
	default:
		return nil, &multierror.Error{Errors: errs}
	}
}

func (r Tiered) GetValue(ctx context.Context, key string, opts ...ropts.Option) ([]byte, error) {
	valInt, err := r.get(ctx, func(ri routing.IpfsRouting) (interface{}, error) {
		return ri.GetValue(ctx, key, opts...)
	})
	val, _ := valInt.([]byte)
	return val, err
}

func (r Tiered) SearchValue(ctx context.Context, key string, opts ...ropts.Option) (<-chan []byte, error) {
	return Parallel{Routers: r.Routers, Validator: r.Validator}.SearchValue(ctx, key, opts...)
}

func (r Tiered) GetPublicKey(ctx context.Context, p peer.ID) (ci.PubKey, error) {
	vInt, err := r.get(ctx, func(ri routing.IpfsRouting) (interface{}, error) {
		return routing.GetPublicKey(ri, ctx, p)
	})
	val, _ := vInt.(ci.PubKey)
	return val, err
}

func (r Tiered) Provide(ctx context.Context, c cid.Cid, local bool) error {
	return Parallel{Routers: r.Routers}.Provide(ctx, c, local)
}

func (r Tiered) FindProvidersAsync(ctx context.Context, c cid.Cid, count int) <-chan pstore.PeerInfo {
	return Parallel{Routers: r.Routers}.FindProvidersAsync(ctx, c, count)
}

func (r Tiered) FindPeer(ctx context.Context, p peer.ID) (pstore.PeerInfo, error) {
	valInt, err := r.get(ctx, func(ri routing.IpfsRouting) (interface{}, error) {
		return ri.FindPeer(ctx, p)
	})
	val, _ := valInt.(pstore.PeerInfo)
	return val, err
}

func (r Tiered) Bootstrap(ctx context.Context) error {
	return Parallel{Routers: r.Routers}.Bootstrap(ctx)
}

var _ routing.IpfsRouting = Tiered{}
