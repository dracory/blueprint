package blogadmin

import (
	"testing"
)

// TestNewBlogAdmin_NilStore tests the New constructor with nil store
func TestNewBlogAdmin_NilStore(t *testing.T) {
	t.Parallel()

	opts := AdminOptions{
		AdminHomeURL: "/admin",
	}

	admin, err := New(opts)
	if err == nil {
		t.Error("Expected error when store is nil")
	}
	if err != ErrStoreRequired {
		t.Errorf("Expected ErrStoreRequired, got: %v", err)
	}
	if admin != nil {
		t.Error("Expected admin to be nil when store is nil")
	}
}

// TestErrStoreRequired tests the error constant
func TestErrStoreRequired(t *testing.T) {
	t.Parallel()

	if ErrStoreRequired == nil {
		t.Error("ErrStoreRequired should not be nil")
	}
	if ErrStoreRequired.Error() != "blog store is required" {
		t.Errorf("ErrStoreRequired message incorrect: %s", ErrStoreRequired.Error())
	}
}

// TestErrLoggerRequired tests the error constant
func TestErrLoggerRequired(t *testing.T) {
	t.Parallel()

	if ErrLoggerRequired == nil {
		t.Error("ErrLoggerRequired should not be nil")
	}
	if ErrLoggerRequired.Error() != "logger is required" {
		t.Errorf("ErrLoggerRequired message incorrect: %s", ErrLoggerRequired.Error())
	}
}

// TestErrFuncLayoutRequired tests the error constant
func TestErrFuncLayoutRequired(t *testing.T) {
	t.Parallel()

	if ErrFuncLayoutRequired == nil {
		t.Error("ErrFuncLayoutRequired should not be nil")
	}
	if ErrFuncLayoutRequired.Error() != "FuncLayout is required" {
		t.Errorf("ErrFuncLayoutRequired message incorrect: %s", ErrFuncLayoutRequired.Error())
	}
}

// TestErrAuthRequired tests the error constant
func TestErrAuthRequired(t *testing.T) {
	t.Parallel()

	if ErrAuthRequired == nil {
		t.Error("ErrAuthRequired should not be nil")
	}
	if ErrAuthRequired.Error() != "authentication required" {
		t.Errorf("ErrAuthRequired message incorrect: %s", ErrAuthRequired.Error())
	}
}
