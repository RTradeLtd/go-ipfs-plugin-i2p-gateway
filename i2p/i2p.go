package i2pgate

import (
	"log"
	"os"

	config "gx/ipfs/QmTbcMKv6GU3fxhnNcbzYChdox9Fdd7VpucM3PQ7UWjX3D/go-ipfs-config"

	"github.com/RTradeLtd/go-ipfs-plugin-i2p-gateway/config"
	plugin "github.com/ipfs/go-ipfs/plugin"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	coreiface "gx/ipfs/QmNmqKNivNTN11HrKWJYt29n6Z2fuzkeDheQV62dbxNuLb/interface-go-ipfs-core"
)

// I2PGatePlugin is a structure containing information which is used for
// setting up an i2p tunnel that connects an IPFS gateway to a tunnel over i2p.
type I2PGatePlugin struct {
	configPath    string
	config        *config.Config
	i2pconfigPath string
	i2pconfig     *i2pgateconfig.Config

	forwardHTTP  string
	forwardRPC   string
	forwardSwarm string
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
	/*i := Setup()
	    if err != nil {
			return nil, err
		}*/
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
	i.forwardSwarm = i.swarmString()
	log.Println("Prepared to forward:", i.forwardRPC, i.forwardHTTP)
	i.i2pconfig, err = i2pgateconfig.ConfigAt(i.configPath)
	if err != nil {
		return nil, err
	}
	err = i2pgateconfig.AddressRPC(i.forwardRPC, i.i2pconfig)
	if err != nil {
		return nil, err
	}
	err = i2pgateconfig.AddressHTTP(i.forwardHTTP, i.i2pconfig)
	if err != nil {
		return nil, err
	}
	err = i2pgateconfig.AddressSwarm(i.forwardSwarm, i.i2pconfig)
	if err != nil {
		return nil, err
	}

	i.i2pconfig, err = i.i2pconfig.Save(i.configPath)
	if err != nil {
		return nil, err
	}

	/*i.i2pconfig, err = i.i2pconfig.Save(i.configPath)
	if err != nil {
		return nil, err
	}*/
	return &i, nil
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

func (i *I2PGatePlugin) swarmString() string {
	swarmaddressbytes := ""
	for _, v := range i.config.Addresses.Swarm {
		swarmaddressbytes += v
	}
	return i2pgateconfig.Unquote(string(swarmaddressbytes))
}

func (i *I2PGatePlugin) idString() string {
	idbytes := i.config.Identity.PeerID
	return i2pgateconfig.Unquote(string(idbytes))
}

// I2PTypeName returns I2PType
func (*I2PGatePlugin) I2PTypeName() string {
	return I2PType
}

func (i *I2PGatePlugin) falseStart() error {
	i2p, err := Setup()
	if err != nil {
		return err
	}

	i2p.transportHTTP()
	i2p.transportSwarm()
	i2p.transportRPC()

	return nil
}

// Start starts the tunnels and also satisfies the Daemon plugin interface
func (i *I2PGatePlugin) Start(coreiface.CoreAPI) error {
	i2p, err := Setup()
	if err != nil {
		return err
	}

	go i2p.transportHTTP()
	go i2p.transportSwarm()
	// only create tunnel if unsafe rpc access is permitted
	if os.Getenv("UNSAFE_RPC_ACCESS") == "yes" {
		go i2p.transportRPC()
	}

	return nil
}

// Close satisfies the Daemon plugin interface
func (*I2PGatePlugin) Close() error {
	return nil
}
