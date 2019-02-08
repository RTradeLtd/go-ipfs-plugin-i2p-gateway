package dht

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	ci "gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	u "gx/ipfs/QmNohiVssaPw3KVLZik59DBVGTSm2dGvYT9eoXt5DQ36Yz/go-ipfs-util"
	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	routing "gx/ipfs/QmRjT8Bkut84fHf9nxMQBxGsqLAkqzMdFaemDK7e61dBNZ/go-libp2p-routing"
	record "gx/ipfs/QmexPd3srWxHC76gW2p5j5tQvwpPuCoW7b9vFhJ8BRPyh9/go-libp2p-record"
)

// Check that GetPublicKey() correctly extracts a public key
func TestPubkeyExtract(t *testing.T) {
	t.Skip("public key extraction for ed25519 keys has been disabled. See https://github.com/libp2p/specs/issues/111")
	ctx := context.Background()
	dht := setupDHT(ctx, t, false)
	defer dht.Close()

	_, pk, err := ci.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	pid, err := peer.IDFromPublicKey(pk)
	if err != nil {
		t.Fatal(err)
	}

	pkOut, err := dht.GetPublicKey(context.Background(), pid)
	if err != nil {
		t.Fatal(err)
	}

	if !pkOut.Equals(pk) {
		t.Fatal("got incorrect public key out")
	}
}

// Check that GetPublicKey() correctly retrieves a public key from the peerstore
func TestPubkeyPeerstore(t *testing.T) {
	ctx := context.Background()
	dht := setupDHT(ctx, t, false)

	r := u.NewSeededRand(15) // generate deterministic keypair
	_, pubk, err := ci.GenerateKeyPairWithReader(ci.RSA, 1024, r)
	if err != nil {
		t.Fatal(err)
	}
	id, err := peer.IDFromPublicKey(pubk)
	if err != nil {
		t.Fatal(err)
	}
	err = dht.peerstore.AddPubKey(id, pubk)
	if err != nil {
		t.Fatal(err)
	}

	rpubk, err := dht.GetPublicKey(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}

	if !pubk.Equals(rpubk) {
		t.Fatal("got incorrect public key")
	}
}

// Check that GetPublicKey() correctly retrieves a public key directly
// from the node it identifies
func TestPubkeyDirectFromNode(t *testing.T) {
	ctx := context.Background()

	dhtA := setupDHT(ctx, t, false)
	dhtB := setupDHT(ctx, t, false)

	defer dhtA.Close()
	defer dhtB.Close()
	defer dhtA.host.Close()
	defer dhtB.host.Close()

	connect(t, ctx, dhtA, dhtB)

	pubk, err := dhtA.GetPublicKey(context.Background(), dhtB.self)
	if err != nil {
		t.Fatal(err)
	}

	id, err := peer.IDFromPublicKey(pubk)
	if err != nil {
		t.Fatal(err)
	}

	if id != dhtB.self {
		t.Fatal("got incorrect public key")
	}
}

