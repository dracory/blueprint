// Package tasks implements various background tasks for the application.
// This file specifically handles sending notification emails to administrators.
package tasks

// EmailToAdminTask handles sending notification emails to system administrators.
//
// =================================================================
// Example Usage:
//
// 1. Direct execution:
//    go run . task EmailToAdminTask --html="<p>Notification content</p>"
//
// 2. Enqueue for background processing:
//    go run . task EmailToAdminTask --html="<p>Notification content</p>" --enqueue=yes
//
// Required Parameters:
// - html: HTML content of the email to be sent
//
// Optional Parameters:
// - enqueue: Set to "yes" to enqueue task instead of executing immediately
// =================================================================

import (
	"errors"
	"project/internal/config"
	"project/internal/emails"

	"github.com/gouniverse/taskstore"
)

// NewEmailToAdminTaskHandler constructs a new task handler for sending emails to admin
// Example usage:
//
//	go run . task EmailToAdminTask --html=HTML
//
// Returns:
//   - taskstore.TaskHandlerInterface: The task handler
func NewEmailToAdminTask() taskstore.TaskHandlerInterface {
	return &emailToAdminTask{}
}

// emailToAdminTask send a notification email to admin
type emailToAdminTask struct {
	taskstore.TaskHandlerBase // Embedded base handler for common task operations
}

// Alias returns the unique identifier for this task handler
// Used when enqueuing and processing tasks
func (handler *emailToAdminTask) Alias() string {
	return "EmailToAdminTask"
}

// Title returns a human-readable title for this task
// Used in task management interfaces
func (handler *emailToAdminTask) Title() string {
	return "Email to Admin"
}

// Description returns a detailed description of the task's purpose
// Used in task management interfaces
func (handler *emailToAdminTask) Description() string {
	return "Sends a notification email to the system administrator"
}

// Enqueue adds a new email task to the task queue
// Parameters:
//   - html: The HTML content of the email to send
//
// Returns:
//   - taskstore.QueueInterface: The enqueued task
//   - error: Any error that occurred during enqueueing
func (handler *emailToAdminTask) Enqueue(html string) (task taskstore.QueueInterface, err error) {
	// Validate task store is initialized
	if config.TaskStore == nil {
		return nil, errors.New("task store is nil")
	}

	// Enqueue task with the provided HTML content
	return config.TaskStore.TaskEnqueueByAlias(handler.Alias(), map[string]any{
		"html": html,
	})
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
func (handler *emailToAdminTask) Handle() bool {
	// Get HTML content from task parameters
	html := handler.GetParam("html")

	// Validate required parameter
	if html == "" {
		handler.LogError("html is required parameter")
		return false
	}

	// Check if task should be enqueued instead of executed directly
	if !handler.HasQueuedTask() && handler.GetParam("enqueue") == "yes" {
		_, err := handler.Enqueue(html)

		if err != nil {
			handler.LogError("Error enqueuing task: " + err.Error())
		} else {
			handler.LogSuccess("Task enqueued.")
		}

		return true
	}

	handler.LogInfo("Parameters ok ...")

	// Send email using the email service
	err := emails.NewEmailNotifyAdmin().Send(html)

	if err != nil {
		handler.LogError("Sending email failed. Code: ")
		handler.LogError("Error: " + err.Error())
		return false
	}

	handler.LogSuccess("Sending email OK.")

	return true
}
