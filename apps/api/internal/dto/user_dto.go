package dto

import (
	"time"

	"github.com/google/uuid"
)

// User request DTOs
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanumunderscore"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50,alphanumunderscore"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=1,max=50,alpha_space"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50,alpha_space"`
}

// User response DTOs
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type AuthError struct {
	Error string `json:"error"`
}
