package aipostgenerator

import (
	"testing"

	"project/internal/testutils"
)

func TestNewAiPostGeneratorController(t *testing.T) {
	// Test with nil registry
	controller := NewAiPostGeneratorController(nil)
	if controller == nil {
		t.Error("NewAiPostGeneratorController() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewAiPostGeneratorController(registry)
	if controller == nil {
		t.Error("NewAiPostGeneratorController() should not return nil")
	}
	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}

func TestConstants(t *testing.T) {
	if ACTION_GENERATE_POST == "" {
		t.Error("ACTION_GENERATE_POST should not be empty")
	}
}
