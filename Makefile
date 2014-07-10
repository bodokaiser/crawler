SHELL := /bin/bash

boot:
	@go run main.go \
		http://www.google.com/

test:
	@go test ./...

test-split:
	@go test -v ./split/...

.PHONY: test-split
