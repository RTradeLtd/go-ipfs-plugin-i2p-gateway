GOCC ?= go
IPFS_PATH ?= $(HOME)/.ipfs

VERSION="0.0.0"

GOPATH=$(shell pwd)/go

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

deps: dep
	$(GX_GO_PATH) get github.com/rtradeltd/go-ipfs-plugin-i2p-gateway

b:
	go build ./i2p

dep:
	go get -u "github.com/ipfs/go-ipfs"
	$(GX_GO_PATH) get "github.com/ipfs/go-ipfs"
	$(GX_GO_PATH) get "github.com/ipfs/go-ipfs-config"
	$(GX_PATH) repo add plugin QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r
	$(GX_PATH) repo add config QmRd5T3VmYoX6jaNoZovFRQcwWHJqHgTVQTs1Qz92ELJ7C
	#$(GX_GO_PATH) import "github.com/ipfs/go-ipfs" #QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r
	$(GX_GO_PATH) import "github.com/ipfs/go-ipfs-config" #QmRd5T3VmYoX6jaNoZovFRQcwWHJqHgTVQTs1Qz92ELJ7C

fmt:
	find ./i2p -name '*.go' -exec gofmt -w {} \;

gx-install:
	$(GX_PATH) install
