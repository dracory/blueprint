package tasks

import (
	"context"
	"project/internal/registry"
	"project/internal/tasks/blind_index_rebuild"
	"project/internal/tasks/clean_up"
	"project/internal/tasks/email_admin"
	"project/internal/tasks/email_admin_new_contact"
	"project/internal/tasks/email_admin_new_user_registered"
	"project/internal/tasks/email_test"
	"project/internal/tasks/hello_world"
	"project/internal/tasks/stats"

	"github.com/dracory/taskstore"
)

// RegisterTasks registers the task handlers to the task store
//
// Parameters:
// - none
//
// Returns:
// - none
func RegisterTasks(registry registry.RegistryInterface) {
	if registry.GetTaskStore() == nil {
		return
	}

	tasks := []taskstore.TaskHandlerInterface{
		blind_index_rebuild.NewBlindIndexRebuildTask(registry),
		clean_up.NewCleanUpTask(registry),
		email_test.NewEmailTestTask(registry),
		email_admin.NewEmailToAdminTask(registry),
		email_admin_new_contact.NewEmailToAdminOnNewContactFormSubmittedTaskHandler(registry),
		email_admin_new_user_registered.NewEmailToAdminOnNewUserRegisteredTaskHandler(registry),
		hello_world.NewHelloWorldTask(registry),
		stats.NewStatsVisitorEnhanceTask(registry),
	}

	for _, task := range tasks {
		err := registry.GetTaskStore().TaskHandlerAdd(context.Background(), task, true)

		if err != nil {
			registry.GetLogger().Error("At registerTaskHandlers", "error", "Error registering task: "+task.Alias()+" - "+err.Error())
		}
	}
}
