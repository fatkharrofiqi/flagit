package service

import (
	"api/internal/dto"
	"api/internal/middleware"
	"api/internal/model"
	"api/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

// authService implements AuthService interface
type authService struct {
	userRepo repository.UserRepository
	jwtSecret string
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user and returns JWT token
func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Check if user already exists by username
	existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingEmail, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &model.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      string(middleware.RoleDeveloper), // Default role
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Active:    true,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID.String(), user.Username, user.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// Login authenticates a user and returns JWT token
func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	// Check if user is active
	if !user.Active {
		return nil, errors.New("user account is inactive")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID.String(), user.Username, user.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *authService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