// Check that GetPublicKey() correctly retrieves a public key
// from the DHT
func TestPubkeyFromDHT(t *testing.T) {
	ctx := context.Background()

	dhtA := setupDHT(ctx, t, false)
	dhtB := setupDHT(ctx, t, false)

	defer dhtA.Close()
	defer dhtB.Close()
	defer dhtA.host.Close()
	defer dhtB.host.Close()

	connect(t, ctx, dhtA, dhtB)

	r := u.NewSeededRand(15) // generate deterministic keypair
	_, pubk, err := ci.GenerateKeyPairWithReader(ci.RSA, 1024, r)
	if err != nil {
		t.Fatal(err)
	}
	id, err := peer.IDFromPublicKey(pubk)
	if err != nil {
		t.Fatal(err)
	}
	pkkey := routing.KeyForPublicKey(id)
	pkbytes, err := pubk.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	// Store public key on node B
	err = dhtB.PutValue(ctx, pkkey, pkbytes)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve public key on node A
	rpubk, err := dhtA.GetPublicKey(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if !pubk.Equals(rpubk) {
		t.Fatal("got incorrect public key")
	}
}

// Check that GetPublicKey() correctly returns an error when the
// public key is not available directly from the node or on the DHT
func TestPubkeyNotFound(t *testing.T) {
	ctx := context.Background()

	dhtA := setupDHT(ctx, t, false)
	dhtB := setupDHT(ctx, t, false)

	defer dhtA.Close()
	defer dhtB.Close()
	defer dhtA.host.Close()
	defer dhtB.host.Close()

	connect(t, ctx, dhtA, dhtB)

	r := u.NewSeededRand(15) // generate deterministic keypair
	_, pubk, err := ci.GenerateKeyPairWithReader(ci.RSA, 1024, r)
	if err != nil {
		t.Fatal(err)
	}
	id, err := peer.IDFromPublicKey(pubk)
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to retrieve public key on node A (should be not found)
	_, err = dhtA.GetPublicKey(ctx, id)
	if err == nil {
		t.Fatal("Expected not found error")
	}
}

// Check that GetPublicKey() returns an error when
// the DHT returns the wrong key
func TestPubkeyBadKeyFromDHT(t *testing.T) {
	ctx := context.Background()

	dhtA := setupDHT(ctx, t, false)
	dhtB := setupDHT(ctx, t, false)

	defer dhtA.Close()
	defer dhtB.Close()
	defer dhtA.host.Close()
	defer dhtB.host.Close()

	connect(t, ctx, dhtA, dhtB)

	r := u.NewSeededRand(15) // generate deterministic keypair
	_, pubk, err := ci.GenerateKeyPairWithReader(ci.RSA, 1024, r)
	if err != nil {
		t.Fatal(err)
	}
	id, err := peer.IDFromPublicKey(pubk)
	if err != nil {
		t.Fatal(err)
	}
	pkkey := routing.KeyForPublicKey(id)

	_, wrongpubk, err := ci.GenerateKeyPairWithReader(ci.RSA, 1024, r)
	if err != nil {
		t.Fatal(err)
	}
	if pubk == wrongpubk {
		t.Fatal("Public keys shouldn't match here")
	}
	wrongbytes, err := wrongpubk.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	// Store incorrect public key on node B
	rec := record.MakePutRecord(pkkey, wrongbytes)
	rec.TimeReceived = u.FormatRFC3339(time.Now())
	err = dhtB.putLocal(pkkey, rec)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve public key from node A
	_, err = dhtA.GetPublicKey(ctx, id)
	if err == nil {
		t.Fatal("Expected error because public key is incorrect")
	}
}

// Check that GetPublicKey() returns the correct value
// when the DHT returns the wrong key but the direct
// connection returns the correct key
func TestPubkeyBadKeyFromDHTGoodKeyDirect(t *testing.T) {
	ctx := context.Background()

	dhtA := setupDHT(ctx, t, false)
	dhtB := setupDHT(ctx, t, false)

	defer dhtA.Close()
	defer dhtB.Close()
	defer dhtA.host.Close()
	defer dhtB.host.Close()

	connect(t, ctx, dhtA, dhtB)

	r := u.NewSeededRand(15) // generate deterministic keypair
	pkkey := routing.KeyForPublicKey(dhtB.self)

	_, wrongpubk, err := ci.GenerateKeyPairWithReader(ci.RSA, 1024, r)
	if err != nil {
		t.Fatal(err)
	}
	wrongbytes, err := wrongpubk.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	// Store incorrect public key on node B
	rec := record.MakePutRecord(pkkey, wrongbytes)
	rec.TimeReceived = u.FormatRFC3339(time.Now())
	err = dhtB.putLocal(pkkey, rec)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve public key from node A
	pubk, err := dhtA.GetPublicKey(ctx, dhtB.self)
	if err != nil {
		t.Fatal(err)
	}

	id, err := peer.IDFromPublicKey(pubk)
	if err != nil {
		t.Fatal(err)
	}

	// The incorrect public key retrieved from the DHT
	// should be ignored in favour of the correct public
	// key retieved from the node directly
	if id != dhtB.self {
		t.Fatal("got incorrect public key")
	}
}

// Check that GetPublicKey() returns the correct value
// when both the DHT returns the correct key and the direct
// connection returns the correct key
func TestPubkeyGoodKeyFromDHTGoodKeyDirect(t *testing.T) {
	ctx := context.Background()

	dhtA := setupDHT(ctx, t, false)
	dhtB := setupDHT(ctx, t, false)

	defer dhtA.Close()
	defer dhtB.Close()
	defer dhtA.host.Close()
	defer dhtB.host.Close()

	connect(t, ctx, dhtA, dhtB)

	pubk := dhtB.peerstore.PubKey(dhtB.self)
	pkbytes, err := pubk.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	// Store public key on node B
	pkkey := routing.KeyForPublicKey(dhtB.self)
	err = dhtB.PutValue(ctx, pkkey, pkbytes)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve public key on node A
	rpubk, err := dhtA.GetPublicKey(ctx, dhtB.self)
	if err != nil {
		t.Fatal(err)
	}

	if !pubk.Equals(rpubk) {
		t.Fatal("got incorrect public key")
	}
}