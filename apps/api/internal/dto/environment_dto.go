package dto

import (
	"github.com/google/uuid"
)

type CreateEnvironmentRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	Name      string    `json:"name" validate:"required,min=1,max=100"`
}

type UpdateEnvironmentRequest struct {
	Name *string `json:"name" validate:"omitempty,min=1,max=100"`
}
