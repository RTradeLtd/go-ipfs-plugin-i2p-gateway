package i2pgate

import (
	"os"
    "strings"
	"github.com/rtradeltd/go-garlic-tcp-transport"
	//"github.com/rtradeltd/go-garlic-tcp-transport/conn"
	"github.com/rtradeltd/go-ipfs-plugin-i2p-gateway/config"

	//TODO: Get a better understanding of gx.
	config "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-config"
	plugin "github.com/ipsn/go-ipfs/plugin"
	fsrepo "github.com/ipsn/go-ipfs/repo/fsrepo"
    peer "github.com/libp2p/go-libp2p-peer")

type I2PGatePlugin struct {
	//*i2ptcp.GarlicTCPTransport
	//*i2ptcpconn.GarlicTCPConn
	configPath    string
	config        *config.Config
	i2pconfigPath string
	i2pconfig     *i2pgateconfig.Config
    id peer.ID

	forwardHTTP string
	forwardRPC  string
}

// I2PType will be used to identify this as the i2p gateway plugin to things
// that use it.
var I2PType = "i2pgate"

var _ plugin.Plugin = (*I2PGatePlugin)(nil)

// Name returns the plugin's name, satisfying the plugin.Plugin interface.
func (*I2PGatePlugin) Name() string {
	return "fwd-i2pgate"
}

// Version returns the plugin's version, satisfying the plugin.Plugin interface.
func (*I2PGatePlugin) Version() string {
	return "0.0.0"
}

// Init initializes plugin, satisfying the plugin.Plugin interface. Put any
// initialization logic here.
func (i *I2PGatePlugin) Init() error {
	var err error
	i.configPath, err = fsrepo.BestKnownPath()
	if err != nil {
		return err
	}
	err = os.Setenv("KEYS_PATH", i.configPath)
	if err != nil {
		return err
	}
	i.config, err = fsrepo.ConfigAt(i.configPath)
	if err != nil {
		return err
	}
	rpcaddressbytes, err := i.config.Addresses.API.MarshalJSON()
	if err != nil {
		return err
	}
	i.forwardRPC = strings.Replace(string(rpcaddressbytes), "\"", "", -1)
	httpaddressbytes, err := i.config.Addresses.Gateway.MarshalJSON()
	if err != nil {
		return err
	}
	i.forwardHTTP = strings.Replace(string(httpaddressbytes), "\"", "", -1)
	i.i2pconfig, err = i2pgateconfig.ConfigAt(i.configPath)
	if err != nil {
		return err
	}
    idbytes := i.config.Identity.PeerID
    idstring := strings.Replace(string(idbytes), "\"", "", -1)
    i.id, err = peer.IDFromString(idstring)
	if err != nil {
		return err
	}
    err = i2pgateconfig.AddressRPC(i.forwardRPC, i.i2pconfig)
	if err != nil {
		return err
	}
	err = i2pgateconfig.AddressHTTP(i.forwardHTTP, i.i2pconfig)
	if err != nil {
		return err
	}
	i.i2pconfig, err = i.i2pconfig.Save(i.configPath)
	if err != nil {
		return err
	}
	go i.transportHTTP()
	go i.transportRPC()
	return nil
}

func (i *I2PGatePlugin) transportHTTP() error {
	GarlicTCPTransport, err := i2ptcp.NewGarlicTCPTransportFromOptions(
        i2ptcp.LocalPeerID(i.id),
		i2ptcp.SAMHost(i.i2pconfig.SAMHost),
		i2ptcp.SAMPort(i.i2pconfig.SAMPort),
		i2ptcp.SAMPass(""),
		i2ptcp.KeysPath(i.configPath+".i2pkeys"),
		i2ptcp.OnlyGarlic(i.i2pconfig.OnlyI2P),
		i2ptcp.GarlicOptions(i.i2pconfig.Print()),
	)
	if err != nil {
		return err
	}
	GarlicTCPConn, err := GarlicTCPTransport.ListenI2P()
	if err != nil {
		return err
	}
	err = i2pgateconfig.ListenerBase32(GarlicTCPConn.Base32(), i.i2pconfig)
	if err != nil {
		return err
	}
	err = i2pgateconfig.ListenerBase64(GarlicTCPConn.Base64(), i.i2pconfig)
	if err != nil {
		return err
	}
	i.i2pconfig, err = i.i2pconfig.Save(i.configPath)
	if err != nil {
		return err
	}
	GarlicTCPConn.ForwardI2P(i.i2pconfig.MaTargetHTTP())
    return nil
}

func (i *I2PGatePlugin) transportRPC() error {
	GarlicTCPTransport, err := i2ptcp.NewGarlicTCPTransportFromOptions(
        i2ptcp.LocalPeerID(i.id),
		i2ptcp.SAMHost(i.i2pconfig.SAMHost),
		i2ptcp.SAMPort(i.i2pconfig.SAMPort),
		i2ptcp.SAMPass(""),
		i2ptcp.KeysPath(i.configPath+".i2pkeys"),
		i2ptcp.OnlyGarlic(i.i2pconfig.OnlyI2P),
		i2ptcp.GarlicOptions(i.i2pconfig.Print()),
	)
	if err != nil {
		return err
	}
	GarlicTCPConn, err := GarlicTCPTransport.ListenI2P()
	if err != nil {
		return err
	}
	err = i2pgateconfig.ListenerBase32RPC(GarlicTCPConn.Base32(), i.i2pconfig)
	if err != nil {
		return err
	}
	err = i2pgateconfig.ListenerBase64RPC(GarlicTCPConn.Base64(), i.i2pconfig)
	if err != nil {
		return err
	}
	i.i2pconfig, err = i.i2pconfig.Save(i.configPath)
	if err != nil {
		return err
	}
	GarlicTCPConn.ForwardI2P(i.i2pconfig.MaTargetRPC())
    return nil
}

// I2PTypeName returns I2PType
func (*I2PGatePlugin) I2PTypeName() string {
	return I2PType
}
