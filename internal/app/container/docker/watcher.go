package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
	"github.com/ivanthescientist/docker-dns-addon/internal/app/container"
	log "github.com/sirupsen/logrus"

	"sync"
)

type Watcher struct {
	doneCh  chan struct{}
	wg      *sync.WaitGroup
	logger  *log.Logger
	handler container.EventHandler
	client  *dockerclient.Client
	errCh   <-chan error
	eventCh <-chan events.Message
}

func NewWatcher(handler container.EventHandler, logger *log.Logger) (*Watcher, error) {
	doneCh := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(1)

	client, err := dockerclient.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		doneCh:  doneCh,
		wg:      wg,
		logger:  logger,
		handler: handler,
		client:  client,
	}, nil
}

func (w *Watcher) Start() error {
	w.logger.Print("Starting docker event watcher")
	filter := filters.NewArgs()
	filters.NewArgs().Add("type", "container")

	w.eventCh, w.errCh = w.client.Events(context.Background(), types.EventsOptions{
		Filters: filter,
	})

	err := w.loadInitialList()
	if err != nil {
		return err
	}

	go w.watch()
	return nil
}

func (w *Watcher) loadInitialList() error {
	containers, err := w.client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, c := range containers {
		var event container.Event
		event.Container, err = w.getContainer(c.ID)
		if err != nil {
			return err
		}
		event.Type = container.EventContainerStarted

		w.handler(event)
	}

	return nil
}

func (w *Watcher) Stop() error {
	w.logger.Print("Stopping docker event watcher")
	err := w.client.Close()
	if err != nil {
		return err
	}

	w.doneCh <- struct{}{}
	w.wg.Wait()
	return nil
}

func (w *Watcher) watch() {
	for {
		select {
		case event := <-w.eventCh:
			w.handleEvent(event)
		case err := <-w.errCh:
			w.logger.Fatalf("Docker event watcher recieved err in event stream: %s", err)
			return
		case <-w.doneCh:
			w.logger.Print("Done stopping docker event watcher")
			w.wg.Done()
			return
		}
	}
}

func (w *Watcher) handleEvent(message events.Message) {
	var event container.Event
	var err error

	switch message.Action {
	case "start":
		event.Type = container.EventContainerStarted
	case "die", "stop", "kill":
		event.Type = container.EventContainerStopped
	default:
		return
	}

	event.Container, err = w.getContainer(message.ID)
	if err != nil {
		w.logger.Errorf("Failed to fetch additional container info [%s]: %s", message.ID, err)
	}

	w.handler(event)
}

func (w *Watcher) getContainer(id string) (container.Container, error) {
	containerInfo, err := w.client.ContainerInspect(context.Background(), id)
	if err != nil {
		return container.Container{}, nil
	}

	name := transformContainerName(containerInfo.Name)

	return container.Container{
		ID:   id,
		Name: name,
		Addr: containerInfo.NetworkSettings.IPAddress,
	}, nil
}
