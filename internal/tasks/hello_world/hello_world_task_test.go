package hello_world

import (
	"reflect"
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestNewHelloWorldTask_InitializesFields(t *testing.T) {
	app := testutils.Setup()

	handler := NewHelloWorldTask(app)

	if handler == nil {
		t.Fatalf("expected handler to be non-nil")
	}

	// verify app is set via reflection since app field is unexported
	v := reflect.ValueOf(handler).Elem().FieldByName("app")
	if !v.IsValid() || v.IsNil() {
		t.Fatalf("expected app to be set on handler")
	}
}

func TestHelloWorldTask_Metadata(t *testing.T) {
	app := testutils.Setup()
	handler := NewHelloWorldTask(app)

	if got, want := handler.Alias(), "HelloWorldTask"; got != want {
		t.Fatalf("Alias() = %q, want %q", got, want)
	}

	if got, want := handler.Title(), "Hello World"; got != want {
		t.Fatalf("Title() = %q, want %q", got, want)
	}

	if got, want := handler.Description(), "Say hello world"; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}
}

func TestHelloWorldTask_Enqueue_AppNil(t *testing.T) {
	// handler with nil app should fail
	handler := &helloWorldTask{}

	if _, err := handler.Enqueue(); err == nil {
		t.Fatalf("expected error when app is nil, got nil")
	}
}

func TestHelloWorldTask_Enqueue_TaskStoreNil(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	handler := NewHelloWorldTask(app)

	if _, err := handler.Enqueue(); err == nil {
		t.Fatalf("expected error when task store is nil, got nil")
	}
}

func TestHelloWorldTask_Handle_EnqueuedTask(t *testing.T) {
	app := testutils.Setup(testutils.WithTaskStore(true))

	if app.GetTaskStore() == nil {
		t.Fatalf("expected task store to be initialized")
	}

	// Register task
	if err := app.GetTaskStore().TaskHandlerAdd(NewHelloWorldTask(app), true); err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %q", err)
	}

	// Enqueue task
	enqueueHandler := NewHelloWorldTask(app)
	queuedTask, err := enqueueHandler.Enqueue()
	if err != nil {
		t.Fatalf("Enqueue() expected nil error, got %q", err)
	}

	// Handle using queued task
	handler := NewHelloWorldTask(app)
	handler.SetQueuedTask(queuedTask)

	if ok := handler.Handle(); !ok {
		t.Fatalf("Handle() expected true, got false")
	}

	details := handler.QueuedTask().Details()
	if details == "" {
		t.Fatalf("Details() should not be empty after successful Handle")
	}

	if !strings.Contains(details, "Hello World!") {
		t.Fatalf("Details() should contain 'Hello World!' but got %q", details)
	}
}
