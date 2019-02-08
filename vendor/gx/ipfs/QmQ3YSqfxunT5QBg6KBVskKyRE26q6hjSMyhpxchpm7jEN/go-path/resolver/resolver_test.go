package resolver_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	path "gx/ipfs/QmQ3YSqfxunT5QBg6KBVskKyRE26q6hjSMyhpxchpm7jEN/go-path"
	"gx/ipfs/QmQ3YSqfxunT5QBg6KBVskKyRE26q6hjSMyhpxchpm7jEN/go-path/resolver"

	ipld "gx/ipfs/QmRL22E4paat7ky7vx9MLpR97JHHbFPrg3ytFQw6qp1y1s/go-ipld-format"
	merkledag "gx/ipfs/QmUtsx89yiCY6F8mbpP6ecXckiSzCBH7EvkKZuZEHBcr1m/go-merkledag"
	dagmock "gx/ipfs/QmUtsx89yiCY6F8mbpP6ecXckiSzCBH7EvkKZuZEHBcr1m/go-merkledag/test"
)

func randNode() *merkledag.ProtoNode {
	node := new(merkledag.ProtoNode)
	node.SetData(make([]byte, 32))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Read(node.Data())
	return node
}

func TestRecurivePathResolution(t *testing.T) {
	ctx := context.Background()
	dagService := dagmock.Mock()

	a := randNode()
	b := randNode()
	c := randNode()

	err := b.AddNodeLink("grandchild", c)
	if err != nil {
		t.Fatal(err)
	}

	err = a.AddNodeLink("child", b)
	if err != nil {
		t.Fatal(err)
	}

	for _, n := range []ipld.Node{a, b, c} {
		err = dagService.Add(ctx, n)
		if err != nil {
			t.Fatal(err)
		}
	}

	aKey := a.Cid()

	segments := []string{aKey.String(), "child", "grandchild"}
	p, err := path.FromSegments("/ipfs/", segments...)
	if err != nil {
		t.Fatal(err)
	}

	resolver := resolver.NewBasicResolver(dagService)
	node, err := resolver.ResolvePath(ctx, p)
	if err != nil {
		t.Fatal(err)
	}

	cKey := c.Cid()
	key := node.Cid()
	if key.String() != cKey.String() {
		t.Fatal(fmt.Errorf(
			"recursive path resolution failed for %s: %s != %s",
			p.String(), key.String(), cKey.String()))
	}

	rCid, rest, err := resolver.ResolveToLastNode(ctx, p)
	if err != nil {
		t.Fatal(err)
	}

	if len(rest) != 0 {
		t.Error("expected rest to be empty")
	}

	if rCid.String() != cKey.String() {
		t.Fatal(fmt.Errorf(
			"ResolveToLastNode failed for %s: %s != %s",
			p.String(), rCid.String(), cKey.String()))
	}

	p2, err := path.FromSegments("/ipfs/", aKey.String())
	if err != nil {
		t.Fatal(err)
	}

	rCid, rest, err = resolver.ResolveToLastNode(ctx, p2)
	if err != nil {
		t.Fatal(err)
	}


	if len(rest) != 0 {
		t.Error("expected rest to be empty")
	}

	if rCid.String() != aKey.String() {
		t.Fatal(fmt.Errorf(
			"ResolveToLastNode failed for %s: %s != %s",
			p.String(), rCid.String(), cKey.String()))
	}
}
