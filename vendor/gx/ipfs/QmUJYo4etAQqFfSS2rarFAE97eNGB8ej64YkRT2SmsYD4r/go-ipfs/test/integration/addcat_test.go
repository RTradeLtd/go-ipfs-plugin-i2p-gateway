package integrationtest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
	"time"

	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core/coreapi"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core/coreapi/interface"
	coreunix "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core/coreunix"
	mock "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core/mock"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/thirdparty/unit"

	random "gx/ipfs/QmSJ9n2s9NUoA9D849W5jj5SJ94nMcZpj1jCgQJieiNqSt/go-random"
	pstore "gx/ipfs/QmTTJcDL3gsnGDALjh2fDGg1onGRUdVgNL2hU2WEZcVrMX/go-libp2p-peerstore"
	mocknet "gx/ipfs/QmUDTcnDp2WssbmiDLC6aYurUeyt7QeRakHUQMxA2mZ5iB/go-libp2p/p2p/net/mock"
	logging "gx/ipfs/QmZChCsSt8DctjceaL56Eibc29CVQq4dGKRXC5JRZ6Ppae/go-log"
	testutil "gx/ipfs/Qma6ESRQTf1ZLPgzpCwDTqQJefPnU6uLvMjP18vK8EWp8L/go-testutil"
)

var log = logging.Logger("epictest")

const kSeed = 1

func Test1KBInstantaneous(t *testing.T) {
	conf := testutil.LatencyConfig{
		NetworkLatency:    0,
		RoutingLatency:    0,
		BlockstoreLatency: 0,
	}

	if err := DirectAddCat(RandomBytes(1*unit.KB), conf); err != nil {
		t.Fatal(err)
	}
}

func TestDegenerateSlowBlockstore(t *testing.T) {
	SkipUnlessEpic(t)
	conf := testutil.LatencyConfig{BlockstoreLatency: 50 * time.Millisecond}
	if err := AddCatPowers(conf, 128); err != nil {
		t.Fatal(err)
	}
}

func TestDegenerateSlowNetwork(t *testing.T) {
	SkipUnlessEpic(t)
	conf := testutil.LatencyConfig{NetworkLatency: 400 * time.Millisecond}
	if err := AddCatPowers(conf, 128); err != nil {
		t.Fatal(err)
	}
}

func TestDegenerateSlowRouting(t *testing.T) {
	SkipUnlessEpic(t)
	conf := testutil.LatencyConfig{RoutingLatency: 400 * time.Millisecond}
	if err := AddCatPowers(conf, 128); err != nil {
		t.Fatal(err)
	}
}

func Test100MBMacbookCoastToCoast(t *testing.T) {
	SkipUnlessEpic(t)
	conf := testutil.LatencyConfig{}.NetworkNYtoSF().BlockstoreSlowSSD2014().RoutingSlow()
	if err := DirectAddCat(RandomBytes(100*1024*1024), conf); err != nil {
		t.Fatal(err)
	}
}

func AddCatPowers(conf testutil.LatencyConfig, megabytesMax int64) error {
	var i int64
	for i = 1; i < megabytesMax; i = i * 2 {
		fmt.Printf("%d MB\n", i)
		if err := DirectAddCat(RandomBytes(i*1024*1024), conf); err != nil {
			return err
		}
	}
	return nil
}

func RandomBytes(n int64) []byte {
	data := new(bytes.Buffer)
	random.WritePseudoRandomBytes(n, data, kSeed)
	return data.Bytes()
}

func DirectAddCat(data []byte, conf testutil.LatencyConfig) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create network
	mn := mocknet.New(ctx)
	mn.SetLinkDefaults(mocknet.LinkOptions{
		Latency: conf.NetworkLatency,
		// TODO add to conf. This is tricky because we want 0 values to be functional.
		Bandwidth: math.MaxInt32,
	})

	adder, err := core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Host:   mock.MockHostOption(mn),
	})
	if err != nil {
		return err
	}
	defer adder.Close()

	catter, err := core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Host:   mock.MockHostOption(mn),
	})
	if err != nil {
		return err
	}
	defer catter.Close()

	catterApi := coreapi.NewCoreAPI(catter)

	err = mn.LinkAll()
	if err != nil {
		return err
	}

	bs1 := []pstore.PeerInfo{adder.Peerstore.PeerInfo(adder.Identity)}
	bs2 := []pstore.PeerInfo{catter.Peerstore.PeerInfo(catter.Identity)}

	if err := catter.Bootstrap(core.BootstrapConfigWithPeers(bs1)); err != nil {
		return err
	}
	if err := adder.Bootstrap(core.BootstrapConfigWithPeers(bs2)); err != nil {
		return err
	}

	added, err := coreunix.Add(adder, bytes.NewReader(data))
	if err != nil {
		return err
	}

	ap, err := iface.ParsePath(added)
	if err != nil {
		return err
	}

	readerCatted, err := catterApi.Unixfs().Get(ctx, ap)
	if err != nil {
		return err
	}

	// verify
	bufout := new(bytes.Buffer)
	io.Copy(bufout, readerCatted)
	if 0 != bytes.Compare(bufout.Bytes(), data) {
		return errors.New("catted data does not match added data")
	}

	cancel()
	return nil
}

func SkipUnlessEpic(t *testing.T) {
	if os.Getenv("IPFS_EPIC_TEST") == "" {
		t.SkipNow()
	}
}
