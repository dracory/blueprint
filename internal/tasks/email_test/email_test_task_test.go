package email_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"project/internal/emails"
	"project/internal/testutils"
)

func TestNewEmailTestTask_InitializesFields(t *testing.T) {
	app := testutils.Setup()

	handler, ok := NewEmailTestTask(app).(*emailTestTask)
	if !ok {
		t.Fatalf("expected *emailTestTask, got different type")
	}

	if handler == nil {
		t.Fatalf("expected task to be non-nil")
	}

	// verify app is set via reflection since app field is unexported
	v := reflect.ValueOf(handler).Elem().FieldByName("app")
	if !v.IsValid() || v.IsNil() {
		t.Fatalf("expected app to be set on task")
	}
}

func TestEmailTestTask_Metadata(t *testing.T) {
	app := testutils.Setup()
	task := NewEmailTestTask(app)

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
	app := testutils.Setup(testutils.WithCfg(cfg))

	handler, ok := NewEmailTestTask(app).(*emailTestTask)
	if !ok {
		t.Fatalf("expected *emailTestTask, got different type")
	}

	if _, err := handler.Enqueue("test@example.com", "<p>test</p>"); err == nil {
		t.Fatalf("expected error when task store is nil, got nil")
	}
}

func TestEmailTestTask_Enqueue_InvalidParams(t *testing.T) {
	app := testutils.Setup(testutils.WithTaskStore(true))
	handler, ok := NewEmailTestTask(app).(*emailTestTask)
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
	server, cleanup := testutils.SetupMailServer(t)
	defer cleanup()

	cfg := testutils.DefaultConf()
	cfg.SetMailDriver("smtp")
	cfg.SetMailHost("127.0.0.1")
	cfg.SetMailPort(server.PortNumber)
	cfg.SetMailUsername("")
	cfg.SetMailPassword("")
	cfg.SetTaskStoreUsed(true)

	app := testutils.Setup(testutils.WithCfg(cfg))

	emails.InitEmailSender(app)

	// Register task so that queued tasks can be processed if needed
	if app.GetTaskStore() == nil {
		t.Fatalf("expected task store to be initialized")
	}

	err := app.GetTaskStore().TaskHandlerAdd(context.Background(), NewEmailTestTask(app), true)
	if err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %q", err)
	}

	// Enqueue task
	enqueueHandler, ok := NewEmailTestTask(app).(*emailTestTask)
	if !ok {
		t.Fatalf("expected *emailTestTask, got different type")
	}

	queuedTask, err := enqueueHandler.Enqueue("test@example.com", "<p>hello</p>")
	if err != nil {
		t.Fatalf("Enqueue() expected nil error, got %q", err)
	}

	// Set queued task and required params on handler
	handlerIface := NewEmailTestTask(app)
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
