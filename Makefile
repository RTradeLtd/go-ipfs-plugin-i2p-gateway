GOCC ?= go
IPFS_PATH ?= $(HOME)/.ipfs

#VERSION="0.0.0"

#GOPATH=$(shell pwd)/go

GX_PATH=$(GOPATH)/bin/gx
GX_GO_PATH=$(GOPATH)/bin/gx-go

.PHONY: install build gx

build: example-plugin.so

install: build
	install -Dm700 example-plugin.so "$(IPFS_PATH)/plugins/example-plugin.so"

gx:
	go get -u github.com/whyrusleeping/gx
	go get -u github.com/whyrusleeping/gx-go

example-plugin.so: plugin.go
	$(GOCC) build -buildmode=plugin
	chmod +x "$@"

docker:
	docker build -t eyedeekay/go-ipfs-plugin-base .
	docker build -f Dockerfile.build -t eyedeekay/go-ipfs-plugin-build .

deps:
	$(GX_GO_PATH) get github.com/rtradeltd/go-ipfs-plugin-i2p-gateway

b:
	go build ./i2p

dep:
	gx-go get "github.com/ipfs/go-ipfs-config"

fmt:
	find ./i2p -name '*.go' -exec gofmt -w {} \;
