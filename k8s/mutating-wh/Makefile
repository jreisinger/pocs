#!make
SHELL := /bin/bash

test:
	go clean -testcache && go test -race ./...
.PHONY:test

e2e-test:
	./e2e/e2e
.PHONY:e2e-test

build: test
	go build
.PHONY:build
