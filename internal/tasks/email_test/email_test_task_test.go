package email_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"project/internal/emails"
	"project/internal/testutils"

	"github.com/dracory/test"
)

func TestNewEmailTestTask_InitializesFields(t *testing.T) {
	registry := testutils.Setup()

	handler, ok := NewEmailTestTask(registry).(*emailTestTask)
	if !ok {
		t.Fatalf("expected *emailTestTask, got different type")
	}

	if handler == nil {
		t.Fatalf("expected task to be non-nil")
	}

	// verify app is set via reflection since app field is unexported
	v := reflect.ValueOf(handler).Elem().FieldByName("registry")
	if !v.IsValid() || v.IsNil() {
		t.Fatalf("expected registry to be set on task")
	}
}

func TestEmailTestTask_Metadata(t *testing.T) {
	registry := testutils.Setup()
	task := NewEmailTestTask(registry)

	if got, want := task.Alias(), "EmailTestTask"; got != want {
		t.Fatalf("Alias() = %q, want %q", got, want)
	}

	if got, want := task.Title(), "Email Test"; got != want {
		t.Fatalf("Title() = %q, want %q", got, want)
	}

	if got, want := task.Description(), "Sends a notification email to the system administrator"; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}
}

func TestEmailTestTask_Enqueue_TaskStoreNil(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(false)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	handler, ok := NewEmailTestTask(registry).(*emailTestTask)
	if !ok {
		t.Fatalf("expected *emailTestTask, got different type")
	}

	if _, err := handler.Enqueue("test@example.com", "<p>test</p>"); err == nil {
		t.Fatalf("expected error when task store is nil, got nil")
	}
}

func TestEmailTestTask_Enqueue_InvalidParams(t *testing.T) {
	registry := testutils.Setup(testutils.WithTaskStore(true))
	handler, ok := NewEmailTestTask(registry).(*emailTestTask)
	if !ok {
		t.Fatalf("expected *emailTestTask, got different type")
	}

	if _, err := handler.Enqueue("", "<p>test</p>"); err == nil {
		t.Fatalf("expected error when 'to' is empty, got nil")
	}

	if _, err := handler.Enqueue("test@example.com", ""); err == nil {
		t.Fatalf("expected error when 'html' is empty, got nil")
	}
}

func TestEmailTestTask_Handle_SendEmail(t *testing.T) {
	// configure mock SMTP server
	server, _, cleanup := test.SetupMailServer(t)
	defer cleanup()

	cfg := testutils.DefaultConf()
	cfg.SetMailDriver("smtp")
	cfg.SetMailHost("127.0.0.1")
	cfg.SetMailPort(server.PortNumber)
	cfg.SetMailUsername("")
	cfg.SetMailPassword("")
	cfg.SetTaskStoreUsed(true)

	registry := testutils.Setup(testutils.WithCfg(cfg))

	emails.InitEmailSender(registry)

	// Register task so that queued tasks can be processed if needed
	if registry.GetTaskStore() == nil {
		t.Fatalf("expected task store to be initialized")
	}

	err := registry.GetTaskStore().TaskHandlerAdd(context.Background(), NewEmailTestTask(registry), true)
	if err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %q", err)
	}

	// Enqueue task
	enqueueHandler, ok := NewEmailTestTask(registry).(*emailTestTask)
	if !ok {
		t.Fatalf("expected *emailTestTask, got different type")
	}

	queuedTask, err := enqueueHandler.Enqueue("test@example.com", "<p>hello</p>")
	if err != nil {
		t.Fatalf("Enqueue() expected nil error, got %q", err)
	}

	// Set queued task and required params on handler
	handlerIface := NewEmailTestTask(registry)
	handler, ok := handlerIface.(*emailTestTask)
	if !ok {
		t.Fatalf("expected *emailTestTask, got different type")
	}

	handler.SetQueuedTask(queuedTask)
	if ok := handler.Handle(); !ok {
		t.Fatalf("Handle() expected true, got false")
	}

	details := handler.QueuedTask().Details()
	if details == "" {
		t.Fatalf("Details() should not be empty after successful Handle")
	}

	if !strings.Contains(details, "Sending email OK.") {
		t.Fatalf("Details() should contain 'Sending email OK.' but got %q", details)
	}
}
