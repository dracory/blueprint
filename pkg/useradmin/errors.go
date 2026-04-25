package useradmin

import "errors"

var (
	// ErrRegistryRequired is returned when registry is not provided
	ErrRegistryRequired = errors.New("registry is required")
)
