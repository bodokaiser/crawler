SHELL := /bin/bash

boot:
	@go run main.go \
		--addr :3000 \
		http://www.google.com/search?q=Hello+World

test-parser:
	@go test -v ./parser

.PHONY: test
