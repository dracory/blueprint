package fileadmin

import "errors"

var (
	// ErrRegistryRequired is returned when the registry is nil
	ErrRegistryRequired = errors.New("registry is required")
)
