package shared

import (
	"testing"
)

// TestConstants verifies constants are defined
func TestConstants(t *testing.T) {
	t.Parallel()
	// This test ensures the constants can be accessed
	_ = CONTROLLER_LOG_MANAGER
}
