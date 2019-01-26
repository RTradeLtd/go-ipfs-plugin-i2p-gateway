// +build !samforwarder

package i2pgate

import (
	"github.com/rtradeltd/go-garlic-tcp-transport"
	"github.com/rtradeltd/go-ipfs-plugin-i2p-gateway/config"
	//TODO: Get a better understanding of gx.
)

func (i *I2PGatePlugin) transportHTTP() error {
	GarlicTCPTransport, err := i2ptcp.NewGarlicTCPTransportFromOptions(
		i2ptcp.LocalPeerID(i.id),
		i2ptcp.SAMHost(i.i2pconfig.SAMHost),
		i2ptcp.SAMPort(i.i2pconfig.SAMPort),
		i2ptcp.SAMPass(""),
		i2ptcp.KeysPath(i.configPath+".i2pkeys"),
		i2ptcp.OnlyGarlic(i.i2pconfig.OnlyI2P),
		i2ptcp.GarlicOptions(i.i2pconfig.Print()),
	)
	if err != nil {
		return err
	}
	GarlicTCPConn, err := GarlicTCPTransport.ListenI2P()
	if err != nil {
		return err
	}
	err = i2pgateconfig.ListenerBase32(GarlicTCPConn.Base32(), i.i2pconfig)
	if err != nil {
		return err
	}
	err = i2pgateconfig.ListenerBase64(GarlicTCPConn.Base64(), i.i2pconfig)
	if err != nil {
		return err
	}
	i.i2pconfig, err = i.i2pconfig.Save(i.configPath)
	if err != nil {
		return err
	}
	GarlicTCPConn.ForwardI2P(i.i2pconfig.MaTargetHTTP())
	return nil
}

func (i *I2PGatePlugin) transportRPC() error {
	GarlicTCPTransport, err := i2ptcp.NewGarlicTCPTransportFromOptions(
		i2ptcp.LocalPeerID(i.id),
		i2ptcp.SAMHost(i.i2pconfig.SAMHost),
		i2ptcp.SAMPort(i.i2pconfig.SAMPort),
		i2ptcp.SAMPass(""),
		i2ptcp.KeysPath(i.configPath+".i2pkeys"),
		i2ptcp.OnlyGarlic(i.i2pconfig.OnlyI2P),
		i2ptcp.GarlicOptions(i.i2pconfig.Print()),
	)
	if err != nil {
		return err
	}
	GarlicTCPConn, err := GarlicTCPTransport.ListenI2P()
	if err != nil {
		return err
	}
	err = i2pgateconfig.ListenerBase32RPC(GarlicTCPConn.Base32(), i.i2pconfig)
	if err != nil {
		return err
	}
	err = i2pgateconfig.ListenerBase64RPC(GarlicTCPConn.Base64(), i.i2pconfig)
	if err != nil {
		return err
	}
	i.i2pconfig, err = i.i2pconfig.Save(i.configPath)
	if err != nil {
		return err
	}
	GarlicTCPConn.ForwardI2P(i.i2pconfig.MaTargetRPC())
	return nil
}
