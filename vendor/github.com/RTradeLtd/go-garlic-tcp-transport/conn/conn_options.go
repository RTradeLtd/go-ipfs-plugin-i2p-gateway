package i2ptcpconn

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	peer "github.com/libp2p/go-libp2p-peer"
	tpt "github.com/libp2p/go-libp2p-transport"
)

// Option is a functional argument to the connection constructor
type Option func(*GarlicTCPConn) error

//Transport sets the parent transport of the connection.
func Transport(t tpt.Transport) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		c.parentTransport = t
		return nil
	}
}

//SAMHost sets the host of the SAM Bridge to use
func SAMHost(s string) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		st := ""
		if !strings.Contains("/ip4/", s) || !strings.Contains("/ip6/", s) {
			ip := net.ParseIP(s)
			if ip.To4() != nil {
				st = "/ip4/" + s + "/"
			}
			if ip.To16() != nil {
				st = "/ip6/" + s + "/"
			}
		}
		c.hostSAM = st
		return nil
	}
}

//SAMPort sets the port of the SAM bridge to use
func SAMPort(s string) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		st := strings.TrimPrefix(s, "/tcp/")
		rt := strings.TrimSuffix(st, "/")
		val, err := strconv.Atoi(rt)
		if err != nil {
			return fmt.Errorf("Connection Construction error: %s", err)
		}
		if val > 0 && val < 65536 {
			c.portSAM = "/tcp/" + rt + "/"
			return nil
		}
		return fmt.Errorf("port is %s invalid", s)
	}
}

//SAMPass sets the password to use when authenticating to the SAM bridge. It's
//ignored for now, and will return an error if it recieves a non-empty string.
func SAMPass(s string) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		if s != "" {
			return fmt.Errorf("SAMPass is unused for now, pass no argument(or empty string). Failing closed")
		}
		return nil
	}
}

//KeysPath sets the path to the keys, if no keys are present, they will be generated.
func KeysPath(s string) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		c.keysPath = s
		return nil
	}
}

//OnlyGarlic indicates that this connection will only be used to serve anonymous
//connections. It does nothing but indicate that for now.
func OnlyGarlic(b bool) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		c.onlyGarlic = b
		return nil
	}
}

// GarlicOptions is a slice of string-formatted options to pass to the SAM API.
func GarlicOptions(s []string) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		for _, v := range s {
			c.garlicOptions = append(c.garlicOptions, v)
		}
		return nil
	}
}

func LocalPeerID(p peer.ID) func(*GarlicTCPConn) error {
	return func(c *GarlicTCPConn) error {
		c.id = p
		return nil
	}
}
