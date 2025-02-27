SHELL:=/bin/bash

all:
	pushd cmd/fac-copy ; go build . ; popd