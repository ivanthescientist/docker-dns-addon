# DockerDNS Addon
[![Go Report Card](https://goreportcard.com/badge/github.com/ivanthescientist/docker-dns-addon)](https://goreportcard.com/report/github.com/ivanthescientist/docker-dns-addon)

## What is this?
A minimalistic DNS server which provides default domains for deployed docker containers (provided they have networking enabled at all). 
The domain is inferred from container name, so if you have a container named `mysql_local` it will be available from domain `mysql_local.docker` 
provided you keep the default domain suffix (which is `.docker.`, notice that trailing dot is part of protocol, browsers just remove it).

## Requirements
- go 1.10+ (only tested version)
- docker (any version compatible with docker-compose 3.0)
- docker-compose (any version supporting 3.0 configuration)
- dep (any compatible version with go 1.10)

## Building
- `make install-dependencies` in the project root to install dependencies
- `make build` to build the project binaries
- `make install` to install built binaries as docker container
- `make rebuild` - clean and then build the project 
- `make reinstall` - remove existing installation and then rebuild and install binaries as docker container

## Installing from Dockerhub
For convenience there is a built image on dockerhub (beware that it is built locally on my machine for now and not on some trusted external CI):
`ivanthescientist/docker-dns-addon`
There is also a version of docker-compose.yaml using it `deployments/docker-compose-dockerhub.yaml` and several make targets: 
- `make install-dockerhub` - installs from dockerhub image
- `make uninstall-dockerhub` - removes deployment from dockerhub
- `make reinstall-dockerhub` - combined uninstall and install, just in case you want to scrap old version and install latest one

## Setup
By default the DNS server is deployed on a static IP `172.100.0.2` using a docker network `172.100.0.0/16`, this can be 
changed in docker-compose along with other configuration options. To start using the service after install it through `make install`
simply add `nameserver 172.100.0.2` to the end of your `/etc/resolvconf/resolv.conf.d/tail` file (create if necessary), then run
`sudo resolvconf -u` to load updated config (it could take a moment or two after updating the config for the system to start using your newly added DNS server). 

## Accessing your services after setup
If you've installed the service correctly, your browser should be able to resolve IPs of your services. The domains will be constructed 
automatically from **Container Name** and provided **Domain Suffix**. For example if you Domain Suffix is `.mydomain.` and you have a container named `mysql_local`
you can access it by domain `mysql_local.mydomain` from your browser or any tool in your system which uses default domain resolution services of your Linux OS. 

If service is deployed on host network, the resolved IP will be `127.0.0.1` always, 
however if any other network configuration is used, first network IP will be used instead. By default docker deploys services into its own networks which typically look
something like `172.x.x.x`, e.g. `172.17.0.0/16`, so the resulting IP of such container could be `172.17.0.2`. 


## Config Options:
  - `SERVER_HOST` - IP address to bind server to, default is `0.0.0.0` for local deployment and `172.100.0.2` for docker deployment
  - `SERVER_PORT` - port to bind server to `53` is default for docker deployment and local deployment uses `5300` because ports below 1024 are only available to root
  - `SERVER_PROTOCOL` - transport protocol to serve DNS from, possible values are `udp` and `tcp`
  - `DOMAIN_SUFFIX` - dot enclosed top level domain (e.g. `.com.`, `.org.` or `.docker.` in the default case), the trailing dot **must** be present. 