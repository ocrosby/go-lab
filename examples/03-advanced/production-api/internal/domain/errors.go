package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInternalError     = errors.New("internal error")
)

// Error checking functions for wrapped errors

// IsUserNotFoundError checks if the error is or wraps ErrUserNotFound
func IsUserNotFoundError(err error) bool {
	return errors.Is(err, ErrUserNotFound)
}

// IsUserAlreadyExistsError checks if the error is or wraps ErrUserAlreadyExists
func IsUserAlreadyExistsError(err error) bool {
	return errors.Is(err, ErrUserAlreadyExists)
}

// IsInvalidInputError checks if the error is or wraps ErrInvalidInput
func IsInvalidInputError(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsInternalError checks if the error is or wraps ErrInternalError
func IsInternalError(err error) bool {
	return errors.Is(err, ErrInternalError)
}
