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

	// Register task (this is needed for the task definition to be available)
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

func TestBlindIndexRebuildTask_Enqueue_EmailIndex(t *testing.T) {
	registry := testutils.Setup(testutils.WithTaskStore(true))
	task := NewBlindIndexRebuildTask(registry)

	// Register task first
	err := registry.GetTaskStore().TaskHandlerAdd(context.Background(), task, true)
	if err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %v", err)
	}

	_, err = task.Enqueue(BlindIndexEmail)
	if err != nil {
		t.Fatalf("Enqueue() with email index expected nil error, got %v", err)
	}
}

func TestBlindIndexRebuildTask_Enqueue_FirstNameIndex(t *testing.T) {
	registry := testutils.Setup(testutils.WithTaskStore(true))
	task := NewBlindIndexRebuildTask(registry)

	// Register task first
	err := registry.GetTaskStore().TaskHandlerAdd(context.Background(), task, true)
	if err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %v", err)
	}

	_, err = task.Enqueue(BlindIndexFirstName)
	if err != nil {
		t.Fatalf("Enqueue() with first_name index expected nil error, got %v", err)
	}
}

func TestBlindIndexRebuildTask_Enqueue_LastNameIndex(t *testing.T) {
	registry := testutils.Setup(testutils.WithTaskStore(true))
	task := NewBlindIndexRebuildTask(registry)

	// Register task first
	err := registry.GetTaskStore().TaskHandlerAdd(context.Background(), task, true)
	if err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %v", err)
	}

	_, err = task.Enqueue(BlindIndexLastName)
	if err != nil {
		t.Fatalf("Enqueue() with last_name index expected nil error, got %v", err)
	}
}

func TestBlindIndexRebuildTask_SetQueuedTask(t *testing.T) {
	registry := testutils.Setup()
	task := NewBlindIndexRebuildTask(registry)

	// SetQueuedTask should not panic even with nil
	task.SetQueuedTask(nil)

	if task.QueuedTask() != nil {
		t.Error("QueuedTask() should return nil after setting nil")
	}
}

func TestBlindIndexRebuildTask_Constants(t *testing.T) {
	if BlindIndexAll != "all" {
		t.Errorf("BlindIndexAll should be 'all', got %q", BlindIndexAll)
	}
	if BlindIndexEmail != "email" {
		t.Errorf("BlindIndexEmail should be 'email', got %q", BlindIndexEmail)
	}
	if BlindIndexFirstName != "first_name" {
		t.Errorf("BlindIndexFirstName should be 'first_name', got %q", BlindIndexFirstName)
	}
	if BlindIndexLastName != "last_name" {
		t.Errorf("BlindIndexLastName should be 'last_name', got %q", BlindIndexLastName)
	}
}
