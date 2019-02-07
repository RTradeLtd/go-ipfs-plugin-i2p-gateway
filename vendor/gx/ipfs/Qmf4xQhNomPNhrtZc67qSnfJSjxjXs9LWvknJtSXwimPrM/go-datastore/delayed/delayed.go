// Package delayed wraps a datastore allowing to artificially
// delay all operations.
package delayed

import (
	delay "gx/ipfs/QmUe1WCHkQaz4UeNKiHDUBV2T6i9prc3DniqyHPXyfGaUq/go-ipfs-delay"
	ds "gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore"
	dsq "gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore/query"
)

// New returns a new delayed datastore.
func New(ds ds.Datastore, delay delay.D) ds.Datastore {
	return &delayed{ds: ds, delay: delay}
}

type delayed struct {
	ds    ds.Datastore
	delay delay.D
}

func (dds *delayed) Put(key ds.Key, value []byte) (err error) {
	dds.delay.Wait()
	return dds.ds.Put(key, value)
}

func (dds *delayed) Get(key ds.Key) (value []byte, err error) {
	dds.delay.Wait()
	return dds.ds.Get(key)
}

func (dds *delayed) Has(key ds.Key) (exists bool, err error) {
	dds.delay.Wait()
	return dds.ds.Has(key)
}

func (dds *delayed) GetSize(key ds.Key) (size int, err error) {
	dds.delay.Wait()
	return dds.ds.GetSize(key)
}

func (dds *delayed) Delete(key ds.Key) (err error) {
	dds.delay.Wait()
	return dds.ds.Delete(key)
}

func (dds *delayed) Query(q dsq.Query) (dsq.Results, error) {
	dds.delay.Wait()
	return dds.ds.Query(q)
}

func (dds *delayed) Batch() (ds.Batch, error) {
	return ds.NewBasicBatch(dds), nil
}

func (dds *delayed) DiskUsage() (uint64, error) {
	dds.delay.Wait()
	return ds.DiskUsage(dds.ds)
}
