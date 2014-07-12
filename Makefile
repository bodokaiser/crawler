SHELL := /bin/bash

boot:
	@go run main.go \
		http://www.google.com/

test:
	@go test ./...

test-httpd:
	@go test -v ./httpd/...

test-parser:
	@go test -v ./parser/...

.PHONY: test-httpd test-parser
