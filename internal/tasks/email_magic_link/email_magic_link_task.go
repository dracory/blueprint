// Package email_magic_link handles sending magic link login emails via the
// background queue.
package email_magic_link

// EmailMagicLinkTask sends magic link emails via the background task queue.
//
// The login token is never stored in the task queue. Instead the task
// receives the nonce under which the login controller stored
// "email|token|return" in the memory cache, and resolves it at execution
// time.
//
// =================================================================
// Example Usage:
//
// 1. Direct execution (explicit parameters):
//    go run ./cmd/server task EmailMagicLinkTask --email="user@example.com" --token="abc123..."
//
// 2. Enqueue for background processing (nonce reference):
//    go run ./cmd/server task EmailMagicLinkTask --nonce="abc123..." --enqueue=yes
//
// Parameters:
// - nonce: The nonce under which "email|token" is stored in the memory cache
//   (used when enqueued by the login controller)
// - email: Email address to send the link to (direct execution only)
// - token: The login token (direct execution only)
//
// Optional Parameters:
// - return: Optional relative path to redirect to after login
// - enqueue: Set to "yes" to enqueue task instead of executing immediately
// =================================================================

import (
	"context"
	"errors"
	"project/internal/app"
	"project/internal/emails"
	"project/internal/links"
	"project/internal/tasks/constants"
	"strings"

	"github.com/dracory/taskstore"
)

// NewEmailMagicLinkTask creates a new task handler for sending magic link emails.
//
// Returns:
//   - taskstore.TaskHandlerInterface: The task handler
func NewEmailMagicLinkTask(app app.AppInterface) taskstore.TaskHandlerInterface {
	return &EmailMagicLinkTask{
		app: app,
	}
}

// EmailMagicLinkTask handles sending magic link emails via background queue.
type EmailMagicLinkTask struct {
	taskstore.TaskHandlerBase // Embedded base handler for common task operations
	app                       app.AppInterface
}

// Alias returns the unique identifier for this task handler.
// Used when enqueuing and processing tasks.
func (handler *EmailMagicLinkTask) Alias() string {
	return constants.EmailMagicLinkTaskAlias
}

// Title returns a human-readable title for this task.
// Used in task management interfaces.
func (handler *EmailMagicLinkTask) Title() string {
	return "Email Magic Link"
}

// Description returns a detailed description of the task's purpose.
// Used in task management interfaces.
func (handler *EmailMagicLinkTask) Description() string {
	return "Sends a magic link login email to a user"
}

// Enqueue adds a new magic link email task to the task queue.
//
// Only the nonce is stored in the task parameters. The email, token and
// optional return URL are resolved from the memory cache at execution time,
// so the login credential is never persisted in the task store.
//
// Parameters:
//   - nonce: The cache key suffix under which "email|token|return" is stored
//
// Returns:
//   - taskstore.TaskQueueInterface: The enqueued task
//   - error: Any error that occurred during enqueueing
func (handler *EmailMagicLinkTask) Enqueue(nonce string) (task taskstore.TaskQueueInterface, err error) {
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

// Handle processes the magic link email task by either:
// 1. Executing the email sending immediately, or
// 2. Enqueuing the task for background processing if --enqueue=yes is specified
//
// Returns:
//   - bool: true if task was processed successfully, false otherwise
func (handler *EmailMagicLinkTask) Handle() bool {
	if handler.app == nil || handler.app.GetConfig() == nil {
		handler.LogError("App/Config is nil. Aborted.")
		return false
	}

	// Get parameters from task
	nonce := handler.GetParam("nonce")
	email := handler.GetParam("email")
	token := handler.GetParam("token")
	returnURL := handler.GetParam("return")

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

		email, token, returnURL = handler.resolveFromCache(nonce)

		if email == "" || token == "" {
			handler.LogError("Magic link token not found or expired for the given nonce")
			return false
		}
	}

	handler.LogInfo("Parameters ok ...")

	magicLinkURL := links.Auth().Auth(map[string]string{"token": token})
	if returnURL != "" {
		magicLinkURL = links.Auth().Auth(map[string]string{"token": token, "return": returnURL})
	}

	// Send magic link email using the email service
	err := emails.NewEmailMagicLink(handler.app).Send(email, magicLinkURL)

	if err != nil {
		handler.LogError("Failed to send magic link email: " + err.Error())
		return false
	}

	handler.LogSuccess("Magic link email sent successfully to " + email)

	return true
}

// resolveFromCache looks up the "email|token|return" value stored by the
// login controller under "magiclink:<nonce>" in the memory cache.
func (handler *EmailMagicLinkTask) resolveFromCache(nonce string) (email, token, returnURL string) {
	cache := handler.app.GetMemoryCache()
	if cache == nil {
		return "", "", ""
	}

	item := cache.Get("magiclink:" + nonce)
	if item == nil {
		return "", "", ""
	}

	stored, ok := item.Value().(string)
	if !ok {
		return "", "", ""
	}

	parts := strings.SplitN(stored, "|", 3)
	if len(parts) < 2 {
		return "", "", ""
	}

	email, token = parts[0], parts[1]
	if len(parts) == 3 {
		returnURL = parts[2]
	}

	return email, token, returnURL
}
