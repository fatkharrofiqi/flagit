package model

import (
	"time"

	"api/internal/sse"
	"github.com/google/uuid"
)

// SSEEvent represents a server-sent event
type SSEEvent struct {
	Type      sse.EventType `json:"type"`
	Timestamp time.Time     `json:"timestamp"`
	Data      interface{}   `json:"data"`
}

// ProjectEvent represents a project-related event payload
type ProjectEvent struct {
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
}

// EnvironmentEvent represents an environment-related event payload
type EnvironmentEvent struct {
	EnvironmentID uuid.UUID `json:"environment_id"`
	ProjectID     uuid.UUID `json:"project_id"`
	Name          string    `json:"name"`
}

// FlagEvent represents a flag-related event payload
type FlagEvent struct {
	FlagID    uuid.UUID `json:"flag_id"`
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
}

// FlagValueEvent represents a flag value-related event payload
type FlagValueEvent struct {
	FlagValueID  uuid.UUID `json:"flag_value_id"`
	FlagID       uuid.UUID `json:"flag_id"`
	EnvironmentID uuid.UUID `json:"environment_id"`
}
