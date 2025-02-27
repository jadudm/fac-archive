SHELL:=/bin/bash

all: clean
	pushd internal/sqlite ; sqlc generate ; popd
	go build -tags "libsqlite3 linux"

clean:
	rm -f fac-tool