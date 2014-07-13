SHELL := /bin/bash

boot:
	@go run main.go \
		http://www.google.com/

test:
	@go test ./...

test-conf:
	@go test -v ./conf/...

test-httpd:
	@go test -v ./httpd/...

test-utils:
	@go test -v ./utils/...

.PHONY: test-httpd test-utils
