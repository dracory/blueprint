package blind_index_rebuild

import (
	"reflect"
	"testing"

	"project/internal/testutils"
)

func TestNewBlindIndexRebuildTask_InitializesFields(t *testing.T) {
	app := testutils.Setup()

	task := NewBlindIndexRebuildTask(app)

	if task.app != app {
		t.Fatalf("expected app to be set on task")
	}

	expected := []string{BlindIndexAll, BlindIndexEmail, BlindIndexFirstName, BlindIndexLastName}
	if !reflect.DeepEqual(task.allowedIndexes, expected) {
		t.Fatalf("unexpected allowed indexes: got %v, want %v", task.allowedIndexes, expected)
	}
}

func TestBlindIndexRebuildTask_Metadata(t *testing.T) {
	app := testutils.Setup()
	task := NewBlindIndexRebuildTask(app)

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
	app := testutils.Setup(testutils.WithCfg(cfg))

	task := NewBlindIndexRebuildTask(app)

	if _, err := task.Enqueue(BlindIndexAll); err == nil {
		t.Fatalf("expected error when task store is nil, got nil")
	}
}

func TestBlindIndexRebuildTask_Enqueue_InvalidIndex(t *testing.T) {
	app := testutils.Setup()
	task := NewBlindIndexRebuildTask(app)

	if _, err := task.Enqueue("invalid"); err == nil {
		t.Fatalf("expected error when invalid index is provided, got nil")
	}

}
