SHELL := /bin/bash

boot:
	@go run main.go \
		http://www.google.com/

test:
	@go test ./...

test-parser:
	@go test -v ./parser

.PHONY: test-parser
