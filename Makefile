SHELL := /bin/bash

test:
	@go test ./...

test-parser:
	@go test -v ./parser

.PHONY: test-parser
