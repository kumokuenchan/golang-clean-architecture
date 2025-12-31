package api

import "errors"

// Common API errors
var (
	ErrIDRequired = errors.New("ID is required")
	ErrInvalidID  = errors.New("Invalid ID")
)