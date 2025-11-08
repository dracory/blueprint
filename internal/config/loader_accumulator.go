package config

import "strings"

// loadAccumulator centralises validation error collection while building a
// configuration instance. Helper methods mirror the existing requireString
// and requireWhen primitives so callers stay concise.
type loadAccumulator struct {
	errs []error
}

// add appends err to the accumulator when it is non-nil.
func (a *loadAccumulator) add(err error) {
	if err != nil {
		a.errs = append(a.errs, err)
	}
}

// mustString returns the value for key via requireString, while recording any
// resulting error for later inspection.
func (a *loadAccumulator) mustString(key, context string) string {
	value, err := requireString(key, context)
	a.add(err)
	return value
}

// mustWhen delegates to requireWhen and records any error produced under the
// supplied condition.
func (a *loadAccumulator) mustWhen(condition bool, key, context, value string) {
	if err := requireWhen(condition, key, context, value); err != nil {
		a.add(err)
	}
}

// err returns a validationError wrapping all collected issues. Nil is returned
// when no errors were recorded.
func (a *loadAccumulator) err() error {
	if len(a.errs) == 0 {
		return nil
	}
	return validationError{errs: a.errs}
}

// validationError aggregates multiple missing/invalid environment errors while
// preserving the existing error semantics.
type validationError struct {
	errs []error
}

func (e validationError) Error() string {
	if len(e.errs) == 0 {
		return "config: validation failed"
	}

	var builder strings.Builder
	builder.WriteString("config: validation failed:\n")
	for i, err := range e.errs {
		if i > 0 {
			builder.WriteByte('\n')
		}
		builder.WriteString(" - ")
		builder.WriteString(err.Error())
	}
	return builder.String()
}

// Errors exposes the accumulated error slice. A defensive copy is returned to
// avoid callers mutating internal state.
func (e validationError) Errors() []error {
	return append([]error(nil), e.errs...)
}
