.PHONY: test clean
default: test

test:
	CGO_ENABLED=1 go test -race ./...

testperfe2e:
	go test -v -run Test_Perf_EndToEnd -timeout 168h ./...

docker:
	docker build -t awakari/core-tests .

staging: docker
	./scripts/staging.sh

release: docker
	./scripts/release.sh

clean:
	go clean
