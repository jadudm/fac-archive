SHELL:=/bin/bash

generate:
	pushd internal/archivedb ; sqlc generate ; popd


all: clean generate
	go build -tags "libsqlite3 linux"

clean:
	rm -f fac-tool