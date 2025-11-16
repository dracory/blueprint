package clean_up

import (
	"reflect"
	"testing"

	"project/internal/testutils"
)

func TestNewCleanUpTask_InitializesFields(t *testing.T) {
	app := testutils.Setup()

	handlerIface := NewCleanUpTask(app)
	handler, ok := handlerIface.(*cleanUpTask)
	if !ok {
		t.Fatalf("expected *cleanUpTask, got different type")
	}

	if handler == nil {
		t.Fatalf("expected handler to be non-nil")
	}

	// verify app is set via reflection since app field is unexported
	v := reflect.ValueOf(handler).Elem().FieldByName("app")
	if !v.IsValid() || v.IsNil() {
		t.Fatalf("expected app to be set on handler")
	}
}

func TestCleanUpTask_Metadata(t *testing.T) {
	app := testutils.Setup()
	handler := NewCleanUpTask(app)

	if got, want := handler.Alias(), "CleanUpTask"; got != want {
		t.Fatalf("Alias() = %q, want %q", got, want)
	}

	if got, want := handler.Title(), "Clean Up"; got != want {
		t.Fatalf("Title() = %q, want %q", got, want)
	}

	if got, want := handler.Description(), "Clean up the database"; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}
}

func TestCleanUpTask_Enqueue_AppNil(t *testing.T) {
	// construct handler with nil app to trigger app check
	handler := &cleanUpTask{}

	if _, err := handler.Enqueue(); err == nil {
		t.Fatalf("expected error when app is nil, got nil")
	}
}

func TestCleanUpTask_Enqueue_TaskStoreNil(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	handlerIface := NewCleanUpTask(app)
	handler, ok := handlerIface.(*cleanUpTask)
	if !ok {
		t.Fatalf("expected *cleanUpTask, got different type")
	}

	if _, err := handler.Enqueue(); err == nil {
		t.Fatalf("expected error when task store is nil, got nil")
	}
}

func TestCleanUpTask_Handle_TaskStoreNil(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	handler := NewCleanUpTask(app)

	if ok := handler.Handle(); !ok {
		t.Fatalf("Handle() expected true when task store is nil, got false")
	}
}

func TestCleanUpTask_Handle_Success(t *testing.T) {
	app := testutils.Setup(testutils.WithTaskStore(true))

	if app.GetTaskStore() == nil {
		t.Fatalf("expected task store to be initialized")
	}

	handler := NewCleanUpTask(app)

	if ok := handler.Handle(); !ok {
		t.Fatalf("Handle() expected true, got false")
	}
}
