all: clean build
.PHONY: all

clean:
	rm -rf bin/*
.PHONY: clean

build:
	go build -o bin/dockerdns cmd/dockerdns/main.go
.PHONY: build