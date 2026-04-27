package blogadmin

import "errors"

// Common errors
var (
	ErrStoreRequired      = errors.New("blog store is required")
	ErrLoggerRequired     = errors.New("logger is required")
	ErrFuncLayoutRequired = errors.New("FuncLayout is required")
	ErrAuthRequired       = errors.New("authentication required")
)
