THIS_FILE := $(lastword $(MAKEFILE_LIST))
UNSAFE=no
GOCC ?= go
IPFS_VERSION=v0.4.19
IPFS_PATH?="$(HOME)/Workspace/ipfs"

# build plugin and ipfs daemon
.PHONY: build
build:
	mkdir -p build
	@$(MAKE) -f $(THIS_FILE) plugin-ipfs

# build the actual plugin
.PHONY: plugin-ipfs
plugin-ipfs:
	$(GOCC) build -o build/go-ipfs-plugin-i2p-gateway.so -buildmode=plugin
	chmod +x "build/go-ipfs-plugin-i2p-gateway.so"

# build ipfs daemon
.PHONY: ipfs
ipfs:
	( cd vendor/github.com/ipfs/go-ipfs/cmd/ipfs; go build -o ../../../../../../build/ipfs )
	install -m755 build/ipfs $(GOPATH)/bin

# clean up files
.PHONY: clean
clean:
	rm -rf build
	find . -name '*.i2pkeys' -exec rm -vf {} \;
	find . -name '*i2pconfig' -exec rm -vf {} \;

# install plugin to ipfs plugin folder
.PHONY: install
install:
	mkdir -p $(IPFS_PATH)/plugins
	install -Dm700 build/go-ipfs-plugin-i2p-gateway.so "$(IPFS_PATH)/plugins/go-ipfs-plugin-i2p-gateway.so"

# grab tooling to deal with gx
.PHONY: gx
gx:
	go get -u github.com/whyrusleeping/gx
	go get -u github.com/whyrusleeping/gx-go
	go get -u github.com/karalabe/ungx

# build libp2p plugin
.PHONY: plugin-libp2p
plugin-libp2p:
	$(GOCC) build -a -tags libp2p -buildmode=plugin

# build i2p  folder
.PHONY: build-i2p
build-i2p:
	go build ./i2p

# format i2p and config golang files
.PHONY: fmt
fmt:
	find ./i2p ./config -name '*.go' -exec gofmt -w {} \;
	gofmt -w *.go

# run tests
.PHONY: test
test:
	go test ./config -v
	go test ./i2p -v

# vet go code
.PHONY: vet
vet:
	go vet ./config
	go vet ./i2p

# import gx ipfs
.PHONY: import
import:
	gx import github.com/ipfs/go-ipfs

# completely rebuild vendor
.PHONY: vendor
vendor: vendor-dep vendor-ipfs

.PHONY: vendor-dep
vendor-dep:
	# Nuke vendor directory
	rm -rf vendor

	# Update standard dependencies
	#dep ensure -v -update
	dep ensure -v

.PHONY: vendor-ipfs
vendor-ipfs:
	# Generate IPFS dependencies
	rm -rf vendor/github.com/ipfs/go-ipfs
	git clone https://github.com/ipfs/go-ipfs.git vendor/github.com/ipfs/go-ipfs --branch $(IPFS_VERSION)
	( cd vendor/github.com/ipfs/go-ipfs ; gx install --local --nofancy )

	rsync -ravhp vendor/github.com/ipfs/go-ipfs/vendor/* vendor/
	rm -rf vendor/github.com/ipfs/go-ipfs/vendor/

	# Remove problematic dependencies
	find . -name test-vectors -type d -exec rm -r {} +
