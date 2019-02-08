package tcp

import (
	"testing"

	ma "gx/ipfs/QmNTCey11oxhb1AxDnQBRHtdhap6Ctud872NjAYPYYXPuc/go-multiaddr"
	tptu "gx/ipfs/QmQwvsMzMDTW2K8ySZYgnTVCkzQXVDxmGB5upvVFwdumJV/go-libp2p-transport-upgrader"
	utils "gx/ipfs/QmUDtgnEr7FFrtK2LQM2dFzTNWghnrApBDcU3iHEJz8eQS/go-libp2p-transport/test"
	insecure "gx/ipfs/QmWmeXRTSyWvjPQZgXjXTj2aP74tMSgJwWi1SAvHsvBJVj/go-conn-security/insecure"
	mplex "gx/ipfs/QmZsejKNkeFSQe5TcmYXJ8iq6qPL1FpsP4eAA8j7RfE7xg/go-smux-multiplex"
)

func TestTcpTransport(t *testing.T) {
	for i := 0; i < 2; i++ {
		ta := NewTCPTransport(&tptu.Upgrader{
			Secure: insecure.New("peerA"),
			Muxer:  new(mplex.Transport),
		})
		tb := NewTCPTransport(&tptu.Upgrader{
			Secure: insecure.New("peerB"),
			Muxer:  new(mplex.Transport),
		})

		zero := "/ip4/127.0.0.1/tcp/0"
		utils.SubtestTransport(t, ta, tb, zero, "peerA")

		envReuseportVal = false
	}
	envReuseportVal = true
}

func TestTcpTransportCantListenUtp(t *testing.T) {
	for i := 0; i < 2; i++ {
		utpa, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/0/utp")
		if err != nil {
			t.Fatal(err)
		}

		tpt := NewTCPTransport(&tptu.Upgrader{
			Secure: insecure.New("peerB"),
			Muxer:  new(mplex.Transport),
		})

		_, err = tpt.Listen(utpa)
		if err == nil {
			t.Fatal("shouldnt be able to listen on utp addr with tcp transport")
		}

		envReuseportVal = false
	}
	envReuseportVal = true
}
