install: build
	cp rtinfo-dashboard ${GOPATH}/bin

generate:
	go generate

build: generate
	go build -o rtinfo-dashboard

test: generate
	go test -race -coverprofile=coverage.txt -covermode=atomic $(shell go list)

.PHONY: install build test rtinfo-dashboard