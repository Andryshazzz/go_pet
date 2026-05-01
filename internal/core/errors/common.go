package core_errors

import "errors"

var (
	ErrNotFounde       = errors.New("Not founde")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrConflict        = errors.New("conflict")
)
