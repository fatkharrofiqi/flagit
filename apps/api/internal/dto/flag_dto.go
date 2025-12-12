package dto

import (
	"github.com/google/uuid"
)

type CreateFlagRequest struct {
	ProjectID   uuid.UUID `json:"project_id" validate:"required"`
	Key         string    `json:"key" validate:"required,min=1,max=100"`
	Description string    `json:"description" validate:"max=500"`
	Type        string    `json:"type" validate:"required,oneof=boolean string number json"`
}

type UpdateFlagRequest struct {
	Key         *string `json:"key" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Type        *string `json:"type" validate:"omitempty,oneof=boolean string number json"`
}

type CreateFlagValueRequest struct {
	FlagID  uuid.UUID `json:"flag_id" validate:"required"`
	EnvID   uuid.UUID `json:"env_id" validate:"required"`
	Value   string    `json:"value" validate:"required"`
	Enabled bool      `json:"enabled"`
}

type UpdateFlagValueRequest struct {
	Value   *string `json:"value"`
	Enabled *bool   `json:"enabled"`
}
