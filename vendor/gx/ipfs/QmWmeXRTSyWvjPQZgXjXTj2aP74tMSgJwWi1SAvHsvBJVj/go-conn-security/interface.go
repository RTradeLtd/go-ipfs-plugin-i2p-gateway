package connsec

import (
	"context"
	"net"

	peer "gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	inet "gx/ipfs/QmZ7cBWUXkyWTMN4qH6NGoyMVs7JugyFChBNP4ZUp5rJHH/go-libp2p-net"
)

// A Transport turns inbound and outbound unauthenticated,
// plain-text connections into authenticated, encrypted connections.
type Transport interface {
	// SecureInbound secures an inbound connection.
	SecureInbound(ctx context.Context, insecure net.Conn) (Conn, error)

	// SecureOutbound secures an outbound connection.
	SecureOutbound(ctx context.Context, insecure net.Conn, p peer.ID) (Conn, error)
}

// Conn is an authenticated, encrypted connection.
type Conn interface {
	net.Conn
	inet.ConnSecurity
}
