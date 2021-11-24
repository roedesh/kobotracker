ENTRY=main.go
BINARY=cryptokobo

GOOS=linux
GOARCH=arm
CGO_ENABLED=1
CC=$(HOME)/x-tools/arm-kobo-linux-gnueabihf/bin/arm-kobo-linux-gnueabihf-gcc
CXX=$(HOME)/x-tools/arm-kobo-linux-gnueabihf/bin/arm-kobo-linux-gnueabihf-g++

TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
VERSION := $(TAG:v%=%)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
ifeq ($(VERSION),)
    VERSION = $(COMMIT)-$(DATE)
endif

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) CC=$(CC) CXX=$(CXX) go build -ldflags="-X 'main.version=$(VERSION)'" -o $(BINARY) $(ENTRY)

package:
	tar -czvf cryptokobo.tar.gz ${BINARY} run.sh assets/*

release: build package clean

clean:
	go clean
	rm -f $(BINARY)
