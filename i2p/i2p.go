package i2p

import (
	"fmt"
	"io"
	"time"

	ds "github.com/ipfs/go-datastore"
	delayed "github.com/ipfs/go-datastore/delayed"
	delay "github.com/ipfs/go-ipfs-delay"
	plugin "github.com/ipfs/go-ipfs/plugin"
	repo "github.com/ipfs/go-ipfs/repo"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
)

type I2PPlugin struct{}

// DatastoreType is this datastore's type name (used to identify the datastore
// in the datastore config).
var DatastoreType = "delaystore"

var _ plugin.PluginDatastore = (*I2PPlugin)(nil)

// Name returns the plugin's name, satisfying the plugin.Plugin interface.
func (*I2PPlugin) Name() string {
	return "ds-delaystore"
}

// Version returns the plugin's version, satisfying the plugin.Plugin interface.
func (*I2PPlugin) Version() string {
	return "0.1.0"
}

// Init initializes plugin, satisfying the plugin.Plugin interface. Put any
// initialization logic here.
func (*I2PPlugin) Init() error {
	return nil
}

// DatastoreTypeName returns the datastore's name. Every datastore
// implementation must have a unique name.
func (*I2PPlugin) DatastoreTypeName() string {
	return DatastoreType
}

type datastoreConfig struct {
	delay time.Duration
	inner fsrepo.DatastoreConfig
}

// DatastoreConfigParser returns a configuration parser for I2P configs.
func (*I2PPlugin) DatastoreConfigParser() fsrepo.ConfigFromMap {
	return func(params map[string]interface{}) (fsrepo.DatastoreConfig, error) {
		var delay time.Duration
		switch d := params["delay"].(type) {
		case string:
			var err error
			delay, err = time.ParseDuration(d)
			if err != nil {
				return nil, fmt.Errorf("delaystore: invalid delay: %s", err)
			}
		case float64:
			delay = time.Duration(d * float64(time.Second))
		case nil:
			return nil, fmt.Errorf("delaystore: no delay configured")
		}
		innerSpec, ok := params["inner"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("delaystore: no inner datastore specified")
		}
		inner, err := fsrepo.AnyDatastoreConfig(innerSpec)
		if err != nil {
			return nil, err
		}
		return &datastoreConfig{
			delay: delay,
			inner: inner,
		}, nil
	}
}

// DiskSpec returns this datastore's config.
func (c *datastoreConfig) DiskSpec() fsrepo.DiskSpec {
	return map[string]interface{}{
		// "type" is *mandatory*
		"type":  DatastoreType,
		"delay": c,
		"inner": c.inner.DiskSpec(),
	}
}

// Create creates or opens the datastore.
func (c *datastoreConfig) Create(path string) (repo.Datastore, error) {
	inner, err := c.inner.Create(path)
	if err != nil {
		return nil, err
	}
	// FIXME: We can return the delayed datastore directly once
	// https://github.com/ipfs/go-datastore/pull/108 is merged.
	return struct {
		ds.Batching
		io.Closer
	}{
		delayed.New(inner, delay.Fixed(c.delay)).(ds.Batching),
		inner,
	}, nil
}
