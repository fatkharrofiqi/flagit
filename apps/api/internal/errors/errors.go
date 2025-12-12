package errors

import (
	"errors"
	"net/http"
)

// AppError represents a custom application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Predefined error types with their HTTP status codes
var (
	// Validation errors
	ErrInvalidRequestBody = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid request body",
	}
	ErrInvalidJSON = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid JSON format",
	}
	ErrMissingRequiredField = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Missing required field",
	}
	ErrInvalidField = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid field value",
	}

	// Authentication errors
	ErrAuthenticationRequired = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Authentication required",
	}
	ErrInvalidCredentials = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid username or password",
	}
	ErrInvalidToken = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid or expired token",
	}
	ErrAccountInactive = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "User account is inactive",
	}

	// Authorization errors
	ErrInsufficientPermissions = &AppError{
		Code:    http.StatusForbidden,
		Message: "Insufficient permissions",
	}
	ErrAccessDenied = &AppError{
		Code:    http.StatusForbidden,
		Message: "Access denied",
	}
	ErrInsufficientRole = &AppError{
		Code:    http.StatusForbidden,
		Message: "Insufficient role",
	}

	// Not found errors
	ErrUserNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "User not found",
	}
	ErrProjectNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "Project not found",
	}
	ErrEnvironmentNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "Environment not found",
	}
	ErrFlagNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "Flag not found",
	}
	ErrFlagValueNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "Flag value not found",
	}

	// Conflict errors
	ErrUsernameExists = &AppError{
		Code:    http.StatusConflict,
		Message: "Username already exists",
	}
	ErrEmailExists = &AppError{
		Code:    http.StatusConflict,
		Message: "Email already exists",
	}
	ErrProjectKeyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "Project key already exists",
	}

	// Business logic errors
	ErrCannotDeleteProjectWithEnvironments = &AppError{
		Code:    http.StatusConflict,
		Message: "Cannot delete project with existing environments",
	}
	ErrCannotDeleteEnvironmentWithFlags = &AppError{
		Code:    http.StatusConflict,
		Message: "Cannot delete environment with existing flags",
	}
	ErrCannotDeleteFlagWithValue = &AppError{
		Code:    http.StatusConflict,
		Message: "Cannot delete flag with existing values",
	}

	// Server errors
	ErrInternalServer = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	}
	ErrDatabaseConnection = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Database connection error",
	}
)

// NewAppError creates a new application error
func NewAppError(code int, message string, details ...string) *AppError {
	err := &AppError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}
