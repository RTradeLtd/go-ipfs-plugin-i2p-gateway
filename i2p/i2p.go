package i2pgate

import (
	"log"
	"os"

	config "gx/ipfs/QmTbcMKv6GU3fxhnNcbzYChdox9Fdd7VpucM3PQ7UWjX3D/go-ipfs-config"

	"github.com/RTradeLtd/go-ipfs-plugin-i2p-gateway/config"
	coreiface "github.com/ipfs/go-ipfs/core/coreapi/interface"
	plugin "github.com/ipfs/go-ipfs/plugin"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
)

// I2PGatePlugin is a structure containing information which is used for
// setting up an i2p tunnel that connects an IPFS gateway to a tunnel over i2p.
type I2PGatePlugin struct {
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

var _ plugin.PluginDaemon = (*I2PGatePlugin)(nil)

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
    i, err = Setup()
	if err != nil {
		return err
	}

	err = i.configGateway()
	if err != nil {
		return err
	}

	i.i2pconfig, err = i.i2pconfig.Save(i.configPath)
	if err != nil {
		return err
	}

	return nil
}

// Setup creates an I2PGatePlugin and config file, but it doesn't start
// any tunnels.
func Setup() (*I2PGatePlugin, error) {
	var err error
	var i I2PGatePlugin
	i.configPath, err = fsrepo.BestKnownPath()
	if err != nil {
		return nil, err
	}
	err = os.Setenv("KEYS_PATH", i.configPath)
	if err != nil {
		return nil, err
	}
	i.config, err = fsrepo.ConfigAt(i.configPath)
	if err != nil {
		return nil, err
	}
	i.forwardRPC = i.rpcString()
	i.forwardHTTP = i.httpString()
	log.Println("Prepared to forward:", i.forwardRPC, i.forwardHTTP)
	i.i2pconfig, err = i2pgateconfig.ConfigAt(i.configPath)
	return &i, nil
}

func (i I2PGatePlugin) configGateway() error {
	err := i2pgateconfig.AddressRPC(i.forwardRPC, i.i2pconfig)
	if err != nil {
		return err
	}
	err = i2pgateconfig.AddressHTTP(i.forwardHTTP, i.i2pconfig)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	log.Println(i.idString())
	i.i2pconfig, err = i.i2pconfig.Save(i.configPath)
	if err != nil {
		return err
	}
	return nil
}

func (i *I2PGatePlugin) rpcString() string {
	rpcaddressbytes, err := i.config.Addresses.API.MarshalJSON()
	if err != nil {
		panic("could not read RPC address, aborting")
	}
	return i2pgateconfig.Unquote(string(rpcaddressbytes))
}

func (i *I2PGatePlugin) httpString() string {
	httpaddressbytes, err := i.config.Addresses.Gateway.MarshalJSON()
	if err != nil {
		panic("could not read HTTP address, aborting")
	}
	return i2pgateconfig.Unquote(string(httpaddressbytes))
}

func (i *I2PGatePlugin) idString() string {
	idbytes := i.config.Identity.PeerID
	return i2pgateconfig.Unquote(string(idbytes))
}

// I2PTypeName returns I2PType
func (*I2PGatePlugin) I2PTypeName() string {
	return I2PType
}

// Start starts the tunnels and also satisfies the Daemon plugin interface
func (i *I2PGatePlugin) Start(coreiface.CoreAPI) error {
	i2p, err := Setup()
	if err != nil {
		return err
	}

	go i2p.transportHTTP()
	go i2p.transportRPC()

	return nil
}

// Close satisfies the Daemon plugin interface
func (*I2PGatePlugin) Close() error {
	return nil
}
