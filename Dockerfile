FROM golang:1.20.6-alpine3.18 AS builder
WORKDIR /go/src/core
COPY . .
RUN apk add protoc protobuf-dev make git gcc musl-dev
ENTRYPOINT ["/usr/bin/make", "test"]
