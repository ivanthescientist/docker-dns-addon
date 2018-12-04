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

install-dockerhub:
	sudo docker-compose -f ./deployments/docker-compose.dockerhub.yaml up -d --build
.PHONY: install-dockerhub

uninstall-dockerhub:
	sudo docker-compose -f ./deployments/docker-compose.dockerhub.yaml down
.PHONY: uninstall-dockerhub

reinstall-dockerhub: uninstall-dockerhub install-dockerhub
.PHONY: reinstall-dockerhub