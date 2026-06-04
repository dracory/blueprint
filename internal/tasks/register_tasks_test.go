package tasks

import (
	"testing"

	"project/internal/testutils"
)

func TestRegisterTasks(t *testing.T) {
	// Test with nil app - this will panic
	defer func() {
		if r := recover(); r != nil {
			t.Log("RegisterTasks() panicked as expected with nil app")
		}
	}()
	RegisterTasks(nil)
}

func TestRegisterTasks_WithoutTaskStore(t *testing.T) {
	// Test with app without task store
	app := testutils.Setup()
	RegisterTasks(app)
	// Should not panic (returns early when task store is nil)
}

func TestRegisterTasks_WithTaskStore(t *testing.T) {
	// Test with app with task store
	app := testutils.Setup(testutils.WithTaskStore(true))
	RegisterTasks(app)
	// Should not panic
}
