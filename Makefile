SHELL := /bin/bash

boot: install
	@go run main.go \
		--selector a \
		--attribute href \
		https://google.com?q=golang

install:
	@go get github.com/puerkitobio/goquery

.PHONY: test
