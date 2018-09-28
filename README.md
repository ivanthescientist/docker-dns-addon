# What is this?
This is an WIP docker addon/daemon that will in the future provide
local DNS server for locally deployed docker containers, the main idea is
to allow host machine to resolve containers by their domain, for now this is done
with `container-name.docker` domain.

At the moment this is still deeply WIP, so there are no docker-compose or even
a Dockerfile to properly deploy this on a local machine, but you can just
run it as root locally and it will start a DNS server on `0.0.0.0:53`, then
you can just add `nameserver 127.0.0.1` to `/etc/resolv.conf` and it should
start resolving domains. It automatically listens to local docker events
and updates its DNS records accordingly whenever a container is started or stopped/dies.
