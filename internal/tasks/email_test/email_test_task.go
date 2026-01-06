// This file specifically handles sending test emails.
package email_test

// EmailTestTask sends test email.
//
// =================================================================
// Example Usage:
//
// 1. Direct execution:
//    go run . task EmailTestTask --html="<p>Notification content</p>" --to="test@example.com"
//
// 2. Enqueue for background processing:
//    go run . task EmailTestTask --html="<p>Notification content</p>" --to="test@example.com" --enqueue=yes
//
// Required Parameters:
// - html: HTML content of the email to be sent
// - to: Email address to send the test email to
//
// Optional Parameters:
// - enqueue: Set to "yes" to enqueue task instead of executing immediately
// =================================================================

import (
	"context"
	"errors"
	"project/internal/emails"
	"project/internal/types"

	"github.com/dracory/taskstore"
)

// NewEmailToAdminTaskHandler constructs a new task handler for sending emails to admin
// Example usage:
//
//	go run . task EmailTestTask --html="<p>Notification content</p>" --to="test@example.com"
//
// Returns:
//   - taskstore.TaskHandlerInterface: The task handler
func NewEmailTestTask(app types.RegistryInterface) taskstore.TaskHandlerInterface {
	return &emailTestTask{
		app: app,
	}
}

// emailTestTask send a notification email to admin
type emailTestTask struct {
	taskstore.TaskHandlerBase // Embedded base handler for common task operations
	app                       types.RegistryInterface
}

// Alias returns the unique identifier for this task handler
// Used when enqueuing and processing tasks
func (handler *emailTestTask) Alias() string {
	return "EmailTestTask"
}

// Title returns a human-readable title for this task
// Used in task management interfaces
func (handler *emailTestTask) Title() string {
	return "Email Test"
}

// Description returns a detailed description of the task's purpose
// Used in task management interfaces
func (handler *emailTestTask) Description() string {
	return "Sends a notification email to the system administrator"
}

// Enqueue adds a new email task to the task queue
// Parameters:
//   - html: The HTML content of the email to send
//
// Returns:
//   - taskstore.TaskQueueInterface: The enqueued task
//   - error: Any error that occurred during enqueueing
func (handler *emailTestTask) Enqueue(toEmail, html string) (task taskstore.TaskQueueInterface, err error) {
	// Validate task store is initialized
	if handler.app == nil || handler.app.GetConfig() == nil {
		return nil, errors.New("app/config is nil")
	}

	if handler.app.GetTaskStore() == nil {
		return nil, errors.New("task store is nil")
	}

	if toEmail == "" {
		return nil, errors.New("to is required parameter")
	}

	if html == "" {
		return nil, errors.New("html is required parameter")
	}

	// Enqueue task with the provided HTML content
	return handler.app.GetTaskStore().TaskDefinitionEnqueueByAlias(
		context.Background(),
		taskstore.DefaultQueueName,
		handler.Alias(),
		map[string]any{
			"to":   toEmail,
			"html": html,
		},
	)
}

// Handle processes the email task by either:
// 1. Executing the email sending immediately, or
// 2. Enqueuing the task for background processing if --enqueue=yes is specified
//
// Workflow:
// 1. Validates required HTML parameter
// 2. Checks if task should be enqueued
// 3. Sends email using the configured email service
// 4. Logs success/error messages appropriately
//
// Returns:
//   - bool: true if task was processed successfully, false otherwise
func (handler *emailTestTask) Handle() bool {
	if handler.app == nil || handler.app.GetConfig() == nil {
		handler.LogError("App/Config is nil. Aborted.")
		return false
	}

	// Get HTML content from task parameters
	toEmail := handler.GetParam("to")
	html := handler.GetParam("html")

	// Validate required parameters
	if toEmail == "" {
		handler.LogError("to is required parameter")
		return false
	}

	if html == "" {
		handler.LogError("html is required parameter")
		return false
	}

	// Check if task should be enqueued instead of executed directly
	if !handler.HasQueuedTask() && handler.GetParam("enqueue") == "yes" {
		_, err := handler.Enqueue(toEmail, html)

		if err != nil {
			handler.LogError("Error enqueuing task: " + err.Error())
		} else {
			handler.LogSuccess("Task enqueued.")
		}

		return true
	}

	handler.LogInfo("Parameters ok ...")

	// Send email using the email service
	err := emails.SendEmail(emails.SendOptions{
		From:     handler.app.GetConfig().GetMailFromAddress(),
		FromName: handler.app.GetConfig().GetMailFromName(),
		To:       []string{toEmail},
		HtmlBody: html,
		Subject:  "Test email",
	})

	if err != nil {
		handler.LogError("Sending email failed. Code: ")
		handler.LogError("Error: " + err.Error())
		return false
	}

	handler.LogSuccess("Sending email OK.")

	return true
}
