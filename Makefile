SHELL:=/bin/bash

generate:
	pushd internal/archivedb ; sqlc generate ; popd

windows: generate
	rm -f fac-tool-win-amd.exe
	rm -f fac-tool-win-arm.exe
	GOOS=windows GOARCH=amd64 go build -o fac-tool-win-amd.exe
	GOOS=windows GOARCH=arm64 go build -o fac-tool-win-arm.exe

mac: generate
	rm -f fac-tool-mac-amd
	rm -f fac-tool-mac-arm	
	GOOS=darwin GOARCH=amd64 go build -o fac-tool-mac-amd
	GOOS=darwin GOARCH=arm64 go build -o fac-tool-mac-arm
	
linux: generate
	rm -f fac-tool-linux-amd
	rm -f fac-tool-linux-arm
	GOOS=linux GOARCH=amd64 go build -o fac-tool-linux-amd
	GOOS=linux GOARCH=arm64 go build -o fac-tool-linux-arm

all: linux mac windows