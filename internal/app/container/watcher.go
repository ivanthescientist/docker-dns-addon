package container

type EventHandler func(event Event)

type Watcher interface {
	Start() error
	Stop() error
}
