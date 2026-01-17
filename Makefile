.PHONY: build clean test lint run

BINARY_NAME=relay
VERSION=$(shell cat VERSION)

build:
	go build -o bin/$(BINARY_NAME) ./cmd/relay

run:
	go run ./cmd/relay

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/

.DEFAULT_GOAL := build
