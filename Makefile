GOBIN ?= $(GOPATH)/bin
GOSUM := $(shell which gosum)

all: install

install: go.sum
	GO111MODULE=on go install

build: go.sum
	GO111MODULE=on go build -o commercio-network-chain-installer

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

build-darwin:
	env GO111MODULE=on GOOS=darwin GOARCH=386 go build -o ./build/darwin/commercio-network-chain-installer-darwin-386
	env GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o ./build/darwin/commercio-network-chain-installer-darwin-amd64

build-linux:
	env GO111MODULE=on GOOS=linux GOARCH=386 go build -o ./build/linux/commercio-network-chain-installer-linux-386
	env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ./build/linux/commercio-network-chain-installer-linux-amd64

build-windows:
	env GO111MODULE=on GOOS=windows GOARCH=386 go build -o ./build/windows/commercio-network-chain-installer-windows-386
	env GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o ./build/linux/commercio-network-chain-installer-windows-amd64

build-all:
	make build-darwin
	make build-linux
	make build-windows
