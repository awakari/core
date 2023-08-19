FROM golang:1.20.7-alpine3.18
WORKDIR /go/src/core
COPY . .
RUN apk add --update --no-cache protoc protobuf-dev make git gcc musl-dev
RUN go get github.com/fullstorydev/grpcurl/...
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl
RUN go get -t github.com/awakari/core/test
ENTRYPOINT ["/usr/bin/make", "test"]
