THIS_FILE := $(lastword $(MAKEFILE_LIST))
UNSAFE=no
GOCC ?= go


.PHONY: build
build:
	rm -rf build
	mkdir build
	@$(MAKE) -f $(THIS_FILE) plugin-ipfs
	@$(MAKE) -f $(THIS_FILE) ipfs

# build plugin and ipfs daemon
.PHONY: bbuild
bbuild:
	mkdir build
	@$(MAKE) -f $(THIS_FILE) plugin-ipfs
	@$(MAKE) -f $(THIS_FILE) ipfs

# build the actual plugin
.PHONY: plugin-ipfs
plugin-ipfs:
	$(GOCC) build -o build/go-ipfs-plugin-i2p-gateway.so -buildmode=plugin 
	chmod +x "build/go-ipfs-plugin-i2p-gateway.so"

# build ipfs daemon
.PHONY: ipfs
ipfs:
	( cd vendor/github.com/ipfs/go-ipfs/cmd/ipfs ; go build -o ../../../../../../build/ipfs ; cp ipfs $(GOPATH)/bin)

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

# fetch dependencies
.PHONY: deps
deps:
	go get -u github.com/RTradeLtd/go-garlic-tcp-transport
	go get -u github.com/RTradeLtd/go-ipfs-plugin-i2p-gateway/config
	$(GX_GO_PATH) get github.com/RTradeLtd/go-ipfs-plugin-i2p-gateway

# build i2p  folder
.PHONY: build-i2p
build-i2p:
	go build ./i2p

# format i2p and config golang files
.PHONY: fmt
fmt:
	find ./i2p ./config -name '*.go' -exec gofmt -w {} \;

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