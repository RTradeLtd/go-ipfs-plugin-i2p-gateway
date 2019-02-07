package i2ptcpconn

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"

	"github.com/RTradeLtd/go-garlic-tcp-transport/codec"
	"github.com/RTradeLtd/go-garlic-tcp-transport/common"
	"github.com/RTradeLtd/sam3"
	"github.com/libp2p/go-stream-muxer"
)

// GarlicTCPConn implements a Conn interface
type GarlicTCPConn struct {
	*sam3.SAMConn
	*sam3.SAM
	*sam3.StreamSession
	*sam3.StreamListener
	parentTransport tpt.Transport
	laddr           ma.Multiaddr
	i2pkeys         *sam3.I2PKeys
	id              peer.ID
	rid             peer.ID

	hostSAM string
	portSAM string
	//passSAM string

	keysPath string

	onlyGarlic    bool
	garlicOptions []string
}

// SAMHost returns the IP address of the configured SAM bridge
func (t *GarlicTCPConn) SAMHost() string {
	st := strings.TrimPrefix(t.hostSAM, "/ip4/")
	stt := strings.TrimPrefix(st, "/ip6/")
	rt := strings.TrimSuffix(stt, "/")
	return rt
}

// SAMPort returns the Port of the configured SAM bridge
func (t *GarlicTCPConn) SAMPort() string {
	st := strings.TrimPrefix(t.portSAM, "/tcp/")
	rt := strings.TrimSuffix(st, "/")
	return rt
}

// SAMAddress combines them and returns a full address.
func (t GarlicTCPConn) SAMAddress() string {
	rt := t.SAMHost() + ":" + t.SAMPort()
	fmt.Println(rt)
	return rt
}

// PrintOptions returns the options passed to the SAM bridge as a slice of
// strings.
func (t GarlicTCPConn) PrintOptions() []string {
	return t.garlicOptions
}

// MaBase64 gives us a multiaddr by converting an I2PAddr
func (t GarlicTCPConn) MaBase64() ma.Multiaddr {
	r, err := i2ptcpcodec.FromI2PNetAddrToMultiaddr(t.i2pkeys.Addr())
	if err != nil {
		panic("Critical address error! There is no way this should have occurred")
	}
	return r
}

// Base32 returns the remotely-accessible base32 address of the gateway over i2p
// this is the one you want to use to visit it in the browser.
func (t GarlicTCPConn) Base32() string {
	return t.i2pkeys.Addr().Base32()
}

// Base64 returns the remotely-accessible base64 address of the gateway over I2P
func (t GarlicTCPConn) Base64() string {
	return t.i2pkeys.Addr().Base64()
}

// Transport returns the GarlicTCPTransport to which the GarlicTCPConn belongs
func (t GarlicTCPConn) Transport() tpt.Transport {
	return t.parentTransport
}

// IsClosed says a connection is closed if t.StreamSession is nil because
// Close() nils it if it works. Might need to re-visit that.
func (t GarlicTCPConn) IsClosed() bool {
	if t.StreamSession == nil {
		return true
	}
	return false
}

// AcceptStream lets us streammux
func (t GarlicTCPConn) AcceptStream() (streammux.Stream, error) {
	return t.AcceptI2P()
}

// Dial dials an I2P client connection to an i2p hidden service using a garlic64
// multiaddr and returns a tpt.Conn
func (t GarlicTCPConn) Dial(c context.Context, m ma.Multiaddr, p peer.ID) (tpt.Conn, error) {
	return t.DialI2P(c, m, p)
}

