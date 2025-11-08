package main

import (
	"context"
	"os"
	"testing"

	"project/internal/testutils"
)

func TestCloseResources(t *testing.T) {
	// Should not panic when db handle is nil
	closeResourcesDB(nil)
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
	// Initialize minimal stores for background processes
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Should not panic
	startBackgroundProcesses(ctx, group, app)

	if app.GetCacheStore() == nil {
		t.Errorf("Cache store should not be nil after starting background processes")
	}
	if app.GetSessionStore() == nil {
		t.Errorf("Session store should not be nil after starting background processes")
	}
}
