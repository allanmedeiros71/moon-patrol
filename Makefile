VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
PKG := github.com/allanmedeiros/moonpatrol

.DEFAULT_GOAL := build
.PHONY: build run test lint fmt tools clean

build:
	CGO_ENABLED=1 mise x -- go build -trimpath -ldflags "-s -w -X $(PKG)/internal/buildinfo.Version=$(VERSION)" -o bin/moonpatrol ./cmd/moonpatrol

run: build
	./bin/moonpatrol -windowed

test:
	mise x -- go test -race ./...

lint:
	mise x -- golangci-lint run ./...

fmt:
	mise x -- gofmt -w .

tools:
	mise install

clean:
	rm -rf bin dist coverage.*
