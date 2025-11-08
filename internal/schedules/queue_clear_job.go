package schedules

import (
	"log/slog"
	tasksStats "project/internal/tasks/stats"
	"project/internal/types"

	"github.com/dracory/base/cfmt"
	"github.com/dracory/taskstore"
)

// queueClearJob clears the queue for a specific task
func queueClearJob(app types.AppInterface) {
	if app == nil {
		cfmt.Errorln("QueueClearJob called with nil app; skipping")
		return
	}

	if app.GetTaskStore() == nil {
		cfmt.Warningln("QueueClearJob skipped; task store not configured.")
		return
	}

	alias := tasksStats.NewStatsVisitorEnhanceTask(app).Alias()

	task, err := app.GetTaskStore().TaskFindByAlias(alias)

	if err != nil {
		app.GetLogger().Error("QueueClearJob > Failed to find task",
			slog.String("alias", alias),
			slog.String("error", err.Error()))
		return
	}

	if task == nil {
		app.GetLogger().Error("QueueClearJob > StatsVisitorEnhanceTask task not found.")
		return
	}

	// Find all queued tasks by alias
	queuedTasks, err := app.GetTaskStore().QueueList(taskstore.QueueQuery().
		SetTaskID(task.ID()).
		SetStatus(taskstore.QueueStatusSuccess))

	if err != nil {
		app.GetLogger().Error("QueueClearJob > Failed to list queued tasks",
			slog.String("alias", alias),
			slog.String("error", err.Error()))
		return
	}

	for _, queuedTask := range queuedTasks {
		err := app.GetTaskStore().QueueDelete(queuedTask)
		if err != nil {
			app.GetLogger().Error("QueueClearJob > Failed to delete queued task",
				slog.String("alias", alias),
				slog.String("error", err.Error()))
			return
		}
	}
}
