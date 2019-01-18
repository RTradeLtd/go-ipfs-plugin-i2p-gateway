package i2pgate

import (
	"github.com/rtradeltd/go-garlic-tcp-transport"
	"github.com/rtradeltd/go-garlic-tcp-transport/conn"
	"github.com/rtradeltd/go-ipfs-plugin-i2p-gateway/config"

	//TODO: Fix this. Get a better understanding of gx.
	//config "gx/ipfs/QmRd5T3VmYoX6jaNoZovFRQcwWHJqHgTVQTs1Qz92ELJ7C/go-ipfs-config"
	config "gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config"
	plugin "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/plugin"
	fsrepo "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/repo/fsrepo"
)

type I2PGatePlugin struct {
	*i2ptcp.GarlicTCPTransport
	*i2ptcpconn.GarlicTCPConn
	configPath    string
	config        *config.Config
	i2pconfigPath string
	i2pconfig     *i2pgateconfig.Config

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
	i.config, err = fsrepo.ConfigAt(i.configPath)
	if err != nil {
		return err
	}
	rpcaddressbytes, err := i.config.Addresses.API.MarshalJSON()
	if err != nil {
		return err
	}
	i.forwardRPC = string(rpcaddressbytes)
	httpaddressbytes, err := i.config.Addresses.Gateway.MarshalJSON()
	if err != nil {
		return err
	}
	i.forwardHTTP = string(httpaddressbytes)
	i.i2pconfig, err = i2pgateconfig.ConfigAt(i.configPath)
	if err != nil {
		return err
	}
	err = i2pgateconfig.AddressRPC(i.forwardRPC, i.forwardRPC)
	if err != nil {
		return err
	}
	err = i2pgateconfig.AddressHTTP(i.forwardHTTP, i.forwardHTTP)
	if err != nil {
		return err
	}
	i.i2pconfig, err = i2pgateconfig.WriteConfig(i.configPath, i.i2pconfig)
	if err != nil {
		return err
	}
	i.GarlicTCPTransport, err = i2ptcp.NewGarlicTCPTransportFromOptions(
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
	i.GarlicTCPConn, err = i.GarlicTCPTransport.ListenI2P()
	if err != nil {
		return err
	}
	err = i2pgateconfig.ListenerBase32(i.GarlicTCPConn.Base32())
	if err != nil {
		return err
	}
	err = i2pgateconfig.ListenerBase64(i.GarlicTCPConn.MaBase64())
	if err != nil {
		return err
	}
	return nil
}

// I2PTypeName returns I2PType
func (*I2PGatePlugin) I2PTypeName() string {
	return I2PType
}
