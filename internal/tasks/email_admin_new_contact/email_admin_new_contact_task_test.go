package email_admin_new_contact

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestNewEmailToAdminOnNewContactFormSubmittedTaskHandler_InitializesFields(t *testing.T) {
	app := testutils.Setup()

	handler := NewEmailToAdminOnNewContactFormSubmittedTaskHandler(app)

	if handler == nil {
		t.Fatalf("expected handler to be non-nil")
	}

	// verify app is set via reflection since app field is unexported
	v := reflect.ValueOf(handler).Elem().FieldByName("registry")
	if !v.IsValid() || v.IsNil() {
		t.Fatalf("expected registry to be set on handler")
	}
}

func TestEmailToAdminOnNewContactFormSubmittedTaskHandler_Metadata(t *testing.T) {
	app := testutils.Setup()
	handler := NewEmailToAdminOnNewContactFormSubmittedTaskHandler(app)

	if got, want := handler.Alias(), "email-to-admin-on-new-contact-form-submitted"; got != want {
		t.Fatalf("Alias() = %q, want %q", got, want)
	}

	if got, want := handler.Title(), "Email to Admin on New Contact"; got != want {
		t.Fatalf("Title() = %q, want %q", got, want)
	}

	if got, want := handler.Description(), "Sends a notification email to admin when a new contact form is submitted"; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}
}

func TestEmailToAdminOnNewContactFormSubmittedTaskHandler_Enqueue_AppNil(t *testing.T) {
	// handler with nil app should fail
	handler := &emailToAdminOnNewContactFormSubmittedTaskHandler{}

	if _, err := handler.Enqueue(); err == nil {
		t.Fatalf("expected error when app is nil, got nil")
	}
}

func TestEmailToAdminOnNewContactFormSubmittedTaskHandler_Enqueue_TaskStoreNil(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	handler := NewEmailToAdminOnNewContactFormSubmittedTaskHandler(app)

	if _, err := handler.Enqueue(); err == nil {
		t.Fatalf("expected error when task store is nil, got nil")
	}
}

func TestEmailToAdminOnNewContactFormSubmittedTaskHandler_Handle_SendEmail(t *testing.T) {
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

	if app.GetTaskStore() == nil {
		t.Fatalf("expected task store to be initialized")
	}

	// Register task so that queued tasks can be processed if needed
	if err := app.GetTaskStore().TaskHandlerAdd(context.Background(), NewEmailToAdminOnNewContactFormSubmittedTaskHandler(app), true); err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %q", err)
	}

	// Enqueue task (no extra params required; contact data is derived elsewhere)
	enqueueHandler := NewEmailToAdminOnNewContactFormSubmittedTaskHandler(app)
	queuedTask, err := enqueueHandler.Enqueue()
	if err != nil {
		t.Fatalf("Enqueue() expected nil error, got %q", err)
	}

	// Handle using queued task
	handler := NewEmailToAdminOnNewContactFormSubmittedTaskHandler(app)
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
