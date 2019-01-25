package i2p

import (
	"time"

	plugin "github.com/ipfs/go-ipfs/plugin"
)

type I2PPlugin struct{}

// I2PType is this datastore's type name (used to identify the datastore
// in the datastore config).
var I2PType = "delaystore"

var _ plugin.Plugin = (*I2PPlugin)(nil)

// Name returns the plugin's name, satisfying the plugin.Plugin interface.
func (*I2PPlugin) Name() string {
	return "ds-delaystore"
}

// Version returns the plugin's version, satisfying the plugin.Plugin interface.
func (*I2PPlugin) Version() string {
	return "0.0.1"
}

// Init initializes plugin, satisfying the plugin.Plugin interface. Put any
// initialization logic here.
func (*I2PPlugin) Init() error {
	return nil
}

// I2PTypeName returns the datastore's name. Every datastore
// implementation must have a unique name.
func (*I2PPlugin) I2PTypeName() string {
	return I2PType
}

type I2PConfig struct {
	//delay time.Duration
	//inner fsrepo.DatastoreConfig
}
