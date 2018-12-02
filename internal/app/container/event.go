package container

// EventType represents type of container event within business domain
type EventType int

const (
	// EventContainerStarted indicates container has started and ready to be added to DNS
	EventContainerStarted EventType = iota
	// EventContainerStopped indicates container has stopped and should be removed from DNS
	EventContainerStopped EventType = iota
)

// Event represents generic container event
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
