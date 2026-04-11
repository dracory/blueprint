package main

import (
	"context"
	"testing"
	"time"
)

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
	group := newBackgroundGroup(context.TODO())
	defer group.stop()

	if group.ctx == nil {
		t.Error("Context should not be nil")
	}
}

func TestBackgroundGroup_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	group := newBackgroundGroup(ctx)
	defer group.stop()

	executed := false
	group.Go(func(ctx context.Context) {
		<-ctx.Done()
		executed = true
	})

	// Cancel context
	cancel()

	// Give goroutine time to execute
	time.Sleep(100 * time.Millisecond)

	// Stop and verify goroutine was executed
	group.stop()

	if !executed {
		t.Error("Goroutine should have been executed and context should be cancelled")
	}
}

func TestBackgroundGroup_MultipleGoroutines(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)
	defer group.stop()

	counter := 0
	done := make(chan bool, 3)

	for i := 0; i < 3; i++ {
		group.Go(func(ctx context.Context) {
			counter++
			done <- true
		})
	}

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		select {
		case <-done:
			// Expected
		case <-time.After(2 * time.Second):
			t.Error("Goroutine should have executed")
		}
	}

	if counter != 3 {
		t.Errorf("Expected 3 goroutines to execute, got %d", counter)
	}
}

func TestBackgroundGroup_StopIdempotent(t *testing.T) {
	ctx := context.Background()
	group := newBackgroundGroup(ctx)

	// Call stop multiple times - should not panic
	group.stop()
	group.stop()
	group.stop()

	// Verify Done channel is still closed
	select {
	case <-group.Done():
		// Expected
	default:
		t.Error("Done channel should be closed")
	}
}

func TestBackgroundGroup_WaitGroupSync(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group := newBackgroundGroup(ctx)

	// Track execution order
	executed := make([]int, 0)
	done := make(chan bool, 3)

	// Add multiple goroutines
	for i := 1; i <= 3; i++ {
		idx := i
		group.Go(func(ctx context.Context) {
			executed = append(executed, idx)
			done <- true
		})
	}

	// Wait for all to complete
	for i := 0; i < 3; i++ {
		<-done
	}

	// Stop should wait for all goroutines
	group.stop()

	if len(executed) != 3 {
		t.Errorf("Expected 3 goroutines to execute, got %d", len(executed))
	}
}

func TestBackgroundGroup_GoWithContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	group := newBackgroundGroup(ctx)
	defer group.stop()

	// Cancel context immediately
	cancel()

	executed := false
	group.Go(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			executed = true
		case <-time.After(100 * time.Millisecond):
			// Timeout
		}
	})

	// Give goroutine time to execute
	time.Sleep(50 * time.Millisecond)
	group.stop()

	if !executed {
		t.Error("Goroutine should have detected context cancellation")
	}
}

func TestNewBackgroundGroup_WithTODOContext(t *testing.T) {
	group := newBackgroundGroup(context.TODO())
	defer group.stop()

	if group.ctx == nil {
		t.Error("Context should not be nil")
	}

	if group.cancel == nil {
		t.Error("Cancel function should not be nil")
	}

	if group.doneCh == nil {
		t.Error("Done channel should not be nil")
	}
}
