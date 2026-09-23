// Package email_otp handles sending OTP login emails via the background queue.
package email_otp

// EmailOTPTask sends OTP emails via the background task queue.
//
// The plaintext OTP is never stored in the task queue. Instead the task
// receives the nonce under which the login controller stored "email:otp" in
// the memory cache, and resolves it at execution time.
//
// =================================================================
// Example Usage:
//
// 1. Direct execution (explicit parameters):
//    go run ./cmd/server task EmailOTPTask --email="user@example.com" --otp="123456"
//
// 2. Enqueue for background processing (nonce reference):
//    go run ./cmd/server task EmailOTPTask --nonce="abc123..." --enqueue=yes
//
// Parameters:
// - nonce: The nonce under which "email:otp" is stored in the memory cache
//   (used when enqueued by the login controller)
// - email: Email address to send the OTP to (direct execution only)
// - otp: The 6-digit OTP code (direct execution only)
//
// Optional Parameters:
// - enqueue: Set to "yes" to enqueue task instead of executing immediately
// =================================================================

import (
	"context"
	"errors"
	"project/internal/app"
	"project/internal/emails"
	"project/internal/tasks/constants"
	"strings"

	"github.com/dracory/taskstore"
)

// NewEmailOTPTask creates a new task handler for sending OTP emails.
//
// Returns:
//   - taskstore.TaskHandlerInterface: The task handler
func NewEmailOTPTask(app app.AppInterface) taskstore.TaskHandlerInterface {
	return &EmailOTPTask{
		app: app,
	}
}

// EmailOTPTask handles sending OTP emails via background queue.
type EmailOTPTask struct {
	taskstore.TaskHandlerBase // Embedded base handler for common task operations
	app                       app.AppInterface
}

// Alias returns the unique identifier for this task handler.
// Used when enqueuing and processing tasks.
func (handler *EmailOTPTask) Alias() string {
	return constants.EmailOTPTaskAlias
}

// Title returns a human-readable title for this task.
// Used in task management interfaces.
func (handler *EmailOTPTask) Title() string {
	return "Email OTP"
}

// Description returns a detailed description of the task's purpose.
// Used in task management interfaces.
func (handler *EmailOTPTask) Description() string {
	return "Sends an OTP email to a user for login verification"
}

// Enqueue adds a new OTP email task to the task queue.
//
// Only the nonce is stored in the task parameters. The plaintext email and
// OTP are resolved from the memory cache at execution time, so the code is
// never persisted in the task store.
//
// Parameters:
//   - nonce: The cache key suffix under which "email:otp" is stored
//
// Returns:
//   - taskstore.TaskQueueInterface: The enqueued task
//   - error: Any error that occurred during enqueueing
func (handler *EmailOTPTask) Enqueue(nonce string) (task taskstore.TaskQueueInterface, err error) {
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

// Handle processes the OTP email task by either:
// 1. Executing the email sending immediately, or
// 2. Enqueuing the task for background processing if --enqueue=yes is specified
//
// Returns:
//   - bool: true if task was processed successfully, false otherwise
func (handler *EmailOTPTask) Handle() bool {
	if handler.app == nil || handler.app.GetConfig() == nil {
		handler.LogError("App/Config is nil. Aborted.")
		return false
	}

	// Get parameters from task
	nonce := handler.GetParam("nonce")
	email := handler.GetParam("email")
	otp := handler.GetParam("otp")

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

	// Resolve email+OTP from the memory cache when only the nonce was
	// provided (the normal path — the code is never stored in the queue).
	if email == "" || otp == "" {
		if nonce == "" {
			handler.LogError("nonce (or email and otp) is required parameter")
			return false
		}

		email, otp = handler.resolveFromCache(nonce)

		if email == "" || otp == "" {
			handler.LogError("OTP not found or expired for the given nonce")
			return false
		}
	}

	handler.LogInfo("Parameters ok ...")

	// Send OTP email using the email service
	err := emails.NewEmailOtp(handler.app).Send(email, otp)

	if err != nil {
		handler.LogError("Failed to send OTP email: " + err.Error())
		return false
	}

	handler.LogSuccess("OTP email sent successfully to " + email)

	return true
}

// resolveFromCache looks up the "email:otp" value stored by the login
// controller under "otp:<nonce>" in the memory cache.
func (handler *EmailOTPTask) resolveFromCache(nonce string) (email, otp string) {
	cache := handler.app.GetMemoryCache()
	if cache == nil {
		return "", ""
	}

	item := cache.Get("otp:" + nonce)
	if item == nil {
		return "", ""
	}

	stored, ok := item.Value().(string)
	if !ok {
		return "", ""
	}

	parts := strings.SplitN(stored, ":", 2)
	if len(parts) != 2 {
		return "", ""
	}

	return parts[0], parts[1]
}
