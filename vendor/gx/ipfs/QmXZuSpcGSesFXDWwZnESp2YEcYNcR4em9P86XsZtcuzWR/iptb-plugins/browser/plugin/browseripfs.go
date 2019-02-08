package main

import (
	plugin "gx/ipfs/QmXZuSpcGSesFXDWwZnESp2YEcYNcR4em9P86XsZtcuzWR/iptb-plugins/browser"
	testbedi "gx/ipfs/QmckeQ2zrYLAXoSHYTGn5BDdb22BqbUoHEHm8KZ9YWRxd1/iptb/testbed/interfaces"
)

var PluginName string
var NewNode testbedi.NewNodeFunc

func init() {
	PluginName = plugin.PluginName
	NewNode = plugin.NewNode
}
