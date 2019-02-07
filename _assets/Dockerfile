# Build status-go in a Go builder container
FROM golang:1.9-alpine as builder

ARG build_tags

RUN apk add --no-cache make gcc musl-dev linux-headers

WORKDIR /go/src/github.com/mandrigin/status-go-bots
ADD . .
RUN make allbots

WORKDIR /go/src/github.com/status-im/status-go
ADD . vendor/github.com/status-im/status-go
RUN make statusgo

# Copy the binary to the second image
FROM alpine:latest

RUN apk add --no-cache ca-certificates bash

COPY --from=builder /go/src/github.com/mandrigin/status-go-bots/chanreader /usr/local/bin/

WORKDIR /etc/gs-bots/
COPY --from=builder /go/src/github.com/mandrigin/status-go-bots/_assets   ./_assets

RUN mkdir -p /data/sg_bots
VOLUME "/data/sg_bots"

# 30304 is used for Discovery v5
EXPOSE 8080 8545 30303 30303/udp 30304/udp
