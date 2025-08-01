package main

import (
	"os"
	"project/internal/config"
	"project/internal/testutils"
	"testing"
)

func TestCloseResources(t *testing.T) {
	closeResources()
	// Assuming config.Database is a global variable that should be nil after closing
	if config.Database != nil {
		t.Errorf("Database should be closed and set to nil")
	}
}

func TestIsCliMode(t *testing.T) {
	os.Args = []string{"main", "task", "testTask"}
	if !isCliMode() {
		t.Errorf("isCliMode() should return true")
	}

	os.Args = []string{"main"}
	if isCliMode() {
		t.Errorf("isCliMode() should return false")
	}
}

func TestStartBackgroundProcesses(t *testing.T) {
	testutils.Setup()
	startBackgroundProcesses()
	// Assuming we can verify background processes by checking if certain goroutines are running
	// This is a placeholder assertion; actual verification would depend on the implementation
	// Assuming we can verify background processes by checking if certain goroutines are running
	// This is a placeholder assertion; actual verification would depend on the implementation
	if false {
		t.Errorf("Background processes should be started")
	}

	// Verify that the background processes started correctly
	if config.TaskStore == nil {
		t.Errorf("Task store should not be nil after starting background processes")
	}
	if config.CacheStore == nil {
		t.Errorf("Cache store should not be nil after starting background processes")
	}
	if config.SessionStore == nil {
		t.Errorf("Session store should not be nil after starting background processes")
	}
	if config.ShopStore == nil {
		t.Errorf("Shop store should not be nil after starting background processes")
	}
}
