package service

import (
	"context"
	"errors"
	"api/internal/dto"
	"api/internal/model"
	"api/internal/repository"
	"api/internal/sse"
	"github.com/google/uuid"
)

type environmentService struct {
	envRepo    repository.EnvironmentRepository
	sseService SSEService
}

func NewEnvironmentService(envRepo repository.EnvironmentRepository, sseService SSEService) EnvironmentService {
	return &environmentService{
		envRepo:    envRepo,
		sseService: sseService,
	}
}

func (s *environmentService) CreateEnvironment(ctx context.Context, req *dto.CreateEnvironmentRequest) (*model.Environment, error) {
	if req.ProjectID == uuid.Nil {
		return nil, errors.New("project ID is required")
	}
	
	if req.Name == "" {
		return nil, errors.New("environment name is required")
	}
	
	if len(req.Name) > 100 {
		return nil, errors.New("environment name must be less than 100 characters")
	}

	env := &model.Environment{
		ProjectID: req.ProjectID,
		Name:      req.Name,
	}

	err := s.envRepo.Create(ctx, env)
	if err != nil {
		return nil, err
	}

	// Broadcast SSE event
	eventData := model.EnvironmentEvent{
		EnvironmentID: env.ID,
		ProjectID:     env.ProjectID,
		Name:          env.Name,
	}
	s.sseService.BroadcastEvent(sse.EnvironmentCreated, eventData)

	return env, nil
}

func (s *environmentService) GetEnvironment(ctx context.Context, id uuid.UUID) (*model.Environment, error) {
	env, err := s.envRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if env == nil {
		return nil, errors.New("environment not found")
	}

	return env, nil
}

func (s *environmentService) GetAllEnvironments(ctx context.Context) ([]model.Environment, error) {
	envs, err := s.envRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return envs, nil
}

func (s *environmentService) GetProjectEnvironments(ctx context.Context, projectID uuid.UUID) ([]model.Environment, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project ID is required")
	}

	envs, err := s.envRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return envs, nil
}

func (s *environmentService) UpdateEnvironment(ctx context.Context, id uuid.UUID, req *dto.UpdateEnvironmentRequest) (*model.Environment, error) {
	exists, err := s.envRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if exists == nil {
		return nil, errors.New("environment not found")
	}

	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("environment name cannot be empty")
		}
		
		if len(*req.Name) > 100 {
			return nil, errors.New("environment name must be less than 100 characters")
		}
	}

	env, err := s.envRepo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// Broadcast SSE event
	eventData := model.EnvironmentEvent{
		EnvironmentID: env.ID,
		ProjectID:     env.ProjectID,
		Name:          env.Name,
	}
	s.sseService.BroadcastEvent(sse.EnvironmentUpdated, eventData)

	return env, nil
}

func (s *environmentService) DeleteEnvironment(ctx context.Context, id uuid.UUID) error {
	exists, err := s.envRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if exists == nil {
		return errors.New("environment not found")
	}

	err = s.envRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Broadcast SSE event
	eventData := model.EnvironmentEvent{
		EnvironmentID: exists.ID,
		ProjectID:     exists.ProjectID,
		Name:          exists.Name,
	}
	s.sseService.BroadcastEvent(sse.EnvironmentDeleted, eventData)

	return nil
}
