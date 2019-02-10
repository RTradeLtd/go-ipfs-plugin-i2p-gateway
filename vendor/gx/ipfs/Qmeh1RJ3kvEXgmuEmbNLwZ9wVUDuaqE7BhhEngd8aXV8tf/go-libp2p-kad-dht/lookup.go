package dht

import (
	"context"
	"fmt"
	"strings"

	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	cid "gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	notif "gx/ipfs/QmRjT8Bkut84fHf9nxMQBxGsqLAkqzMdFaemDK7e61dBNZ/go-libp2p-routing/notifications"
	kb "gx/ipfs/QmWNfpvvoMom3j19zgvSh1VMxNuHbLDXTPkxUf8Gpju5H7/go-libp2p-kbucket"
	logging "gx/ipfs/QmcuXC5cxs79ro2cUuHs4HQ2bkDLJUYokwL8aivcX6HW3C/go-log"
	pb "gx/ipfs/Qmeh1RJ3kvEXgmuEmbNLwZ9wVUDuaqE7BhhEngd8aXV8tf/go-libp2p-kad-dht/pb"
)

func tryFormatLoggableKey(k string) (string, error) {
	if len(k) == 0 {
		return "", fmt.Errorf("loggableKey is empty")
	}
	var proto, cstr string
	if k[0] == '/' {
		// it's a path (probably)
		protoEnd := strings.IndexByte(k[1:], '/')
		if protoEnd < 0 {
			return k, fmt.Errorf("loggableKey starts with '/' but is not a path: %x", k)
		}
		proto = k[1 : protoEnd+1]
		cstr = k[protoEnd+2:]
	} else {
		proto = "provider"
		cstr = k
	}

	c, err := cid.Cast([]byte(cstr))
	if err != nil {
		return "", fmt.Errorf("loggableKey could not cast key to a CID: %x %v", k, err)
	}
	return fmt.Sprintf("/%s/%s", proto, c.String()), nil
}

func loggableKey(k string) logging.LoggableMap {
	newKey, err := tryFormatLoggableKey(k)
	if err != nil {
		logger.Debug(err)
	} else {
		k = newKey
	}

	return logging.LoggableMap{
		"key": k,
	}
}

// Kademlia 'node lookup' operation. Returns a channel of the K closest peers
// to the given key
func (dht *IpfsDHT) GetClosestPeers(ctx context.Context, key string) (<-chan peer.ID, error) {
	e := logger.EventBegin(ctx, "getClosestPeers", loggableKey(key))
	tablepeers := dht.routingTable.NearestPeers(kb.ConvertKey(key), AlphaValue)
	if len(tablepeers) == 0 {
		return nil, kb.ErrLookupFailure
	}

	out := make(chan peer.ID, KValue)

	// since the query doesnt actually pass our context down
	// we have to hack this here. whyrusleeping isnt a huge fan of goprocess
	parent := ctx
	query := dht.newQuery(key, func(ctx context.Context, p peer.ID) (*dhtQueryResult, error) {
		// For DHT query command
		notif.PublishQueryEvent(parent, &notif.QueryEvent{
			Type: notif.SendingQuery,
			ID:   p,
		})

		pmes, err := dht.findPeerSingle(ctx, p, peer.ID(key))
		if err != nil {
			logger.Debugf("error getting closer peers: %s", err)
			return nil, err
		}
		peers := pb.PBPeersToPeerInfos(pmes.GetCloserPeers())

		// For DHT query command
		notif.PublishQueryEvent(parent, &notif.QueryEvent{
			Type:      notif.PeerResponse,
			ID:        p,
			Responses: peers,
		})

		return &dhtQueryResult{closerPeers: peers}, nil
	})

	go func() {
		defer close(out)
		defer e.Done()
		// run it!
		res, err := query.Run(ctx, tablepeers)
		if err != nil {
			logger.Debugf("closestPeers query run error: %s", err)
		}

		if res != nil && res.queriedSet != nil {
			sorted := kb.SortClosestPeers(res.queriedSet.Peers(), kb.ConvertKey(key))
			if len(sorted) > KValue {
				sorted = sorted[:KValue]
			}

			for _, p := range sorted {
				out <- p
			}
		}
	}()

	return out, nil
}
