package main

import (
	plugin "gx/ipfs/QmXZuSpcGSesFXDWwZnESp2YEcYNcR4em9P86XsZtcuzWR/iptb-plugins/local"
	testbedi "gx/ipfs/QmckeQ2zrYLAXoSHYTGn5BDdb22BqbUoHEHm8KZ9YWRxd1/iptb/testbed/interfaces"
)

var PluginName string
var NewNode testbedi.NewNodeFunc
var GetAttrList testbedi.GetAttrListFunc
var GetAttrDesc testbedi.GetAttrDescFunc

func init() {
	PluginName = plugin.PluginName
	NewNode = plugin.NewNode
	GetAttrList = plugin.GetAttrList
	GetAttrDesc = plugin.GetAttrDesc
}
