package cidutil

import (
	"math/rand"
	"testing"

	cid "gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	mhash "gx/ipfs/QmerPMzPk1mJVowm8KgmoknWa4yCYvvugMPsgWmDNUvDLW/go-multihash"
)

func TestInlineBuilderSmallValue(t *testing.T) {
	builder := InlineBuilder{cid.V0Builder{}, 64}
	c, err := builder.Sum([]byte("Hello World"))
	if err != nil {
		t.Fatal(err)
	}
	if c.Prefix().MhType != mhash.ID {
		t.Fatal("Inliner builder failed to use ID Multihash on small values")
	}
}

func TestInlinerBuilderLargeValue(t *testing.T) {
	builder := InlineBuilder{cid.V0Builder{}, 64}
	data := make([]byte, 512)
	rand.Read(data)
	c, err := builder.Sum(data)
	if err != nil {
		t.Fatal(err)
	}
	if c.Prefix().MhType == mhash.ID {
		t.Fatal("Inliner builder used ID Multihash on large values")
	}
}
