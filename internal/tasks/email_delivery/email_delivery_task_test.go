package email_delivery

import (
	"testing"

	"project/internal/tasks/constants"
	"project/internal/testutils"

	"github.com/dracory/taskstore"
)

func TestNewAuthEmailDeliveryTask_InitializesFields(t *testing.T) {
	app := testutils.Setup()

	handlerIface := NewAuthEmailDeliveryTask(app)
	handler, ok := handlerIface.(*AuthEmailDeliveryTask)
	if !ok {
		t.Fatalf("expected *AuthEmailDeliveryTask, got different type")
	}

	if handler == nil {
		t.Fatalf("expected handler to be non-nil")
	}
}

func TestAuthEmailDeliveryTask_Metadata(t *testing.T) {
	app := testutils.Setup()
	handler := NewAuthEmailDeliveryTask(app)

	if got, want := handler.Alias(), constants.EmailAuthDeliveryTaskAlias; got != want {
		t.Fatalf("Alias() = %q, want %q", got, want)
	}

	if got, want := handler.Title(), "Auth Email Delivery"; got != want {
		t.Fatalf("Title() = %q, want %q", got, want)
	}

	if handler.Description() == "" {
		t.Fatalf("expected non-empty Description()")
	}
}

func TestAuthEmailDeliveryTask_Enqueue_AppOrConfigNil(t *testing.T) {
	handler := &AuthEmailDeliveryTask{}

	if _, err := handler.Enqueue("key123"); err == nil {
		t.Fatalf("expected error when app/config is nil, got nil")
	}
}

func TestAuthEmailDeliveryTask_Enqueue_TaskStoreNil(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	handlerIface := NewAuthEmailDeliveryTask(app)
	handler, ok := handlerIface.(*AuthEmailDeliveryTask)
	if !ok {
		t.Fatalf("expected *AuthEmailDeliveryTask, got different type")
	}

	if _, err := handler.Enqueue("key123"); err == nil {
		t.Fatalf("expected error when task store is disabled, got nil")
	}
}

func TestAuthEmailDeliveryTask_Enqueue_DeliveryKeyRequired(t *testing.T) {
	app := testutils.Setup()
	handler := NewAuthEmailDeliveryTask(app).(*AuthEmailDeliveryTask)

	if _, err := handler.Enqueue(""); err == nil {
		t.Fatalf("expected error when delivery key is empty, got nil")
	}
}

func TestAuthEmailDeliveryTask_Handle_MissingParams(t *testing.T) {
	app := testutils.Setup()
	handler := NewAuthEmailDeliveryTask(app).(*AuthEmailDeliveryTask)

	if handler.Handle() {
		t.Fatalf("expected Handle() = false with no parameters, got true")
	}
}

func TestAuthEmailDeliveryTask_Handle_UnknownDeliveryKey(t *testing.T) {
	app := testutils.Setup()
	handler := NewAuthEmailDeliveryTask(app).(*AuthEmailDeliveryTask)
	queuedTask, err := taskstore.NewTaskQueue().SetParametersMap(map[string]string{
		"delivery_key": "nonexistent-key",
	})
	if err != nil {
		t.Fatalf("SetParametersMap() unexpected error: %v", err)
	}
	handler.SetQueuedTask(queuedTask)

	if handler.Handle() {
		t.Fatalf("expected Handle() = false for an unknown delivery key, got true")
	}
}
