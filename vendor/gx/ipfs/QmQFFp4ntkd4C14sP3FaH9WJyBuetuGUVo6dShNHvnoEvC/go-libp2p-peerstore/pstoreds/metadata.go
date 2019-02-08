package pstoreds

import (
	"bytes"
	"context"
	"encoding/gob"

	base32 "gx/ipfs/QmfVj3x4D6Jkq9SEoi5n2NmoUomLwoeiwnYz2KQa15wRw6/base32"

	ds "gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore"

	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	pool "gx/ipfs/QmQDvJoB6aJWN3sjr3xsgXqKCXf4jU5zdMXpDMsBkYVNqa/go-buffer-pool"
	pstore "gx/ipfs/QmQFFp4ntkd4C14sP3FaH9WJyBuetuGUVo6dShNHvnoEvC/go-libp2p-peerstore"
)

// Metadata is stored under the following db key pattern:
// /peers/metadata/<b32 peer id no padding>/<key>
var pmBase = ds.NewKey("/peers/metadata")

type dsPeerMetadata struct {
	ds ds.Datastore
}

var _ pstore.PeerMetadata = (*dsPeerMetadata)(nil)

func init() {
	// Gob registers basic types by default.
	//
	// Register complex types used by the peerstore itself.
	gob.Register(make(map[string]struct{}))
}

// NewPeerMetadata creates a metadata store backed by a persistent db. It uses gob for serialisation.
//
// See `init()` to learn which types are registered by default. Modules wishing to store
// values of other types will need to `gob.Register()` them explicitly, or else callers
// will receive runtime errors.
func NewPeerMetadata(_ context.Context, store ds.Datastore, _ Options) (pstore.PeerMetadata, error) {
	return &dsPeerMetadata{store}, nil
}

func (pm *dsPeerMetadata) Get(p peer.ID, key string) (interface{}, error) {
	k := pmBase.ChildString(base32.RawStdEncoding.EncodeToString([]byte(p))).ChildString(key)
	value, err := pm.ds.Get(k)
	if err != nil {
		if err == ds.ErrNotFound {
			err = pstore.ErrNotFound
		}
		return nil, err
	}

	var res interface{}
	if err := gob.NewDecoder(bytes.NewReader(value)).Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}

func (pm *dsPeerMetadata) Put(p peer.ID, key string, val interface{}) error {
	k := pmBase.ChildString(base32.RawStdEncoding.EncodeToString([]byte(p))).ChildString(key)
	var buf pool.Buffer
	if err := gob.NewEncoder(&buf).Encode(&val); err != nil {
		return err
	}
	return pm.ds.Put(k, buf.Bytes())
}
