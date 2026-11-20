package tasks

import (
	"context"
	"project/internal/app"
	"project/internal/tasks/blind_index_rebuild"
	"project/internal/tasks/clean_up"
	"project/internal/tasks/email_admin"
	"project/internal/tasks/email_admin_new_contact"
	"project/internal/tasks/email_admin_new_user_registered"
	"project/internal/tasks/email_delivery"
	"project/internal/tasks/email_test"
	"project/internal/tasks/hello_world"
	"project/internal/tasks/stats"

	"github.com/dracory/taskstore"
)

// RegisterTasks registers the task handlers to the task store
//
// Parameters:
// - app: the app
//
// Returns:
// - none
func RegisterTasks(app app.AppInterface) {
	if app == nil || app.IsDisabledTaskStore() {
		return
	}

	// Base tasks that are always registered
	tasks := []taskstore.TaskHandlerInterface{
		clean_up.NewCleanUpTask(app),
		email_test.NewEmailTestTask(app),
		email_admin.NewEmailToAdminTask(app),
		email_admin_new_contact.NewEmailToAdminOnNewContactFormSubmittedTaskHandler(app),
		hello_world.NewHelloWorldTask(app),
	}

	// Conditional tasks that are registered based on enabled features
	if app.IsEnabledUserStore() && app.IsEnabledVaultStore() &&
		(app.IsEnabledBlindIndexStoreEmail() ||
			app.IsEnabledBlindIndexStoreFirstName() ||
			app.IsEnabledBlindIndexStoreLastName()) {
		tasks = append(tasks, blind_index_rebuild.NewBlindIndexRebuildTask(app))
	}

	if app.IsEnabledUserStore() {
		tasks = append(tasks, email_admin_new_user_registered.NewEmailToAdminOnNewUserRegisteredTaskHandler(app))
	}

	if app.IsEnabledCacheStore() {
		tasks = append(tasks, email_delivery.NewAuthEmailDeliveryTask(app))
	}

	if app.IsEnabledStatsStore() {
		tasks = append(tasks, stats.NewStatsVisitorEnhanceTask(app))
	}

	for _, task := range tasks {
		err := app.GetTaskStore().TaskHandlerAdd(context.Background(), task, true)

		if err != nil {
			app.GetLogger().Error("At registerTaskHandlers", "error", "Error registering task: "+task.Alias()+" - "+err.Error())
		}
	}
}
