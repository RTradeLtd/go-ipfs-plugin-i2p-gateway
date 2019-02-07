GOCC ?= go
IPFSVERSION=v0.4.18

.PHONY: install build gx

build: go-ipfs-plugin-i2p-gateway.so

clean:
	rm -f go-ipfs-plugin-i2p-gateway.so
	find . -name '*.i2pkeys' -exec rm -vf {} \;
	find . -name '*i2pconfig' -exec rm -vf {} \;

install: build
	mkdir -p $(IPFS_PATH)/plugins
	install -Dm700 go-ipfs-plugin-i2p-gateway.so "$(IPFS_PATH)/plugins/go-ipfs-plugin-i2p-gateway.so"

gx:
	go get -u github.com/whyrusleeping/gx
	go get -u github.com/whyrusleeping/gx-go
	go get -u github.com/karalabe/ungx

go-ipfs-plugin-i2p-gateway.so: plugin.go
	$(GOCC) build -buildmode=plugin
	chmod +x "go-ipfs-plugin-i2p-gateway.so"

plugin-libp2p:
	$(GOCC) build -a -tags libp2p -buildmode=plugin

deps:
	go get -u github.com/RTradeLtd/go-garlic-tcp-transport
	go get -u github.com/RTradeLtd/go-ipfs-plugin-i2p-gateway/config
	$(GX_GO_PATH) get github.com/RTradeLtd/go-ipfs-plugin-i2p-gateway

clobber:
	rm -rf $(GOPATH)

b:
	go build ./i2p

fmt:
	find ./i2p ./config -name '*.go' -exec gofmt -w {} \;

test:
	go test ./config -v
	go test ./i2p -v

vet:
	go vet ./config
	go vet ./i2p

import:
	gx import github.com/ipfs/go-ipfs

.PHONY: vendor

vendor:
	# Nuke vendor directory
	rm -rf vendor

	# Update standard dependencies
	#dep ensure -v -update
	dep ensure -v
	# Generate IPFS dependencies
	rm -rf vendor/github.com/ipfs/go-ipfs
	git clone https://github.com/ipfs/go-ipfs.git vendor/github.com/ipfs/go-ipfs
	( cd vendor/github.com/ipfs/go-ipfs ; gx install --local --nofancy )
	mv vendor/github.com/ipfs/go-ipfs/vendor/* vendor

	# Remove problematic dependencies
	find . -name test-vectors -type d -exec rm -r {} +
