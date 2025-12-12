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

type flagService struct {
	flagRepo      repository.FlagRepository
	flagValueRepo repository.FlagValueRepository
	sseService    SSEService
}

func NewFlagService(flagRepo repository.FlagRepository, flagValueRepo repository.FlagValueRepository, sseService SSEService) FlagService {
	return &flagService{
		flagRepo:      flagRepo,
		flagValueRepo: flagValueRepo,
		sseService:    sseService,
	}
}

func (s *flagService) CreateFlag(ctx context.Context, req *dto.CreateFlagRequest) (*model.Flag, error) {
	if req.ProjectID == uuid.Nil {
		return nil, errors.New("project ID is required")
	}
	
	if req.Key == "" {
		return nil, errors.New("flag key is required")
	}
	
	if len(req.Key) > 100 {
		return nil, errors.New("flag key must be less than 100 characters")
	}
	
	if len(req.Description) > 500 {
		return nil, errors.New("flag description must be less than 500 characters")
	}
	
	if req.Type == "" {
		return nil, errors.New("flag type is required")
	}
	
	validTypes := []string{"boolean", "string", "number", "json"}
	isValidType := false
	for _, validType := range validTypes {
		if req.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return nil, errors.New("flag type must be one of: boolean, string, number, json")
	}

	flag := &model.Flag{
		ProjectID:   req.ProjectID,
		Key:         req.Key,
		Description: req.Description,
		Type:        req.Type,
	}

	err := s.flagRepo.Create(ctx, flag)
	if err != nil {
		return nil, err
	}

	// Broadcast SSE event
	eventData := model.FlagEvent{
		FlagID:    flag.ID,
		ProjectID: flag.ProjectID,
		Name:      flag.Description, // Using description as name since flag model doesn't have name
		Key:       flag.Key,
	}
	s.sseService.BroadcastEvent(sse.FlagCreated, eventData)

	return flag, nil
}

func (s *flagService) GetFlag(ctx context.Context, id uuid.UUID) (*model.Flag, error) {
	flag, err := s.flagRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if flag == nil {
		return nil, errors.New("flag not found")
	}

	return flag, nil
}

func (s *flagService) GetProjectFlags(ctx context.Context, projectID uuid.UUID, includeValues bool) ([]model.FlagWithValues, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project ID is required")
	}

	if includeValues {
		return s.flagRepo.GetWithValuesByProjectID(ctx, projectID)
	}

	flags, err := s.flagRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var flagsWithValues []model.FlagWithValues
	for _, flag := range flags {
		flagsWithValues = append(flagsWithValues, model.FlagWithValues{
			Flag:   flag,
			Values: []model.FlagValue{},
		})
	}

	return flagsWithValues, nil
}

func (s *flagService) UpdateFlag(ctx context.Context, id uuid.UUID, req *dto.UpdateFlagRequest) (*model.Flag, error) {
	exists, err := s.flagRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if exists == nil {
		return nil, errors.New("flag not found")
	}

	if req.Key != nil {
		if *req.Key == "" {
			return nil, errors.New("flag key cannot be empty")
		}
		
		if len(*req.Key) > 100 {
			return nil, errors.New("flag key must be less than 100 characters")
		}
	}

	if req.Description != nil && len(*req.Description) > 500 {
		return nil, errors.New("flag description must be less than 500 characters")
	}

	if req.Type != nil {
		if *req.Type == "" {
			return nil, errors.New("flag type cannot be empty")
		}
		
		validTypes := []string{"boolean", "string", "number", "json"}
		isValidType := false
		for _, validType := range validTypes {
			if *req.Type == validType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			return nil, errors.New("flag type must be one of: boolean, string, number, json")
		}
	}

	flag, err := s.flagRepo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// Broadcast SSE event
	eventData := model.FlagEvent{
		FlagID:    flag.ID,
		ProjectID: flag.ProjectID,
		Name:      flag.Description,
		Key:       flag.Key,
	}
	s.sseService.BroadcastEvent(sse.FlagUpdated, eventData)

	return flag, nil
}

func (s *flagService) DeleteFlag(ctx context.Context, id uuid.UUID) error {
	exists, err := s.flagRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if exists == nil {
		return errors.New("flag not found")
	}

	err = s.flagRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Broadcast SSE event
	eventData := model.FlagEvent{
		FlagID:    exists.ID,
		ProjectID: exists.ProjectID,
		Name:      exists.Description,
		Key:       exists.Key,
	}
	s.sseService.BroadcastEvent(sse.FlagDeleted, eventData)

	return nil
}

// Flag value operations
func (s *flagService) CreateOrUpdateFlagValue(ctx context.Context, req *dto.CreateFlagValueRequest) (*model.FlagValue, error) {
	if req.FlagID == uuid.Nil {
		return nil, errors.New("flag ID is required")
	}
	
	if req.EnvID == uuid.Nil {
		return nil, errors.New("environment ID is required")
	}
	
	if req.Value == "" {
		return nil, errors.New("flag value is required")
	}

	flagValue := &model.FlagValue{
		FlagID:  req.FlagID,
		EnvID:   req.EnvID,
		Value:   req.Value,
		Enabled: req.Enabled,
	}

	err := s.flagValueRepo.Create(ctx, flagValue)
	if err != nil {
		return nil, err
	}

	// Broadcast SSE event
	eventData := model.FlagValueEvent{
		FlagValueID:  flagValue.ID,
		FlagID:       flagValue.FlagID,
		EnvironmentID: flagValue.EnvID,
	}
	s.sseService.BroadcastEvent(sse.FlagValueCreated, eventData)

	return flagValue, nil
}

func (s *flagService) GetFlagValues(ctx context.Context, flagID uuid.UUID) ([]model.FlagValue, error) {
	if flagID == uuid.Nil {
		return nil, errors.New("flag ID is required")
	}

	values, err := s.flagValueRepo.GetByFlagID(ctx, flagID)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (s *flagService) GetEnvironmentFlags(ctx context.Context, envID uuid.UUID) ([]model.FlagValue, error) {
	if envID == uuid.Nil {
		return nil, errors.New("environment ID is required")
	}

	values, err := s.flagValueRepo.GetByEnvID(ctx, envID)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (s *flagService) UpdateFlagValue(ctx context.Context, id uuid.UUID, req *dto.UpdateFlagValueRequest) (*model.FlagValue, error) {
	exists, err := s.flagValueRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if exists == nil {
		return nil, errors.New("flag value not found")
	}

	flagValue, err := s.flagValueRepo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// Broadcast SSE event
	eventData := model.FlagValueEvent{
		FlagValueID:  flagValue.ID,
		FlagID:       flagValue.FlagID,
		EnvironmentID: flagValue.EnvID,
	}
	s.sseService.BroadcastEvent(sse.FlagValueUpdated, eventData)

	return flagValue, nil
}

func (s *flagService) DeleteFlagValue(ctx context.Context, id uuid.UUID) error {
	exists, err := s.flagValueRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if exists == nil {
		return errors.New("flag value not found")
	}

	err = s.flagValueRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Broadcast SSE event
	eventData := model.FlagValueEvent{
		FlagValueID:  exists.ID,
		FlagID:       exists.FlagID,
		EnvironmentID: exists.EnvID,
	}
	s.sseService.BroadcastEvent(sse.FlagValueDeleted, eventData)

	return nil
}
