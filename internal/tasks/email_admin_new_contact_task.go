package tasks

import (
	"errors"
	"project/internal/emails"
	"project/internal/types"

	"github.com/gouniverse/taskstore"
)

// NewEmailToAdminOnNewContactFormSubmittedTaskHandler sends a notification email to admin
// =================================================================
// Example:
//
// go run . task email-to-admin-on-new-contact-form-submitted --html=HTML
//
// =================================================================
func NewEmailToAdminOnNewContactFormSubmittedTaskHandler(app types.AppInterface) *emailToAdminOnNewContactFormSubmittedTaskHandler {
	return &emailToAdminOnNewContactFormSubmittedTaskHandler{
		app: app,
	}
}

// emailToAdminOnNewContactFormSubmittedTaskHandler sends a notification email to admin
type emailToAdminOnNewContactFormSubmittedTaskHandler struct {
	taskstore.TaskHandlerBase
	app types.AppInterface
}

var _ taskstore.TaskHandlerInterface = (*emailToAdminOnNewContactFormSubmittedTaskHandler)(nil) // verify it extends the task interface

func (handler *emailToAdminOnNewContactFormSubmittedTaskHandler) Alias() string {
	return "email-to-admin-on-new-contact-form-submitted"
}

func (handler *emailToAdminOnNewContactFormSubmittedTaskHandler) Title() string {
	return "Email to Admin on New Contact"
}

func (handler *emailToAdminOnNewContactFormSubmittedTaskHandler) Description() string {
	return "Sends a notification email to admin when a new contact form is submitted"
}

func (handler *emailToAdminOnNewContactFormSubmittedTaskHandler) Enqueue() (task taskstore.QueueInterface, err error) {
	if handler.app == nil {
		return nil, errors.New("app is nil")
	}

	if handler.app.GetTaskStore() == nil {
		return nil, errors.New("task store is nil")
	}
	return handler.app.GetTaskStore().TaskEnqueueByAlias(handler.Alias(), map[string]any{})
}

func (handler *emailToAdminOnNewContactFormSubmittedTaskHandler) Handle() bool {
	if !handler.HasQueuedTask() && handler.GetParam("enqueue") == "yes" {
		_, err := handler.Enqueue()

		if err != nil {
			handler.LogError("Error enqueuing task: " + err.Error())
		} else {
			handler.LogSuccess("Task enqueued.")
		}

		return true
	}

	handler.LogInfo("Parameters ok ...")

	// Initialize emails package with config and send using DI
	emails.Init(handler.app.GetConfig())
	err := emails.NewEmailToAdminOnNewContactFormSubmitted(handler.app.GetConfig()).Send()

	if err != nil {
		handler.LogError("Sending email failed. Code: ")
		handler.LogError("Error: " + err.Error())
		return false
	}

	handler.LogSuccess("Sending email OK.")

	return true
}
