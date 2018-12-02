package container

type EventType int

const (
	EventContainerStarted EventType = iota
	EventContainerStopped EventType = iota
)

type Event struct {
	Type      EventType
	Container Container
}

func (e Event) String() string {
	var eventType string
	switch e.Type {
	case EventContainerStarted:
		eventType = "Started"
	case EventContainerStopped:
		eventType = "Stopped"
	}
	return "Container " + eventType + ": " + e.Container.String()
}
