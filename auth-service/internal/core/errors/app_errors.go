package apperrors

import "errors"

// Common application errors used across all layers.
var (
	// ErrNotFound indicates that the requested resource does not exist.
	// Used by repositories when a query returns no results.
	ErrNotFound = errors.New("not found")

	// ErrInvalidArgument indicates that the provided input is invalid.
	// Used by domain validation and request decoding.
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrConflict indicates a resource already exists or a version conflict.
	// Used when creating a resource that violates a uniqueness constraint.
	ErrConflict = errors.New("conflict")

	// ErrUnauthorized indicates that authentication is required or failed.
	// Used when credentials are missing, invalid, or expired.
	ErrUnauthorized = errors.New("unauthorized")
)
