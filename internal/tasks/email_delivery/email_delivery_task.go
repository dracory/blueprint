// Package email_delivery handles sending auth-package emails via the
// background task queue.
package email_delivery

// AuthEmailDeliveryTask sends emails enqueued by the github.com/dracory/auth
// package's FuncEmailEnqueue hook.
//
// The package stores a JSON Delivery (email, subject, body, secret) under
// "delivery:<key>" in the temporary-key store and enqueues only the opaque
// delivery key — credentials and addresses are never persisted in the task
// queue. This task resolves the key at execution time via
// utils.ConsumeDeliveryKey (resolve + delete in one step, so retried jobs
// cannot double-send) and delivers the fully rendered message — the worker
// does not need to know which auth flow produced it.
//
// =================================================================
// Example Usage:
//
// 1. Direct execution (explicit parameters):
//    go run ./cmd/server task AuthEmailDeliveryTask \
//      --email="user@example.com" --subject="Hi" --body="<p>Hi</p>"
//
// 2. Enqueue for background processing (delivery key reference):
//    go run ./cmd/server task AuthEmailDeliveryTask --delivery_key="abc123..." --enqueue=yes
//
// Parameters:
// - delivery_key: The temporary-key suffix under which the JSON Delivery
//   is stored (used when enqueued by the auth package)
// - email: Recipient address (direct execution only)
// - subject: Email subject (direct execution only)
// - body: Rendered HTML body (direct execution only)
//
// Optional Parameters:
// - enqueue: Set to "yes" to enqueue task instead of executing immediately
// =================================================================

import (
	"context"
	"errors"
	"project/internal/app"
	"project/internal/authadapters"
	"project/internal/emails"
	"project/internal/tasks/constants"

	authutils "github.com/dracory/auth/utils"
	"github.com/dracory/taskstore"
)

// NewAuthEmailDeliveryTask creates a new task handler for sending queued
// auth emails.
//
// Returns:
//   - taskstore.TaskHandlerInterface: The task handler
func NewAuthEmailDeliveryTask(app app.AppInterface) taskstore.TaskHandlerInterface {
	return &AuthEmailDeliveryTask{
		app: app,
	}
}

// AuthEmailDeliveryTask handles sending auth emails via background queue.
type AuthEmailDeliveryTask struct {
	taskstore.TaskHandlerBase // Embedded base handler for common task operations
	app                       app.AppInterface
}

// Alias returns the unique identifier for this task handler.
// Used when enqueuing and processing tasks.
func (handler *AuthEmailDeliveryTask) Alias() string {
	return constants.EmailAuthDeliveryTaskAlias
}

// Title returns a human-readable title for this task.
// Used in task management interfaces.
func (handler *AuthEmailDeliveryTask) Title() string {
	return "Auth Email Delivery"
}

// Description returns a detailed description of the task's purpose.
// Used in task management interfaces.
func (handler *AuthEmailDeliveryTask) Description() string {
	return "Sends an email enqueued by the auth package (login code, magic link, password reset)"
}

// Enqueue adds a new auth email delivery task to the task queue.
//
// Only the delivery key is stored in the task parameters. The email,
// subject, and body are resolved from the temporary-key store at execution
// time, so credentials are never persisted in the task store.
//
// Parameters:
//   - deliveryKey: The key suffix under which the JSON Delivery is stored
//
// Returns:
//   - taskstore.TaskQueueInterface: The enqueued task
//   - error: Any error that occurred during enqueueing
func (handler *AuthEmailDeliveryTask) Enqueue(deliveryKey string) (task taskstore.TaskQueueInterface, err error) {
	// Validate task store is initialized
	if handler.app == nil || handler.app.GetConfig() == nil {
		return nil, errors.New("app/config is nil")
	}

	if handler.app.IsDisabledTaskStore() {
		return nil, errors.New("task store is nil")
	}

	if deliveryKey == "" {
		return nil, errors.New("delivery_key is required parameter")
	}

	// Enqueue task with the delivery key reference
	return handler.app.GetTaskStore().TaskDefinitionEnqueueByAlias(
		context.Background(),
		taskstore.DefaultQueueName,
		handler.Alias(),
		map[string]any{
			"delivery_key": deliveryKey,
		},
	)
}

// Handle processes the auth email delivery task by either:
// 1. Executing the email sending immediately, or
// 2. Enqueuing the task for background processing if --enqueue=yes is specified
//
// Returns:
//   - bool: true if task was processed successfully, false otherwise
func (handler *AuthEmailDeliveryTask) Handle() bool {
	if handler.app == nil || handler.app.GetConfig() == nil {
		handler.LogError("App/Config is nil. Aborted.")
		return false
	}

	deliveryKey := handler.GetParam("delivery_key")
	email := handler.GetParam("email")
	subject := handler.GetParam("subject")
	body := handler.GetParam("body")

	// Check if task should be enqueued instead of executed directly
	if !handler.HasQueuedTask() && handler.GetParam("enqueue") == "yes" {
		_, err := handler.Enqueue(deliveryKey)

		if err != nil {
			handler.LogError("Error enqueuing task: " + err.Error())
		} else {
			handler.LogSuccess("Task enqueued.")
		}

		return true
	}

	// Resolve the delivery from the temporary-key store when only the key
	// was provided (the normal path — nothing sensitive enters the queue).
	if deliveryKey != "" {
		tempKeys := &authadapters.TempKeyAdapter{Store: handler.app.GetCacheStore()}
		delivery, ok := authutils.ConsumeDeliveryKey(tempKeys.Get, tempKeys.Delete, tempKeys.Set, deliveryKey)
		if !ok {
			handler.LogError("Delivery not found or expired for the given key")
			return false
		}
		email = delivery.Email
		subject = delivery.Subject
		body = delivery.Body
	}

	if email == "" || subject == "" || body == "" {
		handler.LogError("delivery_key (or email, subject and body) is required parameter")
		return false
	}

	handler.LogInfo("Parameters ok ...")

	err := emails.SendEmail(emails.SendOptions{
		From:     handler.app.GetConfig().GetMailFromAddress(),
		FromName: handler.app.GetConfig().GetMailFromName(),
		To:       []string{email},
		Subject:  subject,
		HtmlBody: body,
	})

	if err != nil {
		handler.LogError("Failed to send email: " + err.Error())
		return false
	}

	handler.LogSuccess("Email sent successfully to " + email)

	return true
}
