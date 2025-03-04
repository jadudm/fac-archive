SHELL:=/bin/bash

generate:
	pushd internal/archivedb ; sqlc generate ; popd

windows: generate
	rm -f fac-archive-win-amd.exe
	rm -f fac-archive-win-arm.exe
	GOOS=windows GOARCH=amd64 go build -o fac-archive-win-amd.exe
	GOOS=windows GOARCH=arm64 go build -o fac-archive-win-arm.exe

mac: generate
	rm -f fac-archive-mac-amd
	rm -f fac-archive-mac-arm	
	GOOS=darwin GOARCH=amd64 go build -o fac-archive-mac-amd
	GOOS=darwin GOARCH=arm64 go build -o fac-archive-mac-arm
	
linux: generate
	rm -f fac-archive-linux-amd
	rm -f fac-archive-linux-arm
	GOOS=linux GOARCH=amd64 go build -o fac-archive-linux-amd
	GOOS=linux GOARCH=arm64 go build -o fac-archive-linux-arm

all: linux mac windows