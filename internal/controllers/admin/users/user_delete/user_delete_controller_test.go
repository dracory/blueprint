package user_delete

import (
	"testing"

	"project/internal/testutils"
)

func TestNewUserDeleteController(t *testing.T) {
	// Test with nil registry
	controller := NewUserDeleteController(nil)
	if controller == nil {
		t.Error("NewUserDeleteController() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewUserDeleteController(registry)
	if controller == nil {
		t.Error("NewUserDeleteController() should not return nil")
	}
}
