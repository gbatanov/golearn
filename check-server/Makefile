.PHONY: prepare-bin, all. build-linux

all: prepare-bin build-linux 

prepare-bin:
	rm -rf ./bin || true
	mkdir -p ./bin || true

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/check-server .

#go build -ldflags "-H=windowsgui" -o bin/check-server.exe .

