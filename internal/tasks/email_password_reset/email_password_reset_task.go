// Package email_password_reset handles sending password reset emails via the
// background queue.
package email_password_reset

// EmailPasswordResetTask sends password reset emails via the background task
// queue.
//
// The reset token is never stored in the task queue. Instead the task
// receives the nonce under which the forgot-password controller stored
// "email|token" in the memory cache, and resolves it at execution time.
//
// =================================================================
// Example Usage:
//
// 1. Direct execution (explicit parameters):
//    go run ./cmd/server task EmailPasswordResetTask --email="user@example.com" --token="abc123..."
//
// 2. Enqueue for background processing (nonce reference):
//    go run ./cmd/server task EmailPasswordResetTask --nonce="abc123..." --enqueue=yes
//
// Parameters:
// - nonce: The nonce under which "email|token" is stored in the memory cache
//   (used when enqueued by the forgot-password controller)
// - email: Email address to send the link to (direct execution only)
// - token: The reset token (direct execution only)
//
// Optional Parameters:
// - enqueue: Set to "yes" to enqueue task instead of executing immediately
// =================================================================

import (
	"context"
	"encoding/json"
	"errors"
	"project/internal/app"
	"project/internal/emails"
	"project/internal/links"
	"project/internal/tasks/constants"

	"github.com/dracory/taskstore"
)

// passwordResetCacheValue mirrors the JSON payload written by the
// forgot-password controller under "passwordreset:<nonce>" in the memory
// cache.
type passwordResetCacheValue struct {
	Email string `json:"email"`
	Token string `json:"token,omitempty"`
}

// NewEmailPasswordResetTask creates a new task handler for sending password
// reset emails.
//
// Returns:
//   - taskstore.TaskHandlerInterface: The task handler
func NewEmailPasswordResetTask(app app.AppInterface) taskstore.TaskHandlerInterface {
	return &EmailPasswordResetTask{
		app: app,
	}
}

// EmailPasswordResetTask handles sending password reset emails via
// background queue.
type EmailPasswordResetTask struct {
	taskstore.TaskHandlerBase // Embedded base handler for common task operations
	app                       app.AppInterface
}

// Alias returns the unique identifier for this task handler.
// Used when enqueuing and processing tasks.
func (handler *EmailPasswordResetTask) Alias() string {
	return constants.EmailPasswordResetTaskAlias
}

// Title returns a human-readable title for this task.
// Used in task management interfaces.
func (handler *EmailPasswordResetTask) Title() string {
	return "Email Password Reset"
}

// Description returns a detailed description of the task's purpose.
// Used in task management interfaces.
func (handler *EmailPasswordResetTask) Description() string {
	return "Sends a password reset email to a user"
}

// Enqueue adds a new password reset email task to the task queue.
//
// Only the nonce is stored in the task parameters. The email and token are
// resolved from the memory cache at execution time, so the reset credential
// is never persisted in the task store.
//
// Parameters:
//   - nonce: The cache key suffix under which "email|token" is stored
//
// Returns:
//   - taskstore.TaskQueueInterface: The enqueued task
//   - error: Any error that occurred during enqueueing
func (handler *EmailPasswordResetTask) Enqueue(nonce string) (task taskstore.TaskQueueInterface, err error) {
	// Validate task store is initialized
	if handler.app == nil || handler.app.GetConfig() == nil {
		return nil, errors.New("app/config is nil")
	}

	if handler.app.IsDisabledTaskStore() {
		return nil, errors.New("task store is nil")
	}

	if nonce == "" {
		return nil, errors.New("nonce is required parameter")
	}

	// Enqueue task with the nonce reference
	return handler.app.GetTaskStore().TaskDefinitionEnqueueByAlias(
		context.Background(),
		taskstore.DefaultQueueName,
		handler.Alias(),
		map[string]any{
			"nonce": nonce,
		},
	)
}

// Handle processes the password reset email task by either:
// 1. Executing the email sending immediately, or
// 2. Enqueuing the task for background processing if --enqueue=yes is specified
//
// Returns:
//   - bool: true if task was processed successfully, false otherwise
func (handler *EmailPasswordResetTask) Handle() bool {
	if handler.app == nil || handler.app.GetConfig() == nil {
		handler.LogError("App/Config is nil. Aborted.")
		return false
	}

	// Get parameters from task
	nonce := handler.GetParam("nonce")
	email := handler.GetParam("email")
	token := handler.GetParam("token")

	// Check if task should be enqueued instead of executed directly
	if !handler.HasQueuedTask() && handler.GetParam("enqueue") == "yes" {
		_, err := handler.Enqueue(nonce)

		if err != nil {
			handler.LogError("Error enqueuing task: " + err.Error())
		} else {
			handler.LogSuccess("Task enqueued.")
		}

		return true
	}

	// Resolve email+token from the memory cache when only the nonce was
	// provided (the normal path — the token is never stored in the queue).
	if email == "" || token == "" {
		if nonce == "" {
			handler.LogError("nonce (or email and token) is required parameter")
			return false
		}

		email, token = handler.resolveFromCache(nonce)

		if email == "" || token == "" {
			handler.LogError("Password reset token not found or expired for the given nonce")
			return false
		}
	}

	handler.LogInfo("Parameters ok ...")

	resetURL := links.Auth().PasswordReset(map[string]string{"token": token})

	// Send password reset email using the email service
	err := emails.NewEmailPasswordReset(handler.app).Send(email, resetURL)

	if err != nil {
		handler.LogError("Failed to send password reset email: " + err.Error())
		return false
	}

	// The email went out — drop the nonce entry so it does not linger in
	// the memory cache until expiry.
	if nonce != "" {
		if cache := handler.app.GetMemoryCache(); cache != nil {
			cache.Delete("passwordreset:" + nonce)
		}
	}

	handler.LogSuccess("Password reset email sent successfully to " + email)

	return true
}

// resolveFromCache looks up the JSON "email/token" value stored by the
// forgot-password controller under "passwordreset:<nonce>" in the memory
// cache.
func (handler *EmailPasswordResetTask) resolveFromCache(nonce string) (email, token string) {
	cache := handler.app.GetMemoryCache()
	if cache == nil {
		return "", ""
	}

	item := cache.Get("passwordreset:" + nonce)
	if item == nil {
		return "", ""
	}

	stored, ok := item.Value().(string)
	if !ok {
		return "", ""
	}

	var value passwordResetCacheValue
	if err := json.Unmarshal([]byte(stored), &value); err != nil {
		return "", ""
	}

	return value.Email, value.Token
}
