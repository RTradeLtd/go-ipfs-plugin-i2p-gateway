package ipns

import (
	"fmt"
	"testing"
	"time"

	ci "gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	u "gx/ipfs/QmNohiVssaPw3KVLZik59DBVGTSm2dGvYT9eoXt5DQ36Yz/go-ipfs-util"
	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
)

func TestEmbedPublicKey(t *testing.T) {

	sr := u.NewTimeSeededRand()
	priv, pub, err := ci.GenerateKeyPairWithReader(ci.RSA, 1024, sr)
	if err != nil {
		t.Fatal(err)
	}

	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}

	e, err := Create(priv, []byte("/a/b"), 0, time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if err := EmbedPublicKey(pub, e); err != nil {
		t.Fatal(err)
	}
	embeddedPk, err := ci.UnmarshalPublicKey(e.PubKey)
	if err != nil {
		t.Fatal(err)
	}
	embeddedPid, err := peer.IDFromPublicKey(embeddedPk)
	if err != nil {
		t.Fatal(err)
	}
	if embeddedPid != pid {
		t.Fatalf("pid mismatch: %s != %s", pid, embeddedPid)
	}
}

func ExampleCreate() {
	// Generate a private key to sign the IPNS record with. Most of the time,
	// however, you'll want to retrieve an already-existing key from IPFS using
	// go-ipfs/core/coreapi CoreAPI.KeyAPI() interface.
	privateKey, _, err := ci.GenerateKeyPair(ci.RSA, 2048)
	if err != nil {
		panic(err)
	}

	// Create an IPNS record that expires in one hour and points to the IPFS address
	// /ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5
	ipnsRecord, err := Create(privateKey, []byte("/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5"), 0, time.Now().Add(1*time.Hour))
	if err != nil {
		panic(err)
	}

	fmt.Println(ipnsRecord)
}
