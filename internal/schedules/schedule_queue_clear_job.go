package schedules

import (
	"context"
	"log/slog"
	tasksStats "project/internal/tasks/stats"
	"project/internal/types"

	"github.com/dracory/base/cfmt"
	"github.com/dracory/taskstore"
)

// queueClearJob clears the queue for a specific task
func queueClearJob(app types.RegistryInterface) {
	if app == nil {
		cfmt.Errorln("QueueClearJob called with nil app; skipping")
		return
	}

	if app.GetTaskStore() == nil {
		cfmt.Warningln("QueueClearJob skipped; task store not configured.")
		return
	}

	alias := tasksStats.NewStatsVisitorEnhanceTask(app).Alias()

	taskDefinition, err := app.GetTaskStore().TaskDefinitionFindByAlias(context.Background(), alias)

	if err != nil {
		app.GetLogger().Error("QueueClearJob > Failed to find task",
			slog.String("alias", alias),
			slog.String("error", err.Error()))
		return
	}

	if taskDefinition == nil {
		app.GetLogger().Error("QueueClearJob > StatsVisitorEnhanceTask task not found.")
		return
	}

	// Find all queued tasks by alias
	queuedTasks, err := app.GetTaskStore().TaskQueueList(
		context.Background(),
		taskstore.TaskQueueQuery().
			SetTaskID(taskDefinition.ID()).
			SetStatus(taskstore.TaskQueueStatusSuccess))

	if err != nil {
		app.GetLogger().Error("QueueClearJob > Failed to list queued tasks",
			slog.String("alias", alias),
			slog.String("error", err.Error()))
		return
	}

	for _, queuedTask := range queuedTasks {
		err := app.GetTaskStore().TaskQueueDelete(context.Background(), queuedTask)
		if err != nil {
			app.GetLogger().Error("QueueClearJob > Failed to delete queued task",
				slog.String("alias", alias),
				slog.String("error", err.Error()))
			return
		}
	}
}
