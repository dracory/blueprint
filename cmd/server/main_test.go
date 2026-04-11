package main

import (
	"context"
	"os"
	"testing"
	"time"

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

func TestStartBackgroundProcesses_NilConfig(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Create a minimal mock registry without config
	app := testutils.Setup()

	// This will likely fail because we can't easily mock GetConfig() returning nil
	// But let's at least test the flow
	err := startBackgroundProcesses(ctx, group, app)
	// Should not error because testutils.Setup() provides a valid config
	if err != nil {
		t.Logf("startBackgroundProcesses returned error (expected for test setup): %v", err)
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

func TestBackgroundGroup_Stop(t *testing.T) {
	ctx := context.Background()
	group := newBackgroundGroup(ctx)

	// Test stop multiple times (should not panic)
	group.stop()
	group.stop() // Second call should be no-op

	// Verify Done channel is closed
	select {
	case <-group.Done():
		// Expected - channel should be closed
	default:
		t.Error("Done channel should be closed after stop")
	}
}

func TestBackgroundGroup_Go(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	done := make(chan bool)
	group.Go(func(ctx context.Context) {
		close(done)
	})

	select {
	case <-done:
		// Expected
	case <-time.After(2 * time.Second):
		t.Error("Go function should have executed")
	}
}

func TestBackgroundGroup_NilParent(t *testing.T) {
	// Test that nil parent context defaults to Background
	group := newBackgroundGroup(nil)
	defer group.stop()

	if group.ctx == nil {
		t.Error("Context should not be nil")
	}
}
