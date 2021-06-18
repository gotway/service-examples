.PHONY: clean deps deps-sync fmt-srv fmt-cli fmt proto build srv run cli cli.mgw

LOCAL_ARGS ?= -server localhost:11000 -tls true
PRO_ARGS ?= -server gotway.duckdns.org:4433 -tls true

all: build
clean:
	rm -rf bin/
deps:
	go get -u -v
deps-sync:
	go mod vendor
fmt-srv: 
	go fmt server/*.go 
fmt-cli: 
	go fmt client/*.go
fmt: fmt-srv fmt-cli
build: fmt-srv clean
	go build -o bin/route -v .
srv: build
	./bin/route
run: srv
cli: fmt-cli
	go run client/*.go $(LOCAL_ARGS)
cli-pro: fmt-cli
	go run client/*.go $(PRO_ARGS)
