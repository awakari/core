FROM golang:1.20.4-alpine3.17 AS builder
WORKDIR /go/src/core
COPY . .
RUN apk add protoc protobuf-dev make git gcc musl-dev
ENTRYPOINT ["/usr/bin/make", "test"]
