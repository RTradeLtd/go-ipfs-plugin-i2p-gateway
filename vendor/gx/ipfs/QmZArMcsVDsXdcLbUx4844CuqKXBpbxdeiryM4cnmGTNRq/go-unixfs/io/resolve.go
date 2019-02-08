package io

import (
	"context"

	dag "gx/ipfs/QmUtsx89yiCY6F8mbpP6ecXckiSzCBH7EvkKZuZEHBcr1m/go-merkledag"
	ft "gx/ipfs/QmZArMcsVDsXdcLbUx4844CuqKXBpbxdeiryM4cnmGTNRq/go-unixfs"
	hamt "gx/ipfs/QmZArMcsVDsXdcLbUx4844CuqKXBpbxdeiryM4cnmGTNRq/go-unixfs/hamt"

	ipld "gx/ipfs/QmRL22E4paat7ky7vx9MLpR97JHHbFPrg3ytFQw6qp1y1s/go-ipld-format"
)

// ResolveUnixfsOnce resolves a single hop of a path through a graph in a
// unixfs context. This includes handling traversing sharded directories.
func ResolveUnixfsOnce(ctx context.Context, ds ipld.NodeGetter, nd ipld.Node, names []string) (*ipld.Link, []string, error) {
	pn, ok := nd.(*dag.ProtoNode)
	if ok {
		fsn, err := ft.FSNodeFromBytes(pn.Data())
		if err != nil {
			// Not a unixfs node, use standard object traversal code
			return nd.ResolveLink(names)
		}

		if fsn.Type() == ft.THAMTShard {
			rods := dag.NewReadOnlyDagService(ds)
			s, err := hamt.NewHamtFromDag(rods, nd)
			if err != nil {
				return nil, nil, err
			}

			out, err := s.Find(ctx, names[0])
			if err != nil {
				return nil, nil, err
			}

			return out, names[1:], nil
		}
	}

	return nd.ResolveLink(names)
}
