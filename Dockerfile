FROM alpine:edge
ARG GOPATH=/usr/src/ipfs-plugin/go
ENV GOPATH=/usr/src/ipfs-plugin/go
RUN apk update
RUN apk add -u go make git ca-certificates musl-dev
RUN go get -u github.com/whyrusleeping/gx
RUN go get -u github.com/whyrusleeping/gx-go
COPY . /usr/src/ipfs-plugin
WORKDIR /usr/src/ipfs-plugin
RUN ls $GOPATH
