package schedules

import (
	"context"
	"testing"
	"time"

	"project/internal/testutils"
)

func TestStartAsync(t *testing.T) {
	// Test with nil registry
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	StartAsync(ctx, nil)
	// Should not panic

	// Test with valid registry
	registry := testutils.Setup()

	// Start the scheduler in a goroutine with a short timeout
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	done := make(chan bool)
	go func() {
		StartAsync(ctx, registry)
		done <- true
	}()

	// Wait for it to complete or timeout
	select {
	case <-done:
		// Expected
	case <-time.After(2 * time.Second):
		t.Error("StartAsync did not complete within timeout")
	}
}

func TestStartAsync_NilRegistry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// This should not panic and should log an error
	StartAsync(ctx, nil)
	// Expected behavior - function returns early
}

func TestNewScheduler(t *testing.T) {
	// Test with nil registry
	scheduler := newScheduler(nil)
	if scheduler == nil {
		t.Error("newScheduler() should not return nil even with nil registry")
	}

	// Clean up
	scheduler.Stop()
	scheduler.Clear()

	// Test with valid registry
	registry := testutils.Setup()
	scheduler = newScheduler(registry)
	if scheduler == nil {
		t.Error("newScheduler() should not return nil with valid registry")
	}

	// Clean up
	scheduler.Stop()
	scheduler.Clear()
}

func TestStartAsync_ContextCancellation(t *testing.T) {
	registry := testutils.Setup()

	// Create a context that cancels immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Should not panic and should return immediately
	StartAsync(ctx, registry)
}

func TestScheduleBlindIndexRebuildTask(t *testing.T) {
	// Test with nil registry - should not panic
	scheduleBlindIndexRebuildTask(nil)
	// Should not panic
}

func TestScheduleCleanUpTask(t *testing.T) {
	// Test with nil registry
	scheduleCleanUpTask(nil)
	// Should not panic

	// Test with valid registry
	registry := testutils.Setup(testutils.WithTaskStore(true))
	scheduleCleanUpTask(registry)
	// Should not panic
}

func TestScheduleStatsVisitorEnhanceTask(t *testing.T) {
	// Test with nil registry
	scheduleStatsVisitorEnhanceTask(nil)
	// Should not panic

	// Test with valid registry
	registry := testutils.Setup(testutils.WithTaskStore(true))
	scheduleStatsVisitorEnhanceTask(registry)
	// Should not panic
}

func TestQueueClearJob(t *testing.T) {
	// Test with nil registry
	queueClearJob(nil)
	// Should not panic

	// Test with valid registry
	registry := testutils.Setup(testutils.WithTaskStore(true))
	queueClearJob(registry)
	// Should not panic
}
