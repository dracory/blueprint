package blind_index_rebuild

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestNewBlindIndexRebuildTask_InitializesFields(t *testing.T) {
	registry := testutils.Setup()

	task := NewBlindIndexRebuildTask(registry)

	if task.registry != registry {
		t.Fatalf("expected registry to be set on task")
	}

	expected := []string{BlindIndexAll, BlindIndexEmail, BlindIndexFirstName, BlindIndexLastName}
	if !reflect.DeepEqual(task.allowedIndexes, expected) {
		t.Fatalf("unexpected allowed indexes: got %v, want %v", task.allowedIndexes, expected)
	}
}

func TestBlindIndexRebuildTask_Metadata(t *testing.T) {
	registry := testutils.Setup()
	task := NewBlindIndexRebuildTask(registry)

	if got, want := task.Alias(), "BlindIndexUpdate"; got != want {
		t.Fatalf("Alias() = %q, want %q", got, want)
	}

	if got, want := task.Title(), "Blind Index Update"; got != want {
		t.Fatalf("Title() = %q, want %q", got, want)
	}

	if got, want := task.Description(), "Truncates a blind index table, and repopulates it with the current data"; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}
}

func TestBlindIndexRebuildTask_Enqueue_TaskStoreNil(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(false)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	task := NewBlindIndexRebuildTask(registry)

	if _, err := task.Enqueue(BlindIndexAll); err == nil {
		t.Fatalf("expected error when task store is nil, got nil")
	}
}

func TestBlindIndexRebuildTask_Enqueue_InvalidIndex(t *testing.T) {
	registry := testutils.Setup()
	task := NewBlindIndexRebuildTask(registry)

	if _, err := task.Enqueue("invalid"); err == nil {
		t.Fatalf("expected error when invalid index is provided, got nil")
	}
}

func TestBlindIndexRebuildTask_Handle(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithTaskStore(true),
		testutils.WithUserStore(true),
	)

	// Register task
	err := registry.GetTaskStore().TaskHandlerAdd(context.Background(), NewBlindIndexRebuildTask(registry), true)
	if err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %q", err)
	}

	// Enqueue task with valid funnel and lead
	queuedTask, err := NewBlindIndexRebuildTask(registry).Enqueue(BlindIndexAll)
	if err != nil {
		t.Fatalf("Enqueue() expected nil error, got %q", err)
	}

	// Set queued task
	task := NewBlindIndexRebuildTask(registry)
	task.SetQueuedTask(queuedTask)

	// ACT
	if ok := task.Handle(); !ok {
		t.Fatalf("Handle() expected true, got false")
	}

	details := task.QueuedTask().Details()
	if details == "" {
		t.Fatalf("Details() should not be empty")
	}

	if !strings.Contains(details, "Rebuilding email index:") {
		t.Fatalf("Details() should contain 'Rebuilding email index:' but got %q", details)
	}

	if !strings.Contains(details, "Rebuilding first name index:") {
		t.Fatalf("Details() should contain 'Rebuilding first name index:' but got %q", details)
	}

	if !strings.Contains(details, "Rebuilding last name index:") {
		t.Fatalf("Details() should contain 'Rebuilding last name index:' but got %q", details)
	}

	if !strings.Contains(details, "Index rebuilt successfully") {
		t.Fatalf("Details() should contain 'Index rebuilt successfully' but got %q", details)
	}
}
