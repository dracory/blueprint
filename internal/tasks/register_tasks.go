package tasks

import (
	"context"
	"project/internal/tasks/blind_index_rebuild"
	"project/internal/tasks/clean_up"
	"project/internal/tasks/email_admin"
	"project/internal/tasks/email_admin_new_contact"
	"project/internal/tasks/email_admin_new_user_registered"
	"project/internal/tasks/email_test"
	"project/internal/tasks/hello_world"
	"project/internal/tasks/stats"
	"project/internal/types"

	"github.com/dracory/taskstore"
)

// RegisterTasks registers the task handlers to the task store
//
// Parameters:
// - none
//
// Returns:
// - none
func RegisterTasks(app types.RegistryInterface) {
	if app.GetTaskStore() == nil {
		return
	}

	tasks := []taskstore.TaskHandlerInterface{
		blind_index_rebuild.NewBlindIndexRebuildTask(app),
		clean_up.NewCleanUpTask(app),
		email_test.NewEmailTestTask(app),
		email_admin.NewEmailToAdminTask(app),
		email_admin_new_contact.NewEmailToAdminOnNewContactFormSubmittedTaskHandler(app),
		email_admin_new_user_registered.NewEmailToAdminOnNewUserRegisteredTaskHandler(app),
		hello_world.NewHelloWorldTask(app),
		stats.NewStatsVisitorEnhanceTask(app),
	}

	for _, task := range tasks {
		err := app.GetTaskStore().TaskHandlerAdd(context.Background(), task, true)

		if err != nil {
			app.GetLogger().Error("At registerTaskHandlers", "error", "Error registering task: "+task.Alias()+" - "+err.Error())
		}
	}
}
