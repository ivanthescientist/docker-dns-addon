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

rebuild: clean build
.PHONY: rebuild

install:
	sudo docker-compose -f ./deployments/docker-compose.yaml up -d --build
.PHONY: install

uninstall:
	sudo docker-compose -f ./deployments/docker-compose.yaml down
.PHONY: uninstall

reinstall: uninstall rebuild install
.PHONY: reinstall