package tasks

import (
	"testing"

	"project/internal/testutils"
)

func TestRegisterTasks(t *testing.T) {
	// Test with nil registry - this will panic
	defer func() {
		if r := recover(); r != nil {
			t.Log("RegisterTasks() panicked as expected with nil registry")
		}
	}()
	RegisterTasks(nil)
}

func TestRegisterTasks_WithoutTaskStore(t *testing.T) {
	// Test with registry without task store
	registry := testutils.Setup()
	RegisterTasks(registry)
	// Should not panic (returns early when task store is nil)
}

func TestRegisterTasks_WithTaskStore(t *testing.T) {
	// Test with registry with task store
	registry := testutils.Setup(testutils.WithTaskStore(true))
	RegisterTasks(registry)
	// Should not panic
}
