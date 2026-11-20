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
	// Check if app is nil or task store is disabled
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

	// Only register blind index rebuild task if user store and vault store are enabled
	// and at least one blind index store is enabled
	if app.IsEnabledUserStore() && app.IsEnabledVaultStore() &&
		(app.IsEnabledBlindIndexStoreEmail() ||
			app.IsEnabledBlindIndexStoreFirstName() ||
			app.IsEnabledBlindIndexStoreLastName()) {
		tasks = append(tasks, blind_index_rebuild.NewBlindIndexRebuildTask(app))
	}

	// Only register new user registered email task if user store is enabled
	if app.IsEnabledUserStore() {
		tasks = append(tasks, email_admin_new_user_registered.NewEmailToAdminOnNewUserRegisteredTaskHandler(app))
	}

	// Only register email delivery task if cache store is enabled
	if app.IsEnabledCacheStore() {
		tasks = append(tasks, email_delivery.NewAuthEmailDeliveryTask(app))
	}

	// Only register stats task if stats store is enabled
	if app.IsEnabledStatsStore() {
		tasks = append(tasks, stats.NewStatsVisitorEnhanceTask(app))
	}

	// Register all tasks
	for _, task := range tasks {
		err := app.GetTaskStore().TaskHandlerAdd(context.Background(), task, true)

		if err != nil {
			app.GetLogger().Error("At registerTaskHandlers", "error", "Error registering task: "+task.Alias()+" - "+err.Error())
		}
	}
}
