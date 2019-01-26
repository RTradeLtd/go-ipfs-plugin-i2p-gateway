package i2pgate

import (
	"log"
	"os"
	"strings"

	//"github.com/rtradeltd/go-ipfs-plugin-i2p-gateway/config"
    "../config"
	//TODO: Get a better understanding of gx.
	config "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-config"
	plugin "github.com/ipsn/go-ipfs/plugin"
	fsrepo "github.com/ipsn/go-ipfs/repo/fsrepo"
	peer "github.com/libp2p/go-libp2p-peer"
)

type I2PGatePlugin struct {
	configPath    string
	config        *config.Config
	i2pconfigPath string
	i2pconfig     *i2pgateconfig.Config
	id            peer.ID

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
	i.forwardRPC = i.rpcString()
	i.forwardHTTP = i.httpString()

	i.i2pconfig, err = i2pgateconfig.ConfigAt(i.configPath)
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
	go i.transportHTTP()
	go i.transportRPC()
	return nil
}

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
	i.id, err = peer.IDFromString(i.idString())
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
	return unquote(string(rpcaddressbytes))
}

func (i *I2PGatePlugin) httpString() string {
	httpaddressbytes, err := i.config.Addresses.Gateway.MarshalJSON()
	if err != nil {
		panic("could not read HTTP address, aborting")
	}
	return unquote(string(httpaddressbytes))
}

func (i *I2PGatePlugin) idString() string {
	idbytes := i.config.Identity.PeerID
	return unquote(string(idbytes))
}

// I2PTypeName returns I2PType
func (*I2PGatePlugin) I2PTypeName() string {
	return I2PType
}

func unquote(s string) string{
    return strings.Replace(s, "\"", "", -1)
}
