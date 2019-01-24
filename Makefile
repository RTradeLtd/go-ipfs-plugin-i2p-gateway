GOCC ?= go
IPFS_PATH ?= $(HOME)/.ipfs

VERSION="0.0.0"

#GOPATH=$(shell pwd)/go

GX_PATH=$(GOPATH)/bin/gx
GX_GO_PATH=$(GOPATH)/bin/gx-go

.PHONY: install build gx

build: example-plugin.so

install: build
	mkdir -p $(IPFS_PATH)/plugins
	install -Dm700 go-ipfs-plugin-i2p-gateway.so "$(IPFS_PATH)/plugins/example-plugin.so"

gx:
	go get -u github.com/whyrusleeping/gx
	go get -u github.com/whyrusleeping/gx-go

example-plugin.so: plugin.go
	$(GOCC) build -buildmode=plugin
	chmod +x "go-ipfs-plugin-i2p-gateway.so"

docker:
	docker build -t eyedeekay/go-ipfs-plugin-base .
	docker build -f Dockerfile.build -t eyedeekay/go-ipfs-plugin-build .

deps:
	$(GX_GO_PATH) get github.com/rtradeltd/go-ipfs-plugin-i2p-gateway

b:
	go build ./i2p

dep:
	go get -u "github.com/rtradeltd/go-garlic-tcp-transport"
	$(GX_GO_PATH) get "github.com/ipfs/go-ipfs"
	$(GX_GO_PATH) get "github.com/ipfs/go-ipfs-config"

dep2:
	$(GX_PATH) repo add plugin QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r
	$(GX_PATH) repo add config QmRd5T3VmYoX6jaNoZovFRQcwWHJqHgTVQTs1Qz92ELJ7C67ty78

import:
	$(GX_PATH) import "github.com/ipfs/go-ipfs"
	$(GX_PATH) import "github.com/ipfs/go-ipfs-config"

setup:
	$(GX_PATH) import QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r
	$(GX_PATH) import QmRd5T3VmYoX6jaNoZovFRQcwWHJqHgTVQTs1Qz92ELJ7C

dep3:
	$(GX_GO_PATH) import --yesall "github.com/ipfs/go-ipfs"
	$(GX_GO_PATH) import --yesall "github.com/ipfs/go-ipfs-config"

fmt:
	find ./i2p ./config -name '*.go' -exec gofmt -w {} \;

gx-install:
	$(GX_PATH) install


