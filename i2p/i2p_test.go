package i2pgate

import (
	"log"
	"os"
	"testing"

	"github.com/eyedeekay/sam-forwarder"
	"github.com/rtradeltd/go-ipfs-plugin-i2p-gateway/config"
)

var configPath = "./"

// Test_Network tries to create a config file
func Test_Network(t *testing.T) {

	err := os.Setenv("IPFS_PATH", configPath)

	i, err := Setup(&I2PGatePlugin{})
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
	transportHTTP(i2pconfig)
	transportRPC(i2pconfig)
}

func transportHTTP(i2pconfig *i2pgateconfig.Config) error {
	host, err := i2pconfig.HTTPHost()
	if err != nil {
		return err
	}
	port, err := i2pconfig.HTTPPort()
	if err != nil {
		return err
	}
	GarlicForwarder, err := samforwarder.NewSAMForwarderFromOptions(
		samforwarder.SetSAMHost(i2pconfig.HostSAM()),
		samforwarder.SetSAMPort(i2pconfig.PortSAM()),
		samforwarder.SetHost(host),
		samforwarder.SetPort(port),
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
	log.Println("waiting for i2p forwarding")
	for true {
		if len(GarlicForwarder.Base32()) > 51 {
			log.Println(GarlicForwarder.Base32())
		} else {
			log.Println("waiting for i2p forwarding")
		}

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
	host, err := i2pconfig.RPCHost()
	if err != nil {
		return err
	}
	port, err := i2pconfig.RPCPort()
	if err != nil {
		return err
	}
	GarlicForwarder, err := samforwarder.NewSAMForwarderFromOptions(
		samforwarder.SetSAMHost(i2pconfig.HostSAM()),
		samforwarder.SetSAMPort(i2pconfig.PortSAM()),
		samforwarder.SetHost(host),
		samforwarder.SetPort(port),
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
	log.Println("waiting for i2p forwarding")
	go GarlicForwarder.Serve()
	log.Println(GarlicForwarder.Base32())
	log.Println("len", len(GarlicForwarder.Base32()))
	for true {
		if len(GarlicForwarder.Base32()) > 51 {
			log.Println(GarlicForwarder.Base32())
		} else {
			log.Println("waiting for i2p forwarding")
		}

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