// DialI2P helps with Dial and returns a GarlicTCPConn
func (t GarlicTCPConn) DialI2P(c context.Context, m ma.Multiaddr, p peer.ID) (*GarlicTCPConn, error) {
	var err error
	t.SAMConn, err = t.StreamSession.DialContextI2P(c, "", m.String())
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// OpenStream lets us streammux
func (t GarlicTCPConn) OpenStream() (streammux.Stream, error) {
	return t.DialI2P(nil, t.RemoteMultiaddr(), t.RemotePeer())
}

// LocalMultiaddr returns the local multiaddr for this connection
func (t GarlicTCPConn) LocalMultiaddr() ma.Multiaddr {
	return t.laddr
}

// RemoteMultiaddr returns the remote multiaddr for this connection
func (t GarlicTCPConn) RemoteMultiaddr() ma.Multiaddr {
	return t.MaBase64()
}

// LocalPrivateKey returns the local private key used for the peer.ID
func (t GarlicTCPConn) LocalPrivateKey() crypto.PrivKey {
	return nil
}

// RemotePeer returns the remote peer.ID used for IPFS
func (t GarlicTCPConn) RemotePeer() peer.ID {
	return t.rid
}

// RemotePublicKey returns the remote public key used for communicating with the
// peer. It returns nil for now, security is provided solely by i2p for now.
func (t GarlicTCPConn) RemotePublicKey() crypto.PubKey {
	return nil
}

// LocalPeer returns the local peer.ID used for IPFS
func (t GarlicTCPConn) LocalPeer() peer.ID {
	return t.id
}

// Close ends a SAM session associated with a transport
func (t GarlicTCPConn) Close() error {
	err := t.StreamSession.Close()
	if err == nil {
		t.StreamSession = nil
	}
	return err
}

// Reset lets us streammux, I need to re-examine how to implement it.
func (t GarlicTCPConn) Reset() error {
	return t.Close()
}

// GetI2PKeys loads the i2p address keys and returns them.
func (t GarlicTCPConn) GetI2PKeys() (*sam3.I2PKeys, error) {
	return i2phelpers.LoadKeys(t.keysPath)
}

// Accept implements a listener
func (t GarlicTCPConn) Accept() (tpt.Conn, error) {
	return t.AcceptI2P()
}

// AcceptI2P helps with Accept
func (t GarlicTCPConn) AcceptI2P() (*GarlicTCPConn, error) {
	var err error
	t.SAMConn, err = t.StreamListener.AcceptI2P()
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Listen implements a listener
func (t GarlicTCPConn) Listen() (tpt.Conn, error) {
	return t.ListenI2P()
}

// ListenI2P helps with Listen
func (t GarlicTCPConn) ListenI2P() (*GarlicTCPConn, error) {
	var err error
	t.StreamListener, err = t.StreamSession.Listen()
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (t GarlicTCPConn) forward(conn *GarlicTCPConn) {
	//var request *http.Request
	var err error
	var client net.Conn
	if client, err = net.Dial("tcp", t.Addr().String()); err != nil {
		panic("Dial failed: %v" + err.Error())
	}
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

// Forward sets a local multiaddr and forwards the service at that address to
// I2P by calling ForwardI2P. You must already be listening on I2P before calling
// this function.
func (t GarlicTCPConn) Forward(addr ma.Multiaddr) {
	t.ForwardI2P(addr)
}

// ForwardI2P sets a local multiaddr and forwards the service at that address to
// I2P, it's what Forward calls.
func (t GarlicTCPConn) ForwardI2P(addr ma.Multiaddr) {
	var err error
	t.laddr = addr
	t.StreamListener, err = t.StreamSession.Listen()
	if err != nil {
		panic(err.Error())
	}
	for {
		conn, err := t.AcceptI2P()
		if err != nil {
			panic("ERROR: failed to accept listener: %v" + err.Error())
		}
		go t.forward(conn)
	}
}

// Addr returns the net.Addr version of the local Multiaddr
func (t GarlicTCPConn) Addr() net.Addr {
	ra, _ := manet.ToNetAddr(t.Multiaddr())
	return ra
}

// Multiaddr returns the local Multiaddr
func (t GarlicTCPConn) Multiaddr() ma.Multiaddr {
	return t.laddr
}

// NewGarlicTCPConn creates an I2P Connection struct from a fixed list of arguments
func NewGarlicTCPConn(transport tpt.Transport, host, port, pass string, keysPath string, onlyGarlic bool, options []string) (*GarlicTCPConn, error) {
	return NewGarlicTCPConnFromOptions(
		Transport(transport),
		SAMHost(host),
		SAMPort(port),
		SAMPass(pass),
		KeysPath(keysPath),
		OnlyGarlic(onlyGarlic),
		GarlicOptions(options),
	)
}

// NewGarlicTCPConnPeer creates an I2P Connection struct from a fixed list of
// arguments with a local peer.ID
func NewGarlicTCPConnPeer(transport tpt.Transport, id peer.ID, host, port, pass string, keysPath string, onlyGarlic bool, options []string) (*GarlicTCPConn, error) {
	return NewGarlicTCPConnFromOptions(
		Transport(transport),
		LocalPeerID(id),
		SAMHost(host),
		SAMPort(port),
		SAMPass(pass),
		KeysPath(keysPath),
		OnlyGarlic(onlyGarlic),
		GarlicOptions(options),
	)
}

// NewGarlicTCPConnFromOptions creates a GarlicTCPConn using function arguments
func NewGarlicTCPConnFromOptions(opts ...func(*GarlicTCPConn) error) (*GarlicTCPConn, error) {
	var t GarlicTCPConn
	t.hostSAM = "127.0.0.1"
	t.portSAM = "7656"
	//t.passSAM = ""
	t.keysPath = ""
	t.onlyGarlic = false
	t.garlicOptions = []string{}
	for _, o := range opts {
		if err := o(&t); err != nil {
			return nil, err
		}
	}
	var err error
	t.SAM, err = sam3.NewSAM(t.SAMAddress())
	if err != nil {
		return nil, err
	}
	t.i2pkeys, err = t.GetI2PKeys()
	if err != nil {
		return nil, err
	}
	t.StreamSession, err = t.SAM.NewStreamSession(i2phelpers.RandTunName(), *t.i2pkeys, t.PrintOptions())
	if err != nil {
		return nil, err
	}
	return &t, nil
}
