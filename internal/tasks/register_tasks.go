package tasks

import (
	taskStats "project/internal/tasks/stats"
	"project/internal/types"

	"github.com/gouniverse/taskstore"
)

// RegisterTasks registers the task handlers to the task store
//
// Parameters:
// - none
//
// Returns:
// - none
func RegisterTasks(app types.AppInterface) {
	if app.GetTaskStore() == nil {
		return
	}

	tasks := []taskstore.TaskHandlerInterface{
		NewBlindIndexRebuildTask(app),
		NewCleanUpTask(app),
		NewEmailToAdminTask(app),
		NewEmailToAdminOnNewContactFormSubmittedTaskHandler(app),
		NewEmailToAdminOnNewUserRegisteredTaskHandler(app),
		NewHelloWorldTask(app),
		taskStats.NewStatsVisitorEnhanceTask(app),
	}

	for _, task := range tasks {
		err := app.GetTaskStore().TaskHandlerAdd(task, true)

		if err != nil {
			app.GetLogger().Error("At registerTaskHandlers", "error", "Error registering task: "+task.Alias()+" - "+err.Error())
		}
	}
}
