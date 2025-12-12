package validation

import (
	"api/internal/errors"
	"regexp"
	"strings"
)

// Validator provides validation functions
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateUsername validates username
func (v *Validator) ValidateUsername(username string) error {
	if strings.TrimSpace(username) == "" {
		return errors.ErrMissingRequiredField
	}
	
	if len(username) < 3 || len(username) > 50 {
		return errors.ErrInvalidField
	}
	
	// Allow alphanumeric, underscore, and hyphen
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	if !matched {
		return errors.NewAppError(400, "Username can only contain letters, numbers, underscores, and hyphens")
	}
	
	return nil
}

// ValidateEmail validates email
func (v *Validator) ValidateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return errors.ErrMissingRequiredField
	}
	
	// Basic email validation
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if !matched {
		return errors.ErrInvalidField
	}
	
	return nil
}

// ValidatePassword validates password
func (v *Validator) ValidatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return errors.ErrMissingRequiredField
	}
	
	if len(password) < 8 {
		return errors.ErrInvalidField
	}
	
	// Basic password strength requirements
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	
	if !hasUpper || !hasLower || !hasNumber {
		return errors.NewAppError(400, "Password must contain at least one uppercase letter, one lowercase letter, and one number")
	}
	
	return nil
}

// ValidateName validates first/last name
func (v *Validator) ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.ErrMissingRequiredField
	}
	
	if len(name) < 1 || len(name) > 50 {
		return errors.ErrInvalidField
	}
	
	// Allow letters, spaces, hyphens, and apostrophes
	matched, _ := regexp.MatchString(`^[a-zA-Z\s'-]+$`, name)
	if !matched {
		return errors.NewAppError(400, "Name can only contain letters, spaces, hyphens, and apostrophes")
	}
	
	return nil
}

// ValidateRequired checks if field is required and not empty
func (v *Validator) ValidateRequired(field, fieldName string) error {
	if strings.TrimSpace(field) == "" {
		return errors.NewAppError(400, fieldName+" is required")
	}
	return nil
}

// ValidateLength checks field length
func (v *Validator) ValidateLength(field string, fieldName string, min, max int) error {
	if len(field) < min || len(field) > max {
		return errors.NewAppError(400, fieldName+" must be between "+string(min)+" and "+string(max)+" characters")
	}
	return nil
}

// ValidateUUID validates UUID string
func (v *Validator) ValidateUUID(uuid string) error {
	if strings.TrimSpace(uuid) == "" {
		return errors.ErrMissingRequiredField
	}
	
	// Basic UUID validation (format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)
	matched, _ := regexp.MatchString(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, uuid)
	if !matched {
		return errors.ErrInvalidField
	}
	
	return nil
}
