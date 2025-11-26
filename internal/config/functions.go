package config

import (
	"fmt"
	"strings"

	"github.com/dracory/env"
)

// MissingEnvError describes an unset required environment
// variable with optional context to aid debugging.
type MissingEnvError struct {
	Key     string
	Context string
}

// Error returns the formatted error describing the missing
// environment variable and optional context.
func (e MissingEnvError) Error() string {
	if strings.TrimSpace(e.Context) == "" {
		return fmt.Sprintf("config: required env %q is missing", e.Key)
	}

	return fmt.Sprintf(
		"config: required env %q is missing: %s",
		e.Key,
		e.Context,
	)
}

// requireString trims and retrieves the environment value for the
// provided key, returning a typed MissingEnvError when the value is absent.
func requireString(key, context string) (string, error) {
	value := strings.TrimSpace(env.GetString(key))

	if err := ensureRequired(value, key, context); err != nil {
		return "", err
	}

	return value, nil
}

// requireWhen validates that the given value is present when the
// supplied condition evaluates to true.
func requireWhen(condition bool, key, context, value string) error {
	if !condition {
		return nil
	}

	return ensureRequired(value, key, context)
}

// ensureRequired returns a MissingEnvError when the supplied value
// is blank after trimming whitespace.
func ensureRequired(value, key, context string) error {
	if strings.TrimSpace(value) != "" {
		return nil
	}

	return MissingEnvError{Key: key, Context: context}
}
