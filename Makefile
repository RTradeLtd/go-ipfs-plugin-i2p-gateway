GOCC ?= go
IPFS_PATH ?= $(HOME)/.ipfs

VERSION="0.0.0"

GOPATH=$(shell pwd)/go

GX_PATH=$(GOPATH)/bin/gx
UNGX_PATH=$(GOPATH)/bin/ungx
GX_GO_PATH=$(GOPATH)/bin/gx-go

.PHONY: install build gx

build: go-ipfs-plugin-i2p-gateway.so

clean:
	rm -f go-ipfs-plugin-i2p-gateway.so

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

docker:
	docker build -t eyedeekay/go-ipfs-plugin-base .
	docker build -f Dockerfile.build -t eyedeekay/go-ipfs-plugin-build .

deps:
	go get -u github.com/rtradeltd/go-ipfs-plugin-i2p-gateway/config
	go get -u github.com/rtradeltd/go-ipfs-plugin-i2p-gateway/i2p

clobber:
	rm -rf $(GOPATH)

b:
	go build ./i2p

fmt:
	find ./i2p ./config -name '*.go' -exec gofmt -w {} \;

gx-install:
	$(GX_PATH) install

test:
	go test ./config
	go test ./i2p -v

vet:
	go vet ./config
	go vet ./i2p
