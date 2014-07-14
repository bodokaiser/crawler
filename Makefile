SHELL := /bin/bash

test:
	@go test ./...

test-httpd:
	@go test -v ./httpd/...

test-utils:
	@go test -v ./utils/...

cmd-cli:
	@go run ./cmd/cli/main.go \
		--url http://www.satisfeet.me

cmd-http:
	@go run ./cmd/http/main.go

.PHONY: test-httpd test-utils
