TARGET=bin/cryptokobo

GOOS=linux
GOARCH=arm
CGO_ENABLED=1
CC=$(HOME)/x-tools/arm-kobo-linux-gnueabihf/bin/arm-kobo-linux-gnueabihf-gcc
CXX=$(HOME)/x-tools/arm-kobo-linux-gnueabihf/bin/arm-kobo-linux-gnueabihf-g++

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) CC=$(CC) CXX=$(CXX) go build -o $(TARGET) main.go

clean:
	go clean
	rm -f $(TARGET)
