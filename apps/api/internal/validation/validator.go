package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	v := validator.New()

	// Register custom validators
	v.RegisterValidation("alphanumunderscore", validateAlphaNumUnderscore)
	v.RegisterValidation("alpha_space", validateAlphaSpace)

	return &Validator{
		validate: v,
	}
}

// Validate validates a struct using validator tags
func (v *Validator) Validate(s interface{}) error {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	// Return the first validation error
	return ErrValidationFailed{Message: err.Error()}
}

// validateAlphaNumUnderscore validates that a string contains only alphanumeric characters and underscores
func validateAlphaNumUnderscore(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	for _, c := range field {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

// validateAlphaSpace validates that a string contains only alphabetic characters and spaces
func validateAlphaSpace(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	for _, c := range field {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == ' ' || c == '-' || c == '\'') {
			return false
		}
	}
	return true
}

// ErrValidationFailed represents a validation error
type ErrValidationFailed struct {
	Message string
}

func (e ErrValidationFailed) Error() string {
	// Clean up the error message
	msg := e.Message

	// Remove "Key: '' Error:Field validation for '' failed on the '' tag'" prefix
	if strings.Contains(msg, "Error:Field validation for") {
		parts := strings.Split(msg, "Error:Field validation for")
		if len(parts) > 1 {
			fieldParts := strings.Split(parts[1], " failed on the")
			if len(fieldParts) > 1 {
				fieldName := strings.TrimSpace(fieldParts[0])
				if fieldName[0] == '\'' {
					fieldName = fieldName[1 : len(fieldName)-1]
				}

				tagParts := strings.Split(fieldParts[1], "' tag")
				if len(tagParts) > 0 {
					tagName := strings.TrimSpace(tagParts[0])
					if tagName[0] == '\'' {
						tagName = tagName[1 : len(tagName)-1]
					}

					// Create a more user-friendly message
					switch tagName {
					case "required":
						return fmt.Sprintf("%s is required", fieldName)
					case "email":
						return fmt.Sprintf("%s must be a valid email address", fieldName)
					case "min":
						return fmt.Sprintf("%s is too short", fieldName)
					case "max":
						return fmt.Sprintf("%s is too long", fieldName)
					case "alphanumunderscore":
						return fmt.Sprintf("%s can only contain letters, numbers, underscores, and hyphens", fieldName)
					case "alpha_space":
						return fmt.Sprintf("%s can only contain letters, spaces, hyphens, and apostrophes", fieldName)
					default:
						return fmt.Sprintf("%s is invalid", fieldName)
					}
				}
			}
		}
	}

	return msg
}
