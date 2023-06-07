.PHONY: test clean
default: test

test:
	CGO_ENABLED=1 go test -race ./...

docker:
	docker build -t awakari/core-tests .

run: docker
	docker run \
		--name awakari-core-tests \
		--network host \
		awakari/core-tests

staging: docker
	./scripts/staging.sh

release: docker
	./scripts/release.sh

clean:
	go clean
