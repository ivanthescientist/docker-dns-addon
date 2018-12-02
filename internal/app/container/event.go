package container

// EventType represents type of container event within business domain
type EventType int

const (
	// Container started event type
	EventContainerStarted EventType = iota
	// Container stopped event type
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
