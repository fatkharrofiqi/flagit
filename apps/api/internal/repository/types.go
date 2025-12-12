package repository

import (
	"context"
	"api/internal/dto"
	"api/internal/model"
	"github.com/google/uuid"
)

// Repository interfaces for dependency injection

type ProjectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error)
	GetAll(ctx context.Context) ([]model.Project, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateProjectRequest) (*model.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type EnvironmentRepository interface {
	Create(ctx context.Context, env *model.Environment) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Environment, error)
	GetAll(ctx context.Context) ([]model.Environment, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Environment, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateEnvironmentRequest) (*model.Environment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type FlagRepository interface {
	Create(ctx context.Context, flag *model.Flag) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Flag, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Flag, error)
	GetWithValuesByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.FlagWithValues, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateFlagRequest) (*model.Flag, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type FlagValueRepository interface {
	Create(ctx context.Context, flagValue *model.FlagValue) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.FlagValue, error)
	GetByFlagID(ctx context.Context, flagID uuid.UUID) ([]model.FlagValue, error)
	GetByEnvID(ctx context.Context, envID uuid.UUID) ([]model.FlagValue, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateFlagValueRequest) (*model.FlagValue, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
