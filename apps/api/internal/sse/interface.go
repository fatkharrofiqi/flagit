package sse

// EventType represents the type of event being broadcast
type EventType string

const (
	// Project events
	ProjectCreated   EventType = "project.created"
	ProjectUpdated   EventType = "project.updated"
	ProjectDeleted   EventType = "project.deleted"

	// Environment events
	EnvironmentCreated EventType = "environment.created"
	EnvironmentUpdated EventType = "environment.updated"
	EnvironmentDeleted EventType = "environment.deleted"

	// Flag events
	FlagCreated       EventType = "flag.created"
	FlagUpdated       EventType = "flag.updated"
	FlagDeleted       EventType = "flag.deleted"
	FlagValueCreated  EventType = "flag.value.created"
	FlagValueUpdated  EventType = "flag.value.updated"
	FlagValueDeleted  EventType = "flag.value.deleted"
)

// Service defines the interface for server-sent events
type Service interface {
	BroadcastEvent(eventType EventType, data interface{})
}
