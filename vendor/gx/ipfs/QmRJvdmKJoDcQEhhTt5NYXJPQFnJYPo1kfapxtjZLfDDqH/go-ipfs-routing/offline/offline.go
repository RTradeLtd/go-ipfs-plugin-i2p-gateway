// Package offline implements IpfsRouting with a client which
// is only able to perform offline operations.
package offline

import (
	"bytes"
	"context"
	"errors"
	"time"

	"gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	pstore "gx/ipfs/QmQFFp4ntkd4C14sP3FaH9WJyBuetuGUVo6dShNHvnoEvC/go-libp2p-peerstore"
	cid "gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	routing "gx/ipfs/QmRjT8Bkut84fHf9nxMQBxGsqLAkqzMdFaemDK7e61dBNZ/go-libp2p-routing"
	ropts "gx/ipfs/QmRjT8Bkut84fHf9nxMQBxGsqLAkqzMdFaemDK7e61dBNZ/go-libp2p-routing/options"
	dshelp "gx/ipfs/QmauEMWPoSqggfpSDHMMXuDn12DTd7TaFBvn39eeurzKT2/go-ipfs-ds-help"
	proto "gx/ipfs/QmdxUuburamoF6zF9qjeQC4WYcWGbWuRmdLacMEsW8ioD8/gogo-protobuf/proto"
	record "gx/ipfs/QmexPd3srWxHC76gW2p5j5tQvwpPuCoW7b9vFhJ8BRPyh9/go-libp2p-record"
	pb "gx/ipfs/QmexPd3srWxHC76gW2p5j5tQvwpPuCoW7b9vFhJ8BRPyh9/go-libp2p-record/pb"
	ds "gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore"
)

// ErrOffline is returned when trying to perform operations that
// require connectivity.
var ErrOffline = errors.New("routing system in offline mode")

// NewOfflineRouter returns an IpfsRouting implementation which only performs
// offline operations. It allows to Put and Get signed dht
// records to and from the local datastore.
func NewOfflineRouter(dstore ds.Datastore, validator record.Validator) routing.IpfsRouting {
	return &offlineRouting{
		datastore: dstore,
		validator: validator,
	}
}

// offlineRouting implements the IpfsRouting interface,
// but only provides the capability to Put and Get signed dht
// records to and from the local datastore.
type offlineRouting struct {
	datastore ds.Datastore
	validator record.Validator
}

func (c *offlineRouting) PutValue(ctx context.Context, key string, val []byte, _ ...ropts.Option) error {
	if err := c.validator.Validate(key, val); err != nil {
		return err
	}
	if old, err := c.GetValue(ctx, key); err == nil {
		// be idempotent to be nice.
		if bytes.Equal(old, val) {
			return nil
		}
		// check to see if the older record is better
		i, err := c.validator.Select(key, [][]byte{val, old})
		if err != nil {
			// this shouldn't happen for validated records.
			return err
		}
		if i != 0 {
			return errors.New("can't replace a newer record with an older one")
		}
	}
	rec := record.MakePutRecord(key, val)
	data, err := proto.Marshal(rec)
	if err != nil {
		return err
	}

	return c.datastore.Put(dshelp.NewKeyFromBinary([]byte(key)), data)
}

func (c *offlineRouting) GetValue(ctx context.Context, key string, _ ...ropts.Option) ([]byte, error) {
	buf, err := c.datastore.Get(dshelp.NewKeyFromBinary([]byte(key)))
	if err != nil {
		return nil, err
	}

	rec := new(pb.Record)
	err = proto.Unmarshal(buf, rec)
	if err != nil {
		return nil, err
	}
	val := rec.GetValue()

	err = c.validator.Validate(key, val)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (c *offlineRouting) SearchValue(ctx context.Context, key string, _ ...ropts.Option) (<-chan []byte, error) {
	out := make(chan []byte, 1)
	go func() {
		defer close(out)
		v, err := c.GetValue(ctx, key)
		if err == nil {
			out <- v
		}
	}()
	return out, nil
}

func (c *offlineRouting) FindPeer(ctx context.Context, pid peer.ID) (pstore.PeerInfo, error) {
	return pstore.PeerInfo{}, ErrOffline
}

func (c *offlineRouting) FindProvidersAsync(ctx context.Context, k cid.Cid, max int) <-chan pstore.PeerInfo {
	out := make(chan pstore.PeerInfo)
	close(out)
	return out
}

func (c *offlineRouting) Provide(_ context.Context, k cid.Cid, _ bool) error {
	return ErrOffline
}

func (c *offlineRouting) Ping(ctx context.Context, p peer.ID) (time.Duration, error) {
	return 0, ErrOffline
}

func (c *offlineRouting) Bootstrap(context.Context) error {
	return nil
}

// ensure offlineRouting matches the IpfsRouting interface
var _ routing.IpfsRouting = &offlineRouting{}
