SHELL := /bin/bash

test:
	@go test ./...

test-net:
	@go test -v ./net/...

test-text:
	@go test -v ./text/...

cmd-cli:
	@go run ./cmd/cli/main.go \
		--url http://www.satisfeet.me

cmd-http:
	@go run ./cmd/http/main.go

.PHONY: test-net text-text
