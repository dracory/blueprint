package email_admin_new_user_registered

import (
	"reflect"
	"strings"
	"testing"

	"project/internal/emails"
	"project/internal/testutils"
)

func TestNewEmailToAdminOnNewUserRegisteredTaskHandler_InitializesFields(t *testing.T) {
	app := testutils.Setup()

	handler := NewEmailToAdminOnNewUserRegisteredTaskHandler(app)

	if handler == nil {
		t.Fatalf("expected handler to be non-nil")
	}

	// verify app is set via reflection since app field is unexported
	v := reflect.ValueOf(handler).Elem().FieldByName("app")
	if !v.IsValid() || v.IsNil() {
		t.Fatalf("expected app to be set on handler")
	}
}

func TestEmailToAdminOnNewUserRegisteredTaskHandler_Metadata(t *testing.T) {
	app := testutils.Setup()
	handler := NewEmailToAdminOnNewUserRegisteredTaskHandler(app)

	if got, want := handler.Alias(), "email-to-admin-on-new-user-registered"; got != want {
		t.Fatalf("Alias() = %q, want %q", got, want)
	}

	if got, want := handler.Title(), "Email to Admin on New User"; got != want {
		t.Fatalf("Title() = %q, want %q", got, want)
	}

	if got, want := handler.Description(), "When a new user is registered to the application an email should be sent to the admin"; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}
}

func TestEmailToAdminOnNewUserRegisteredTaskHandler_Enqueue_AppOrConfigNil(t *testing.T) {
	// handler with nil app should fail
	handler := &emailToAdminOnNewUserRegisteredTaskHandler{}

	if _, err := handler.Enqueue(testutils.USER_01); err == nil {
		t.Fatalf("expected error when app/config is nil, got nil")
	}
}

func TestEmailToAdminOnNewUserRegisteredTaskHandler_Enqueue_TaskStoreNil(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(false)
	cfg.SetUserStoreUsed(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	handler := NewEmailToAdminOnNewUserRegisteredTaskHandler(app)

	if _, err := handler.Enqueue(testutils.USER_01); err == nil {
		t.Fatalf("expected error when task store is nil, got nil")
	}
}

func TestEmailToAdminOnNewUserRegisteredTaskHandler_Handle_MissingUserID(t *testing.T) {
	app := testutils.Setup(testutils.WithTaskStore(true), testutils.WithUserStore(true))
	handler := NewEmailToAdminOnNewUserRegisteredTaskHandler(app)

	// no user_id provided via params or queue
	if ok := handler.Handle(); ok {
		t.Fatalf("Handle() expected false when user_id is missing, got true")
	}
}

func TestEmailToAdminOnNewUserRegisteredTaskHandler_Handle_SendEmail(t *testing.T) {
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
	cfg.SetUserStoreUsed(true)

	app := testutils.Setup(testutils.WithCfg(cfg))

	emails.InitEmailSender(app)

	if app.GetTaskStore() == nil {
		t.Fatalf("expected task store to be initialized")
	}

	if app.GetUserStore() == nil {
		t.Fatalf("expected user store to be initialized")
	}

	// Seed a user for the happy path
	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser() expected nil error, got %q", err)
	}

	if user == nil {
		t.Fatalf("SeedUser() expected non-nil user")
	}

	// Register task
	if err := app.GetTaskStore().TaskHandlerAdd(NewEmailToAdminOnNewUserRegisteredTaskHandler(app), true); err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %q", err)
	}

	// Enqueue task with user ID
	enqueueHandler := NewEmailToAdminOnNewUserRegisteredTaskHandler(app)
	queuedTask, err := enqueueHandler.Enqueue(user.ID())
	if err != nil {
		t.Fatalf("Enqueue() expected nil error, got %q", err)
	}

	// Handle using queued task (parameters come from queue details)
	handler := NewEmailToAdminOnNewUserRegisteredTaskHandler(app)
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
