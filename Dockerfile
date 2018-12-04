FROM debian:stretch-slim
RUN  apt-get update ; apt-get install -y ca-certificates
COPY bin/dockerdns /dockerdns
RUN mkdir log ; chmod 775 /dockerdns
EXPOSE 53 5300
ENTRYPOINT /dockerdns