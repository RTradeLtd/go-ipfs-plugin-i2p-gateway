package main

import (
	plugin "github.com/ipsn/go-ipfs/plugin"
	i2p "github.com/rtradeltd/go-ipfs-plugin-i2p-gateway/i2p"
)

// Plugins is an exported list of plugins that will be loaded by go-ipfs.
var Plugins = []plugin.Plugin{
	&i2p.I2PGatePlugin{},
}
