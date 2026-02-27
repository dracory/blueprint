package email_admin

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"project/internal/emails"
	"project/internal/testutils"
)

func TestNewEmailToAdminTask_InitializesFields(t *testing.T) {
	registry := testutils.Setup()

	handlerIface := NewEmailToAdminTask(registry)
	handler, ok := handlerIface.(*emailToAdminTask)
	if !ok {
		t.Fatalf("expected *emailToAdminTask, got different type")
	}

	if handler == nil {
		t.Fatalf("expected handler to be non-nil")
	}

	// verify registry is set via reflection since registry field is unexported
	v := reflect.ValueOf(handler).Elem().FieldByName("registry")
	if !v.IsValid() || v.IsNil() {
		t.Fatalf("expected registry to be set on handler")
	}
}

func TestEmailToAdminTask_Metadata(t *testing.T) {
	registry := testutils.Setup()
	handler := NewEmailToAdminTask(registry)

	if got, want := handler.Alias(), "EmailToAdminTask"; got != want {
		t.Fatalf("Alias() = %q, want %q", got, want)
	}

	if got, want := handler.Title(), "Email to Admin"; got != want {
		t.Fatalf("Title() = %q, want %q", got, want)
	}

	if got, want := handler.Description(), "Sends a notification email to the system administrator"; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}
}

func TestEmailToAdminTask_Enqueue_AppOrConfigNil(t *testing.T) {
	// handler with nil app should fail
	handler := &emailToAdminTask{}

	if _, err := handler.Enqueue("<p>test</p>"); err == nil {
		t.Fatalf("expected error when app/config is nil, got nil")
	}
}

func TestEmailToAdminTask_Enqueue_TaskStoreNil(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(false)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	handlerIface := NewEmailToAdminTask(registry)
	handler, ok := handlerIface.(*emailToAdminTask)
	if !ok {
		t.Fatalf("expected *emailToAdminTask, got different type")
	}

	if _, err := handler.Enqueue("<p>test</p>"); err == nil {
		t.Fatalf("expected error when task store is nil, got nil")
	}
}

func TestEmailToAdminTask_Handle_MissingHtml(t *testing.T) {
	registry := testutils.Setup(testutils.WithTaskStore(true))
	handlerIface := NewEmailToAdminTask(registry)
	handler, ok := handlerIface.(*emailToAdminTask)
	if !ok {
		t.Fatalf("expected *emailToAdminTask, got different type")
	}

	// no html param provided via queue or params
	if ok := handler.Handle(); ok {
		t.Fatalf("Handle() expected false when html is missing, got true")
	}
}

func TestEmailToAdminTask_Handle_SendEmail(t *testing.T) {
	// configure mock SMTP server
	server, cleanup := testutils.SetupMailServer(t)
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

	if registry.GetTaskStore() == nil {
		t.Fatalf("expected task store to be initialized")
	}

	// Register task
	if err := registry.GetTaskStore().TaskHandlerAdd(context.Background(), NewEmailToAdminTask(registry), true); err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %q", err)
	}

	// Enqueue task with HTML
	enqueueHandler, ok := NewEmailToAdminTask(registry).(*emailToAdminTask)
	if !ok {
		t.Fatalf("expected *emailToAdminTask, got different type")
	}

	queuedTask, err := enqueueHandler.Enqueue("<p>hello admin</p>")
	if err != nil {
		t.Fatalf("Enqueue() expected nil error, got %q", err)
	}

	// Handle using queued task (parameters come from queue details)
	handlerIface := NewEmailToAdminTask(registry)
	handler, ok := handlerIface.(*emailToAdminTask)
	if !ok {
		t.Fatalf("expected *emailToAdminTask, got different type")
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
