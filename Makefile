VERSION := $(shell git describe --tags --always --dirty)

.PHONY: build test vet check

build:
	go build -ldflags "-X main.version=$(VERSION)" .

test:
	go test ./...

vet:
	go vet ./...

check:
	gofmt -w .
	go test -race ./...
	go vet ./...
	go build -ldflags "-X main.version=$(VERSION)" .
