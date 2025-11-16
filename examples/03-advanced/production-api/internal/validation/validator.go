// Package validation provides input validation utilities for the User Management API.
package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ocrosby/go-lab/projects/api/internal/domain"
)

// Validator provides centralized validation logic
type Validator struct {
	emailRegex *regexp.Regexp
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

// ValidateNonEmpty validates that a field is not empty
func (v *Validator) ValidateNonEmpty(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// ValidateEmail validates email format
func (v *Validator) ValidateEmail(email string) error {
	if err := v.ValidateNonEmpty(email, "email"); err != nil {
		return err
	}

	if !v.emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateUserCreation validates user creation parameters
func (v *Validator) ValidateUserCreation(email, name string) error {
	if err := v.ValidateEmail(email); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := v.ValidateNonEmpty(name, "name"); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// ValidateUserUpdate validates user update parameters
func (v *Validator) ValidateUserUpdate(id, name string) error {
	if err := v.ValidateNonEmpty(id, "user ID"); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := v.ValidateNonEmpty(name, "name"); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// ValidateUserID validates user ID parameter
func (v *Validator) ValidateUserID(id string) error {
	if err := v.ValidateNonEmpty(id, "user ID"); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// ValidationError wraps validation errors to domain errors
func (v *Validator) WrapValidationError(err error) error {
	if err == nil {
		return nil
	}

	// Return domain error for consistency with the rest of the application
	return domain.ErrInvalidInput
}
