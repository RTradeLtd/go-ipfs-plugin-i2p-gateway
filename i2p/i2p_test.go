package i2pgate

import (
	"os"
	"testing"
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
    i.falseStart()

    err = i.Init()
	if err != nil {
		t.Fatal(err)
	}
}
