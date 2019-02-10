package main

import (
	"flag"
	"fmt"
	"log"

	p2pd "gx/ipfs/QmTSBjpm7866oJKqcbuuicx6v3ypjgzoZ47y6xDNnWbB3o/go-libp2p-daemon"

	crypto "gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
)

func main() {
	file := flag.String("f", "identity", "output key file")
	ktype := flag.String("t", "rsa", "key type; rsa or ed25519")
	bits := flag.Int("b", 2048, "key size in bits (for rsa)")
	flag.Parse()

	var typ int

	switch *ktype {
	case "rsa":
		typ = crypto.RSA

	case "ed25519":
		typ = crypto.Ed25519

	default:
		log.Fatalf("Unknown key type %s; must be rsa or ed25519", *ktype)
	}

	priv, pub, err := crypto.GenerateKeyPair(typ, *bits)
	if err != nil {
		log.Fatal(err)
	}

	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Peer ID: %s\n", id.Pretty())

	err = p2pd.WriteIdentity(priv, *file)
	if err != nil {
		log.Fatal(err)
	}
}
