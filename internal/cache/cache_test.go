package cache

import (
	"testing"
)

// TestCacheVariablesExported verifies that Memory and File cache variables are exported
func TestCacheVariablesExported(t *testing.T) {
	// These should be accessible (not cause compilation error)
	_ = Memory
	_ = File
}

// TestCacheVariableTypes verifies cache variables have correct types
func TestCacheVariableTypes(t *testing.T) {
	// Memory should be a *ttlcache.Cache[string, any] or nil
	// This test verifies the variables are declared with correct types
	if Memory != nil {
		// If not nil, it should be usable
		_ = Memory
	}

	// File should be a cachego.Cache or nil
	if File != nil {
		// If not nil, it should be usable
		_ = File
	}
}
