package i2pgate

import (
	fsrepo "github.com/ipsn/go-ipfs/repo/fsrepo"
	//peer "github.com/libp2p/go-libp2p-peer"
    "log"
	"github.com/rtradeltd/go-ipfs-plugin-i2p-gateway/config"
	"os"
	"testing"
    "github.com/eyedeekay/sam-forwarder"
)

var configPath = "./"

/*
func Test_Config(t *testing.T) {
	i := &I2PGatePlugin{}
	if err := i.test(); err != nil {
		t.Fatal(err)
	}
}

func (i *I2PGatePlugin) test() error {
	i.configPath = "./"
	err := os.Setenv("KEYS_PATH", i.configPath)
	if err != nil {
		return err
	}
	i.config, err = fsrepo.ConfigAt(i.configPath)
	if err != nil {
		return err
	}

	i.forwardHTTP = i.httpString()
    i.forwardRPC = i.rpcString()

	i.i2pconfig, err = i2pgateconfig.ConfigAt(i.configPath)
	if err != nil {
		return err
	}

	i.id, err = peer.IDFromString(i.idString())
	if err != nil {
		return err
	}
	err = i.configGateway()
	if err != nil {
		return err
	}
	i.i2pconfig, err = i.i2pconfig.Save(i.configPath)
	if err != nil {
		return err
	}
	go i.transportHTTP()
    go i.transportRPC()
	return nil
}
*/

// Test_config tries to create a config file
func Test_Config(t *testing.T) {

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
	i2pconfig, err := i2pgateconfig.ConfigAt(configPath)
	if err != nil {
		t.Fatal(err)
	}
	err = i2pgateconfig.AddressRPC(forwardRPC, i2pconfig)
	if err != nil {
		t.Fatal(err)
	}
	err = i2pgateconfig.AddressHTTP(forwardHTTP, i2pconfig)
	if err != nil {
		t.Fatal(err)
	}
	_, err = i2pconfig.Save(configPath)
	if err != nil {
		t.Fatal(err)
	}
	transportHTTP(i2pconfig)
	transportRPC(i2pconfig)
}

func transportHTTP(i2pconfig *i2pgateconfig.Config) error {
	GarlicForwarder, err := samforwarder.NewSAMForwarderFromOptions(
		samforwarder.SetSAMHost(i2pconfig.HostSAM()),
		samforwarder.SetSAMPort(i2pconfig.PortSAM()),
		samforwarder.SetHost(i2pconfig.HTTPHost()),
		samforwarder.SetPort(i2pconfig.HTTPPort()),
		samforwarder.SetType("server"),
		samforwarder.SetSaveFile(true),
		samforwarder.SetName("ipfs-gateway-http"),
		samforwarder.SetInLength(i2pconfig.InLength),
		samforwarder.SetOutLength(i2pconfig.OutLength),
		samforwarder.SetInVariance(i2pconfig.InVariance),
		samforwarder.SetOutVariance(i2pconfig.OutVariance),
		samforwarder.SetInQuantity(i2pconfig.InQuantity),
		samforwarder.SetOutQuantity(i2pconfig.OutQuantity),
		samforwarder.SetInBackups(i2pconfig.InBackupQuantity),
		samforwarder.SetOutBackups(i2pconfig.OutBackupQuantity),
		samforwarder.SetAllowZeroIn(i2pconfig.InAllowZeroHop),
		samforwarder.SetAllowZeroOut(i2pconfig.OutAllowZeroHop),
		samforwarder.SetCompress(i2pconfig.UseCompression),
		samforwarder.SetReduceIdle(i2pconfig.ReduceIdle),
		samforwarder.SetReduceIdleTimeMs(i2pconfig.ReduceIdleTime),
		samforwarder.SetReduceIdleQuantity(i2pconfig.ReduceIdleQuantity),
		samforwarder.SetAccessListType(i2pconfig.AccessListType),
		samforwarder.SetAccessList(i2pconfig.AccessList),
		samforwarder.SetEncrypt(i2pconfig.EncryptLeaseSet),
		samforwarder.SetLeaseSetKey(i2pconfig.EncryptedLeaseSetKey),
		samforwarder.SetLeaseSetPrivateKey(i2pconfig.EncryptedLeaseSetPrivateKey),
		samforwarder.SetLeaseSetPrivateSigningKey(i2pconfig.EncryptedLeaseSetPrivateSigningKey),
		samforwarder.SetMessageReliability(i2pconfig.MessageReliability),
	)
	if err != nil {
		return err
	}
	go GarlicForwarder.Serve()
	for len(GarlicForwarder.Base32()) < 51 {
		log.Println("Waiting for i2p destination to be generated(HTTP)")
	}
	err = i2pgateconfig.ListenerBase32(GarlicForwarder.Base32(), i2pconfig)
	if err != nil {
		return err
	}
	err = i2pgateconfig.ListenerBase64(GarlicForwarder.Base64(), i2pconfig)
	if err != nil {
		return err
	}
	i2pconfig, err = i2pconfig.Save(configPath)
	if err != nil {
		return err
	}
	return nil
}

