package admin

import (
	"testing"

	"project/internal/testutils"
)

func TestTaskRoutes(t *testing.T) {
	// Test with nil registry - should return error
	routes, err := TaskRoutes(nil)
	if err == nil {
		t.Error("TaskRoutes(nil) should return an error")
	}
	if routes != nil {
		t.Error("TaskRoutes(nil) should return nil routes")
	}
}

func TestTaskRoutesWithRegistry(t *testing.T) {
	// Test with valid registry
	registry := testutils.Setup()
	routes, err := TaskRoutes(registry)
	if err != nil {
		t.Errorf("TaskRoutes(registry) should not return an error: %v", err)
	}
	if routes == nil {
		t.Error("TaskRoutes() should not return nil")
	}
	if len(routes) != 2 {
		t.Errorf("TaskRoutes(registry) should return 2 routes, got %d", len(routes))
	}
}
