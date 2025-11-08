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
	userStore, sessionStore, cacheStore, cleanup := testutils.SetupTestAuth(t)
	defer cleanup()

	// Build application via testutils to ensure cfg, DB, and stores are initialized
	application := testutils.Setup()

	// Inject required stores
	application.SetUserStore(userStore)
	application.SetSessionStore(sessionStore)
	application.SetCacheStore(cacheStore)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Should not panic
	startBackgroundProcesses(ctx, group, application)

	if application.GetCacheStore() == nil {
		t.Errorf("Cache store should not be nil after starting background processes")
	}
	if application.GetSessionStore() == nil {
		t.Errorf("Session store should not be nil after starting background processes")
	}
}
