.PHONY: test clean
default: test

proto:
	go install github.com/golang/protobuf/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	PATH=${PATH}:~/go/bin protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative \
		api/grpc/messages/*.proto \
		api/grpc/subscriptions/*.proto \
		api/grpc/writer/*.proto

test: proto
	go test -race ./...

docker:
	docker build -t awakari/core-tests .

run: docker
	docker run \
		-d \
		--name awakari-core-tests \
		-p 50051:50051 \
		--expose 50051 \
		awakari/core-tests

staging: docker
	./scripts/staging.sh

release: docker
	./scripts/release.sh

clean:
	go clean
