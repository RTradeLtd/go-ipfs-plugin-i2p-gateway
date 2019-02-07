package manet

import (
	"net"
	"testing"

	ma "gx/ipfs/QmNTCey11oxhb1AxDnQBRHtdhap6Ctud872NjAYPYYXPuc/go-multiaddr"
)

func TestRegisterSpec(t *testing.T) {
	cm := NewCodecMap()
	myproto := &NetCodec{
		ProtocolName:     "test",
		NetAddrNetworks:  []string{"test", "iptest", "blahtest"},
		ConvertMultiaddr: func(a ma.Multiaddr) (net.Addr, error) { return nil, nil },
		ParseNetAddr:     func(a net.Addr) (ma.Multiaddr, error) { return nil, nil },
	}

	cm.RegisterNetCodec(myproto)

	_, ok := cm.addrParsers["test"]
	if !ok {
		t.Fatal("myproto not properly registered")
	}

	_, ok = cm.addrParsers["iptest"]
	if !ok {
		t.Fatal("myproto not properly registered")
	}

	_, ok = cm.addrParsers["blahtest"]
	if !ok {
		t.Fatal("myproto not properly registered")
	}

	_, ok = cm.maddrParsers["test"]
	if !ok {
		t.Fatal("myproto not properly registered")
	}

	_, ok = cm.maddrParsers["iptest"]
	if ok {
		t.Fatal("myproto not properly registered")
	}

	_, ok = cm.maddrParsers["blahtest"]
	if ok {
		t.Fatal("myproto not properly registered")
	}
}
