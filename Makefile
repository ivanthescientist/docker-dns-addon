all: clean build
.PHONY: all

install-dependencies:
	dep ensure
.PHONY: install-dependencies

clean:
	rm -rf bin/*
.PHONY: clean

build:
	go build -o bin/dockerdns cmd/dockerdns/main.go
.PHONY: build

install: clean build
	sudo docker-compose -f ./deployments/docker-compose.yaml up -d --build
.PHONY: install

uninstall:
	sudo docker-compose -f ./deployments/docker-compose.yaml down
.PHONY: uninstall