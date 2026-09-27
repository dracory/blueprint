package email_magic_link

import (
	"context"
	"strings"
	"testing"
	"time"

	"project/internal/emails"
	"project/internal/tasks/constants"
	"project/internal/testutils"

	"github.com/dracory/test"
)

func TestNewEmailMagicLinkTask_InitializesFields(t *testing.T) {
	app := testutils.Setup()

	handlerIface := NewEmailMagicLinkTask(app)
	handler, ok := handlerIface.(*EmailMagicLinkTask)
	if !ok {
		t.Fatalf("expected *EmailMagicLinkTask, got different type")
	}

	if handler == nil {
		t.Fatalf("expected handler to be non-nil")
	}
}

func TestEmailMagicLinkTask_Metadata(t *testing.T) {
	app := testutils.Setup()
	handler := NewEmailMagicLinkTask(app)

	if got, want := handler.Alias(), constants.EmailMagicLinkTaskAlias; got != want {
		t.Fatalf("Alias() = %q, want %q", got, want)
	}

	if got, want := handler.Title(), "Email Magic Link"; got != want {
		t.Fatalf("Title() = %q, want %q", got, want)
	}

	if got, want := handler.Description(), "Sends a magic link login email to a user"; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}
}

func TestEmailMagicLinkTask_Enqueue_AppOrConfigNil(t *testing.T) {
	handler := &EmailMagicLinkTask{}

	if _, err := handler.Enqueue("nonce123"); err == nil {
		t.Fatalf("expected error when app/config is nil, got nil")
	}
}

func TestEmailMagicLinkTask_Enqueue_TaskStoreNil(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	handlerIface := NewEmailMagicLinkTask(app)
	handler, ok := handlerIface.(*EmailMagicLinkTask)
	if !ok {
		t.Fatalf("expected *EmailMagicLinkTask, got different type")
	}

	if _, err := handler.Enqueue("nonce123"); err == nil {
		t.Fatalf("expected error when task store is nil, got nil")
	}
}

func TestEmailMagicLinkTask_Enqueue_MissingParams(t *testing.T) {
	app := testutils.Setup(testutils.WithTaskStore(true))
	handler, ok := NewEmailMagicLinkTask(app).(*EmailMagicLinkTask)
	if !ok {
		t.Fatalf("expected *EmailMagicLinkTask, got different type")
	}

	if _, err := handler.Enqueue(""); err == nil {
		t.Fatalf("expected error when nonce is empty, got nil")
	}
}

func TestEmailMagicLinkTask_Handle_MissingParams(t *testing.T) {
	app := testutils.Setup(testutils.WithTaskStore(true))
	handler, ok := NewEmailMagicLinkTask(app).(*EmailMagicLinkTask)
	if !ok {
		t.Fatalf("expected *EmailMagicLinkTask, got different type")
	}

	if ok := handler.Handle(); ok {
		t.Fatalf("Handle() expected false when nonce is missing, got true")
	}
}

func TestEmailMagicLinkTask_Handle_SendEmail(t *testing.T) {
	server, _, cleanup := test.SetupMailServer(t)
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

	if app.IsDisabledTaskStore() {
		t.Fatalf("expected task store to be initialized")
	}

	if err := app.GetTaskStore().TaskHandlerAdd(
		context.Background(), NewEmailMagicLinkTask(app), true); err != nil {
		t.Fatalf("TaskHandlerAdd() expected nil error, got %q", err)
	}

	// Seed the memory cache the way the login controller does
	nonce := "testnonce123"
	app.GetMemoryCache().Set("magiclink:"+nonce,
		`{"email":"user@test.com","token":"token123"}`, 15*time.Minute)

	handler, ok := NewEmailMagicLinkTask(app).(*EmailMagicLinkTask)
	if !ok {
		t.Fatalf("expected *EmailMagicLinkTask, got different type")
	}

	queuedTask, err := handler.Enqueue(nonce)
	if err != nil {
		t.Fatalf("Enqueue() expected nil error, got %q", err)
	}

	// The plaintext token must never appear in the persisted task parameters
	details := queuedTask.GetParameters()
	if strings.Contains(details, "token123") {
		t.Fatalf("task parameters must not contain the plaintext token, got %q", details)
	}

	handler.SetQueuedTask(queuedTask)
	if ok := handler.Handle(); !ok {
		t.Fatalf("Handle() expected true, got false")
	}

	details = handler.QueuedTask().GetDetails()
	if !strings.Contains(details, "Magic link email sent successfully") {
		t.Fatalf("Details() should contain 'Magic link email sent successfully' but got %q", details)
	}
}
