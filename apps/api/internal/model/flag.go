package model

import (
	"time"

	"github.com/google/uuid"
)

type Flag struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ProjectID   uuid.UUID `json:"project_id" db:"project_id"`
	Key         string    `json:"key" db:"key"`
	Description string    `json:"description" db:"description"`
	Type        string    `json:"type" db:"type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type FlagValue struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FlagID    uuid.UUID `json:"flag_id" db:"flag_id"`
	EnvID     uuid.UUID `json:"env_id" db:"env_id"`
	Value     string    `json:"value" db:"value"`
	Enabled   bool      `json:"enabled" db:"enabled"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type FlagWithValues struct {
	Flag
	Values []FlagValue `json:"values,omitempty"`
}
