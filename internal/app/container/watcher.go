package container

// EventHandler represents container event handler function type
type EventHandler func(event Event)

// Watcher represents generic container event Watcher, it should be possible to create any impl for it (e.g. docker, rkt)
type Watcher interface {
	Start() error
	Stop() error
}