func transportRPC(i2pconfig *i2pgateconfig.Config) error {
	GarlicForwarder, err := samforwarder.NewSAMForwarderFromOptions(
		samforwarder.SetSAMHost(i2pconfig.HostSAM()),
		samforwarder.SetSAMPort(i2pconfig.PortSAM()),
		samforwarder.SetHost(i2pconfig.RPCHost()),
		samforwarder.SetPort(i2pconfig.RPCPort()),
		samforwarder.SetType("server"),
		samforwarder.SetSaveFile(true),
		samforwarder.SetName("ipfs-gateway-rpc"),
		samforwarder.SetInLength(i2pconfig.InLength),
		samforwarder.SetOutLength(i2pconfig.OutLength),
		samforwarder.SetInVariance(i2pconfig.InVariance),
		samforwarder.SetOutVariance(i2pconfig.OutVariance),
		samforwarder.SetInQuantity(i2pconfig.InQuantity),
		samforwarder.SetOutQuantity(i2pconfig.OutQuantity),
		samforwarder.SetInBackups(i2pconfig.InBackupQuantity),
		samforwarder.SetOutBackups(i2pconfig.OutBackupQuantity),
		samforwarder.SetAllowZeroIn(i2pconfig.InAllowZeroHop),
		samforwarder.SetAllowZeroOut(i2pconfig.OutAllowZeroHop),
		samforwarder.SetCompress(i2pconfig.UseCompression),
		samforwarder.SetReduceIdle(i2pconfig.ReduceIdle),
		samforwarder.SetReduceIdleTimeMs(i2pconfig.ReduceIdleTime),
		samforwarder.SetReduceIdleQuantity(i2pconfig.ReduceIdleQuantity),
		samforwarder.SetAccessListType(i2pconfig.AccessListType),
		samforwarder.SetAccessList(i2pconfig.AccessList),
		samforwarder.SetEncrypt(i2pconfig.EncryptLeaseSet),
		samforwarder.SetLeaseSetKey(i2pconfig.EncryptedLeaseSetKey),
		samforwarder.SetLeaseSetPrivateKey(i2pconfig.EncryptedLeaseSetPrivateKey),
		samforwarder.SetLeaseSetPrivateSigningKey(i2pconfig.EncryptedLeaseSetPrivateSigningKey),
		samforwarder.SetMessageReliability(i2pconfig.MessageReliability),
	)
	if err != nil {
		return err
	}
	go GarlicForwarder.Serve()
	for len(GarlicForwarder.Base32()) < 51 {
		log.Println("Waiting for i2p destination to be generated(RPC)")
	}
	err = i2pgateconfig.ListenerBase32RPC(GarlicForwarder.Base32(), i2pconfig)
	if err != nil {
		return err
	}
	err = i2pgateconfig.ListenerBase64RPC(GarlicForwarder.Base64(), i2pconfig)
	if err != nil {
		return err
	}
	i2pconfig, err = i2pconfig.Save(configPath)
	if err != nil {
		return err
	}
	return nil
}
