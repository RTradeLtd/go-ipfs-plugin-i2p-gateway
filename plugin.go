package main

import (
	i2p "github.com/rtradeltd/go-ipfs-plugin-i2p-gateway/i2p"
	plugin "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/plugin"
)

// Plugins is an exported list of plugins that will be loaded by go-ipfs.
var Plugins = []plugin.Plugin{
	&i2p.I2PGatePlugin{},
}
