package i2ptcpcodec

import (
	"github.com/RTradeLtd/sam3"
	"net"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

// FromMultiaddrToNetAddr wraps around FromMultiaddrToI2PNetAddr to work with manet.NetCodec
func FromMultiaddrToNetAddr(from ma.Multiaddr) (net.Addr, error) {
	return FromMultiaddrToI2PNetAddr(from)
}

// FromMultiaddrToI2PNetAddr converts a ma.Multiaddr to a sam3.I2PAddr
func FromMultiaddrToI2PNetAddr(from ma.Multiaddr) (sam3.I2PAddr, error) {
	return sam3.NewI2PAddrFromString(from.String())
}

// FromNetAddrToMultiaddr wraps around FromI2PNetAddrToMultiaddr to work with manet.NetCodec
func FromNetAddrToMultiaddr(from net.Addr) (ma.Multiaddr, error) {
	return FromI2PNetAddrToMultiaddr(from.(sam3.I2PAddr))
}

// FromI2PNetAddrToMultiaddr converts a sam3.I2PAddr to a ma.Multiaddr
func FromI2PNetAddrToMultiaddr(from sam3.I2PAddr) (ma.Multiaddr, error) {
	return ma.NewMultiaddr("/garlic64/" + from.Base64())
}

func NewGarlicTCPNetCodec() manet.NetCodec {

	var fromNetAddr manet.FromNetAddrFunc
	fromNetAddr = FromNetAddrToMultiaddr

	var toMultiAddr manet.ToNetAddrFunc
	toMultiAddr = FromMultiaddrToNetAddr

	return manet.NetCodec{
		//NetAddrNetworks: ,
		ProtocolName: "garlic64",
		// ParseNetAddr parses a net.Addr belonging to this type into a multiaddr
		ParseNetAddr: fromNetAddr,
		// ConvertMultiaddr converts a multiaddr of this type back into a net.Addr
		ConvertMultiaddr: toMultiAddr,
		Protocol:         ma.ProtocolWithName("garlic64"),
	}
}
