# Build status-go in a Go builder container
FROM golang:1.9-alpine as builder

ARG build_tags

RUN apk add --no-cache make gcc musl-dev linux-headers


WORKDIR /go/src/github.com/status-im/status-go
ADD vendor/github.com/status-im/status-go .
RUN make statusgo


WORKDIR /go/src/github.com/mandrigin/status-go-bots
ADD cmd cmd
ADD vendor vendor
ADD Makefile Makefile
RUN make allbots

ADD _assets _assets
ADD run-reader.sh run-reader.sh

# Copy the binary to the second image
FROM alpine:latest

RUN apk add --no-cache ca-certificates bash

COPY --from=builder /go/src/github.com/status-im/status-go/build/bin/statusd /usr/local/bin/
COPY --from=builder /go/src/github.com/mandrigin/status-go-bots/chanreader /usr/local/bin/

WORKDIR /etc/gs-bots/
COPY --from=builder /go/src/github.com/mandrigin/status-go-bots/_assets   ./_assets
COPY --from=builder /go/src/github.com/status-im/status-go/ethdenver.json ./_assets/ethdenver.json
COPY --from=builder /go/src/github.com/mandrigin/status-go-bots/run-reader.sh ./_assets/run-reader.sh

RUN chmod a+x ./_assets/run-reader.sh

RUN mkdir -p /data/sg_bots
VOLUME "/data/sg_bots"

# 30304 is used for Discovery v5
EXPOSE 8080 8545 30303 30303/udp 30304/udp
