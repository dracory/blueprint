package main

import (
	"context"
	"os"
	"testing"

	"project/internal/testutils"
)

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
		testutils.WithTaskStore(true),
		testutils.WithUserStore(true),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Should not panic
	if err := startBackgroundProcesses(ctx, group, app); err != nil {
		t.Fatalf("startBackgroundProcesses returned error: %v", err)
	}

	if app.GetCacheStore() == nil {
		t.Errorf("Cache store should not be nil after starting background processes")
	}
	if app.GetSessionStore() == nil {
		t.Errorf("Session store should not be nil after starting background processes")
	}
}

func TestStartBackgroundProcesses_NilRegistry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	err := startBackgroundProcesses(ctx, group, nil)
	if err == nil {
		t.Error("startBackgroundProcesses with nil registry should return error")
	}
}

func TestBackgroundGroup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Test that Done() channel works
	select {
	case <-group.Done():
		t.Error("Background group should not be done immediately")
	default:
		// Expected
	}
}
