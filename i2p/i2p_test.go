package i2pgate

import (
	"os"
	"testing"

	//"github.com/rtradeltd/go-ipfs-plugin-i2p-gateway/config"
)

var configPath = "./"

// Test_Network tries to create a config file
func Test_Network(t *testing.T) {

	err := os.Setenv("IPFS_PATH", configPath)

    var i *I2PGatePlugin
    err = i.Init()
    if err != nil {
		t.Fatal(err)
	}
	/*
    err = i.transportHTTP()
	if err != nil {
		t.Fatal(err)
	}
	err = i.transportRPC()
	if err != nil {
		t.Fatal(err)
	}
    */
}
