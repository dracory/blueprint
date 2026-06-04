package shopadmin

import "errors"

var (
	// ErrRegistryRequired is returned when app is not provided
	ErrRegistryRequired = errors.New("app is required")
)
