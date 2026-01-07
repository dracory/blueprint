package schedules

import (
	"context"
	"log/slog"
	"project/internal/registry"
	tasksStats "project/internal/tasks/stats"

	"github.com/dracory/base/cfmt"
	"github.com/dracory/taskstore"
)

// queueClearJob clears the queue for a specific task
func queueClearJob(registry registry.RegistryInterface) {
	if registry == nil {
		cfmt.Errorln("QueueClearJob called with nil registry; skipping")
		return
	}

	if registry.GetTaskStore() == nil {
		cfmt.Warningln("QueueClearJob skipped; task store not configured.")
		return
	}

	alias := tasksStats.NewStatsVisitorEnhanceTask(registry).Alias()

	taskDefinition, err := registry.GetTaskStore().TaskDefinitionFindByAlias(context.Background(), alias)

	if err != nil {
		registry.GetLogger().Error("QueueClearJob > Failed to find task",
			slog.String("alias", alias),
			slog.String("error", err.Error()))
		return
	}

	if taskDefinition == nil {
		registry.GetLogger().Error("QueueClearJob > StatsVisitorEnhanceTask task not found.")
		return
	}

	// Find all queued tasks by alias
	queuedTasks, err := registry.GetTaskStore().TaskQueueList(
		context.Background(),
		taskstore.TaskQueueQuery().
			SetTaskID(taskDefinition.ID()).
			SetStatus(taskstore.TaskQueueStatusSuccess))

	if err != nil {
		registry.GetLogger().Error("QueueClearJob > Failed to list queued tasks",
			slog.String("alias", alias),
			slog.String("error", err.Error()))
		return
	}

	for _, queuedTask := range queuedTasks {
		err := registry.GetTaskStore().TaskQueueDelete(context.Background(), queuedTask)
		if err != nil {
			registry.GetLogger().Error("QueueClearJob > Failed to delete queued task",
				slog.String("alias", alias),
				slog.String("error", err.Error()))
			return
		}
	}
}
