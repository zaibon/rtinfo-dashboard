install: build
	cp rtinfo-dashboard ${GOPATH}/bin

generate:
	go generate

build: generate
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o rtinfo-dashboard

docker: build
	docker build -t zaibon/rtinfo-dashboard .

test: generate
	go test -race -coverprofile=coverage.txt -covermode=atomic $(shell go list)

clean:
	rm -f rtinfo-dashboard bindata.go

.PHONY: install build docker test rtinfo-dashboard clean
