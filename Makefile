export GOPATH := $(shell pwd)
default: build

init:
	@rm -f bin/main
	@cd src/main && go get

build: init
	rm -f bin/T9client
	go build -o bin/T9client src/main/main.go 

linux: init
	rm -f bin/T9clientLinux
	GOOS=linux go build -o bin/T9clientLinux src/main/main.go 

osx: init
	rm -f bin/T9clientOSX
	GOOS=darwin go build -o bin/T9clientOSX src/main/main.go 

windows: init
	rm -f bin/T9clientWindows.exe
	GOOS=windows go build -o bin/T9clientWindows.exe src/main/main.go

all: buildLinux buildOSX buildWindows

run: build
	bin/T9client
