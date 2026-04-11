package user_create

import (
	"testing"

	"project/internal/testutils"
)

func TestNewUserCreateController(t *testing.T) {
	// Test with nil registry
	controller := NewUserCreateController(nil)
	if controller == nil {
		t.Error("NewUserCreateController() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewUserCreateController(registry)
	if controller == nil {
		t.Error("NewUserCreateController() should not return nil")
	}
}
