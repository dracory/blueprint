package tasks

import (
	"context"
	"errors"
	"project/internal/emails"
	"project/internal/types"

	"github.com/dracory/taskstore"
)

// NewEmailToAdminOnNewUserRegisteredRegisterTask sends an email to admin email when new user is registered
// =================================================================
// Example:
//
// go run . task email-to-admin-on-new-user-register --userID=12345678
//
// =================================================================
func NewEmailToAdminOnNewUserRegisteredTaskHandler(app types.AppInterface) *emailToAdminOnNewUserRegisteredTaskHandler {
	return &emailToAdminOnNewUserRegisteredTaskHandler{
		app: app,
	}
}

// emailToAdminOnNewUserRegisteredTaskHandler sends a notification email to admin
type emailToAdminOnNewUserRegisteredTaskHandler struct {
	taskstore.TaskHandlerBase
	app types.AppInterface
}

var _ taskstore.TaskHandlerInterface = (*emailToAdminOnNewUserRegisteredTaskHandler)(nil) // verify it extends the task interface

func (handler *emailToAdminOnNewUserRegisteredTaskHandler) Alias() string {
	return "email-to-admin-on-new-user-registered"
}

func (handler *emailToAdminOnNewUserRegisteredTaskHandler) Title() string {
	return "Email to Admin on New User"
}

func (handler *emailToAdminOnNewUserRegisteredTaskHandler) Description() string {
	return "When a new user is registered to the application an email should be sent to the admin"
}

func (handler *emailToAdminOnNewUserRegisteredTaskHandler) Enqueue(userID string) (task taskstore.QueueInterface, err error) {
	if handler.app == nil || handler.app.GetConfig() == nil {
		return nil, errors.New("app/config is nil")
	}

	if handler.app.GetTaskStore() == nil {
		return nil, errors.New("task store is nil")
	}

	return handler.app.GetTaskStore().TaskEnqueueByAlias(handler.Alias(), map[string]any{
		"user_id": userID,
	})
}

func (handler *emailToAdminOnNewUserRegisteredTaskHandler) Handle() bool {
	if handler.app == nil || handler.app.GetConfig() == nil {
		handler.LogError("App/Config is nil. Aborted.")
		return false
	}

	if handler.app.GetUserStore() == nil {
		handler.LogError("User store is nil. Aborted.")
		return false
	}

	userID := handler.GetParam("user_id")

	if userID == "" {
		handler.LogError("user_id is required parameter")
		return false
	}

	if !handler.HasQueuedTask() && handler.GetParam("enqueue") == "yes" {
		_, err := handler.Enqueue(userID)

		if err != nil {
			handler.LogError("Error enqueuing task: " + err.Error())
		} else {
			handler.LogSuccess("Task enqueued.")
		}

		return true
	}

	handler.LogInfo("Parameters ok ...")

	user, errUser := handler.app.GetUserStore().UserFindByID(context.Background(), userID)

	if errUser != nil {
		handler.LogError("Error getting user: " + errUser.Error())
		return false
	}

	if user == nil {
		handler.LogError("User not found")
		return false
	}

	err := emails.NewEmailToAdminOnNewUserRegistered(handler.app.GetConfig()).Send(user.ID())

	if err != nil {
		handler.LogError("Sending email failed. Code: ")
		handler.LogError("Error: " + err.Error())
		return false
	}

	handler.LogSuccess("Sending email OK.")

	return true
}
