package admin

import (
	"testing"

	"project/internal/testutils"
)

func TestTaskRoutes(t *testing.T) {
	// Test with nil registry - expect panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("TaskRoutes(nil) should panic")
		}
	}()
	TaskRoutes(nil)
}

func TestTaskRoutesWithRegistry(t *testing.T) {
	// Test with valid registry
	registry := testutils.Setup()
	routes := TaskRoutes(registry)
	if routes == nil {
		t.Error("TaskRoutes() should not return nil")
	}
	if len(routes) != 2 {
		t.Errorf("TaskRoutes(registry) should return 2 routes, got %d", len(routes))
	}
}
