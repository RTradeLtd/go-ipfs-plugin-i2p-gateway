package i2pgate

import (
	"os"
	"testing"

	"github.com/rtradeltd/go-ipfs-plugin-i2p-gateway/config"
)

var configPath = "./"

// Test_Network tries to create a config file
func Test_Network(t *testing.T) {

	err := os.Setenv("IPFS_PATH", configPath)

	i, err := Setup()
	if err != nil {
		t.Fatal(err)
	}

	i2pconfig, err := i2pgateconfig.ConfigAt(i.configPath)
	if err != nil {
		t.Fatal(err)
	}
	err = i2pgateconfig.AddressRPC(i.forwardRPC, i2pconfig)
	if err != nil {
		t.Fatal(err)
	}
	err = i2pgateconfig.AddressHTTP(i.forwardHTTP, i2pconfig)
	if err != nil {
		t.Fatal(err)
	}
	_, err = i2pconfig.Save(configPath)
	if err != nil {
		t.Fatal(err)
	}
	err = i.transportHTTP()
	if err != nil {
		t.Fatal(err)
	}
	err = i.transportRPC()
	if err != nil {
		t.Fatal(err)
	}
}
