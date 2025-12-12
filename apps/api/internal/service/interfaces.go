package service

import (
	"api/internal/dto"
	"api/internal/model"
	"api/internal/sse"
	"context"

	"github.com/google/uuid"
)

type ProjectService interface {
	CreateProject(ctx context.Context, req *dto.CreateProjectRequest) (*model.Project, error)
	GetProject(ctx context.Context, id uuid.UUID) (*model.Project, error)
	GetAllProjects(ctx context.Context) ([]model.Project, error)
	UpdateProject(ctx context.Context, id uuid.UUID, req *dto.UpdateProjectRequest) (*model.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
}

type EnvironmentService interface {
	CreateEnvironment(ctx context.Context, req *dto.CreateEnvironmentRequest) (*model.Environment, error)
	GetEnvironment(ctx context.Context, id uuid.UUID) (*model.Environment, error)
	GetAllEnvironments(ctx context.Context) ([]model.Environment, error)
	GetProjectEnvironments(ctx context.Context, projectID uuid.UUID) ([]model.Environment, error)
	UpdateEnvironment(ctx context.Context, id uuid.UUID, req *dto.UpdateEnvironmentRequest) (*model.Environment, error)
	DeleteEnvironment(ctx context.Context, id uuid.UUID) error
}

type FlagService interface {
	CreateFlag(ctx context.Context, req *dto.CreateFlagRequest) (*model.Flag, error)
	GetFlag(ctx context.Context, id uuid.UUID) (*model.Flag, error)
	GetProjectFlags(ctx context.Context, projectID uuid.UUID, includeValues bool) ([]model.FlagWithValues, error)
	UpdateFlag(ctx context.Context, id uuid.UUID, req *dto.UpdateFlagRequest) (*model.Flag, error)
	DeleteFlag(ctx context.Context, id uuid.UUID) error

	// Flag value operations
	CreateOrUpdateFlagValue(ctx context.Context, req *dto.CreateFlagValueRequest) (*model.FlagValue, error)
	GetFlagValues(ctx context.Context, flagID uuid.UUID) ([]model.FlagValue, error)
	GetEnvironmentFlags(ctx context.Context, envID uuid.UUID) ([]model.FlagValue, error)
	UpdateFlagValue(ctx context.Context, id uuid.UUID, req *dto.UpdateFlagValueRequest) (*model.FlagValue, error)
	DeleteFlagValue(ctx context.Context, id uuid.UUID) error
}

type SSEService interface {
	BroadcastEvent(eventType sse.EventType, data interface{})
}
