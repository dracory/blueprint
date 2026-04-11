package ai_tools

import (
	"testing"

	"project/internal/testutils"
)

// TestNewAiToolsController verifies controller can be created
func TestNewAiToolsController(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	controller := NewAiToolsController(app)

	if controller == nil {
		t.Error("NewAiToolsController() returned nil")
	}
}

// TestAiToolsControllerRegistry verifies controller has registry
func TestAiToolsControllerRegistry(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	controller := NewAiToolsController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestAiToolsControllerHandlerExists verifies Handler method exists
func TestAiToolsControllerHandlerExists(t *testing.T) {
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	controller := NewAiToolsController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}
