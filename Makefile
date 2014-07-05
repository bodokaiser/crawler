SHELL := /bin/bash

boot: install
	@go run main.go \
		--addr :3000 \
		http://www.google.com/search?q=Hello+World

install:
	@go get github.com/puerkitobio/goquery

.PHONY: test
