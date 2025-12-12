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

type projectService struct {
	projectRepo repository.ProjectRepository
	sseService  SSEService
}

func NewProjectService(projectRepo repository.ProjectRepository, sseService SSEService) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		sseService:  sseService,
	}
}

func (s *projectService) CreateProject(ctx context.Context, req *dto.CreateProjectRequest) (*model.Project, error) {
	if req.Name == "" {
		return nil, errors.New("project name is required")
	}
	
	if len(req.Name) > 100 {
		return nil, errors.New("project name must be less than 100 characters")
	}
	
	if len(req.Description) > 500 {
		return nil, errors.New("project description must be less than 500 characters")
	}

	project := &model.Project{
		Name:        req.Name,
		Description: req.Description,
	}

	err := s.projectRepo.Create(ctx, project)
	if err != nil {
		return nil, err
	}

	// Broadcast SSE event
	eventData := model.ProjectEvent{
		ProjectID: project.ID,
		Name:      project.Name,
	}
	s.sseService.BroadcastEvent(sse.ProjectCreated, eventData)

	return project, nil
}

func (s *projectService) GetProject(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if project == nil {
		return nil, errors.New("project not found")
	}

	return project, nil
}

func (s *projectService) GetAllProjects(ctx context.Context) ([]model.Project, error) {
	projects, err := s.projectRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *projectService) UpdateProject(ctx context.Context, id uuid.UUID, req *dto.UpdateProjectRequest) (*model.Project, error) {
	exists, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if exists == nil {
		return nil, errors.New("project not found")
	}

	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("project name cannot be empty")
		}
		
		if len(*req.Name) > 100 {
			return nil, errors.New("project name must be less than 100 characters")
		}
	}

	if req.Description != nil && len(*req.Description) > 500 {
		return nil, errors.New("project description must be less than 500 characters")
	}

	project, err := s.projectRepo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// Broadcast SSE event
	eventData := model.ProjectEvent{
		ProjectID: project.ID,
		Name:      project.Name,
	}
	s.sseService.BroadcastEvent(sse.ProjectUpdated, eventData)

	return project, nil
}

func (s *projectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	exists, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if exists == nil {
		return errors.New("project not found")
	}

	err = s.projectRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Broadcast SSE event
	eventData := model.ProjectEvent{
		ProjectID: exists.ID,
		Name:      exists.Name,
	}
	s.sseService.BroadcastEvent(sse.ProjectDeleted, eventData)

	return nil
}
