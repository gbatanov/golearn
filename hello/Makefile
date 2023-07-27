.PHONY: all  prepare-bin build-linux

all: prepare-bin build-linux 

prepare-bin:
	rm -rf ./bin || true
	mkdir -p ./bin || true

build-linux:
	go build  -o bin/hello hello.go
