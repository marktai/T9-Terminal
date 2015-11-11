export GOPATH := $(shell pwd)
default: build

init:
	@rm -f bin/main
	@cd src/main && go get

build: init
	@go build -o bin/T9client src/main/main.go 

buildOSX: init
	@ GOOS=darwin go build -o bin/T9clientOSX src/main/main.go 

run: build
	@bin/T9client
