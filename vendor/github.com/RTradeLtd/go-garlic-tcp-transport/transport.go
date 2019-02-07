package i2ptcp

import (
	"context"
	"strings"

	peer "github.com/libp2p/go-libp2p-peer"

	"github.com/RTradeLtd/go-garlic-tcp-transport/common"
	"github.com/RTradeLtd/go-garlic-tcp-transport/conn"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
)

// GarlicTCPTransport is a libp2p interface to an i2p TCP-like tunnel created
// via the SAM bridge
type GarlicTCPTransport struct {
	hostSAM string
	portSAM string
	passSAM string
	id      peer.ID

	keysPath string

	onlyGarlic    bool
	garlicOptions []string
}

func (t *GarlicTCPTransport) SAMHost() string {
	st := strings.TrimPrefix(t.hostSAM, "/ip4/")
	stt := strings.TrimPrefix(st, "/ip6/")
	rt := strings.TrimSuffix(stt, "/")
	return rt
}

func (t *GarlicTCPTransport) SAMPort() string {
	st := strings.TrimPrefix(t.portSAM, "/tcp/")
	rt := strings.TrimSuffix(st, "/")
	return rt
}

func (t GarlicTCPTransport) SAMAddress() string {
	return t.SAMHost() + ":" + t.SAMPort()
}

func (t GarlicTCPTransport) PrintOptions() []string {
	return t.garlicOptions
}

// CanDial implements transport.CanDial
func (t GarlicTCPTransport) CanDial(m ma.Multiaddr) bool {
	return t.Matches(m)
}

// CanDialI2P is a special CanDial function that only returns true if it's an
// i2p address.
func (t GarlicTCPTransport) CanDialI2P(m ma.Multiaddr) bool {
	return t.MatchesI2P(m)
}

// Matches returns true if the address is a valid garlic TCP multiaddr
func (t *GarlicTCPTransport) Matches(a ma.Multiaddr) bool {
	return i2phelpers.IsValidGarlicMultiAddr(a)
}

// Matches returns true if the address is a valid garlic TCP multiaddr
func (t *GarlicTCPTransport) MatchesI2P(a ma.Multiaddr) bool {
	return i2phelpers.IsValidGarlicMultiAddr(a)
}

// Dial returns a new GarlicConn
func (t GarlicTCPTransport) Dial(c context.Context, m ma.Multiaddr, p peer.ID) (tpt.Conn, error) {
	conn, err := i2ptcpconn.NewGarlicTCPConn(t, t.hostSAM, t.portSAM, t.passSAM, t.keysPath, t.onlyGarlic, t.PrintOptions())
	if err != nil {
		return nil, err
	}
	return conn.DialI2P(c, m, p)
}

// Listen implements a connection, but addr is IGNORED here, it's drawn from the
//transport keys
func (t GarlicTCPTransport) Listen(addr ma.Multiaddr) (tpt.Listener, error) {
	return t.ListenI2P()
}

// ListenI2P is like Listen, but it returns the GarlicTCPConn and doesn't
//require a multiaddr
func (t GarlicTCPTransport) ListenI2P() (*i2ptcpconn.GarlicTCPConn, error) {
	conn, err := i2ptcpconn.NewGarlicTCPConnPeer(t, t.id, t.SAMHost(), t.SAMPort(), t.passSAM, t.keysPath, t.onlyGarlic, t.PrintOptions())
	if err != nil {
		return nil, err
	}
	return conn.ListenI2P()
}

// Protocols need only return this I think
func (t GarlicTCPTransport) Protocols() []int {
	return []int{ma.P_GARLIC64}
}

// Proxy always returns false, we're using the SAM bridge to make our requests
func (t GarlicTCPTransport) Proxy() bool {
	return false
}

// NewGarlicTransport initializes a GarlicTransport for libp2p
func NewGarlicTCPTransport(host, port, pass string, keysPath string, onlyGarlic bool, options []string) (tpt.Transport, error) {
	return NewGarlicTCPTransportFromOptions(
		SAMHost(host),
		SAMPort(port),
		SAMPass(pass),
		KeysPath(keysPath),
		OnlyGarlic(onlyGarlic),
		GarlicOptions(options),
	)
}

// NewGarlicTransportPeer initializes a GarlicTransport for libp2p with a local peer.ID
func NewGarlicTCPTransportPeer(id peer.ID, host, port, pass string, keysPath string, onlyGarlic bool, options []string) (tpt.Transport, error) {
	return NewGarlicTCPTransportFromOptions(
		LocalPeerID(id),
		SAMHost(host),
		SAMPort(port),
		SAMPass(pass),
		KeysPath(keysPath),
		OnlyGarlic(onlyGarlic),
		GarlicOptions(options),
	)
}

func NewGarlicTCPTransportFromOptions(opts ...func(*GarlicTCPTransport) error) (*GarlicTCPTransport, error) {
	var g GarlicTCPTransport
	g.hostSAM = "127.0.0.1"
	g.portSAM = "7656"
	g.passSAM = ""
	g.keysPath = ""
	g.onlyGarlic = false
	g.garlicOptions = []string{}
	for _, o := range opts {
		if err := o(&g); err != nil {
			return nil, err
		}
	}
	return &g, nil
}
