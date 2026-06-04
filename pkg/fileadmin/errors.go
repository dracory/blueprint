package fileadmin

import "errors"

var (
	// ErrRegistryRequired is returned when the app is nil
	ErrRegistryRequired = errors.New("app is required")
)
