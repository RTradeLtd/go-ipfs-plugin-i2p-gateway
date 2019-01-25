package i2pgateconfig

import (
	"os"
	"testing"

	"github.com/rtradeltd/go-garlic-tcp-transport"

	fsrepo "github.com/ipsn/go-ipfs/repo/fsrepo"
)

// Test_config tries to create a config file
func Test_Config(t *testing.T) {
	configPath := "./"
	err := os.Setenv("KEYS_PATH", configPath)
	if err != nil {
		t.Fatal("")
	}
	config, err := fsrepo.ConfigAt(configPath)
	if err != nil {
		t.Fatal(err)
	}
	rpcaddressbytes, err := config.Addresses.API.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	forwardRPC := string(rpcaddressbytes)
	httpaddressbytes, err := config.Addresses.Gateway.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	forwardHTTP := string(httpaddressbytes)
	i2pconfig, err := ConfigAt(configPath)
	if err != nil {
		t.Fatal(err)
	}
	err = AddressRPC(forwardRPC, i2pconfig)
	if err != nil {
		t.Fatal(err)
	}
	err = AddressHTTP(forwardHTTP, i2pconfig)
	if err != nil {
		t.Fatal(err)
	}
	_, err = i2pconfig.Save(configPath)
	if err != nil {
		t.Fatal(err)
	}
	GarlicTCPTransport, err := i2ptcp.NewGarlicTCPTransportFromOptions(
		i2ptcp.SAMHost(i2pconfig.SAMHost),
		i2ptcp.SAMPort(i2pconfig.SAMPort),
		i2ptcp.SAMPass(""),
		i2ptcp.KeysPath(configPath+".i2pkeys"),
		i2ptcp.OnlyGarlic(i2pconfig.OnlyI2P),
		i2ptcp.GarlicOptions(i2pconfig.Print()),
	)
	if err != nil {
		t.Fatal(err)
	}
	GarlicTCPConn, err := GarlicTCPTransport.ListenI2P()
	if err != nil {
		t.Fatal(err)
	}
	err = ListenerBase32(GarlicTCPConn.Base32(), i2pconfig)
	if err != nil {
		t.Fatal(err)
	}
	err = ListenerBase64(GarlicTCPConn.Base64(), i2pconfig)
	if err != nil {
		t.Fatal(err)
	}
	_, err = i2pconfig.Save(configPath)
	if err != nil {
		t.Fatal(err)
	}
}