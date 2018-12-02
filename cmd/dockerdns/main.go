package main

import (
	"github.com/ivanthescientist/docker-dns-addon/internal/app/config"
	"github.com/ivanthescientist/docker-dns-addon/internal/app/container"
	"github.com/ivanthescientist/docker-dns-addon/internal/app/container/docker"
	"github.com/ivanthescientist/docker-dns-addon/internal/app/dns"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var cfg *config.Config
var logger *log.Logger

func init() {
	cfg = config.GetFromEnv()
	logger = log.New()
}

func main() {
	var err error
	var watcher container.Watcher
	var domainRegistry *dns.DomainRegistry
	var server *dns.Server

	domainRegistry = dns.NewDomainRegistry(logger, cfg.DomainSuffix)

	watcher, err = docker.NewWatcher(domainRegistry.HandleEvent, logger)
	if err != nil {
		logger.Fatal(err)
	}

	err = watcher.Start()
	if err != nil {
		logger.Fatal(err)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		logger.Info("Shutting down...")

		err = watcher.Stop()
		if err != nil {
			logger.Error(err)
		}

		err = server.Shutdown()
		if err != nil {
			logger.Error(err)
		}
	}()

	server = dns.NewServer(logger, cfg.ServerHost, cfg.ServerPort, cfg.ServerProtocol, domainRegistry)
	err = server.ListenAndServe()
	if err != nil {
		logger.Info(err)
	}
}
