VERSION := $(shell git describe --exact-match --tags HEAD)
ifeq ($(VERSION),)
	VERSION := 'v0.0.0'
endif
VERSION_FLAGS := -X "github.com/saladtechnologies/salad-cloud-job-queue-worker/pkg/workers.Version=$(VERSION)"
LDFLAGS := -ldflags='$(VERSION_FLAGS)'

.PHONY: build
build: clean
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./build/ $(LDFLAGS) ./cmd/salad-http-job-queue-worker
	cd ./build; tar -czvf salad-http-job-queue-worker_x86_64.tar.gz salad-http-job-queue-worker
	cd ./build; sha256sum salad-http-job-queue-worker_x86_64.tar.gz > salad-http-job-queue-worker_x86_64.tar.gz.sha256

.PHONY: clean
clean:
	rm -rf ./build

.PHONY: generate
generate:
	buf generate
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run ./...
