package tasks

import (
	"errors"
	"project/internal/types"

	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/taskstore"
	"github.com/spf13/cast"
)

func NewCleanUpTask(app types.AppInterface) taskstore.TaskHandlerInterface {
	return &cleanUpTask{
		app: app,
	}
}

type cleanUpTask struct {
	taskstore.TaskHandlerBase
	app types.AppInterface
}

var _ taskstore.TaskHandlerInterface = (*cleanUpTask)(nil) // verify it extends the task interface

func (t *cleanUpTask) Alias() string {
	return "CleanUpTask"
}

func (t *cleanUpTask) Title() string {
	return "Clean Up"
}

func (t *cleanUpTask) Description() string {
	return "Clean up the database"
}

func (t *cleanUpTask) Enqueue() (task taskstore.QueueInterface, err error) {
	if t.app == nil {
		return nil, errors.New("app is nil")
	}

	if t.app.GetTaskStore() == nil {
		return nil, errors.New("task store is nil")
	}

	return t.app.GetTaskStore().TaskEnqueueByAlias(t.Alias(), map[string]any{})
}

func (t *cleanUpTask) Handle() bool {
	// Defensive: if TaskStore isn't initialized (e.g., during early DI cutover), skip work gracefully
	if t.app.GetTaskStore() == nil {
		t.LogInfo("TaskStore not configured; skipping CleanUpTask run.")
		return true
	}

	if !t.HasQueuedTask() && t.GetParam("enqueue") == "yes" {
		_, err := t.Enqueue()

		if err != nil {
			t.LogError("Error enqueuing task: " + err.Error())
		} else {
			t.LogSuccess("Task enqueued.")
		}

		return true
	}

	purgeSince := carbon.Now(carbon.UTC).SubMinutes(30).ToDateTimeString()

	purgeTasks, err := t.app.GetTaskStore().QueueList(taskstore.QueueQuery().
		SetStatus(taskstore.QueueStatusSuccess).
		SetCreatedAtLte(purgeSince))

	if err != nil {
		t.LogError("Error purging tasks: " + err.Error())
		return false
	}

	t.LogInfo("Purging " + cast.ToString(len(purgeTasks)) + " tasks older than " + purgeSince + " ...")

	for _, purgeTask := range purgeTasks {
		err := t.app.GetTaskStore().QueueDeleteByID(purgeTask.ID())

		if err != nil {
			t.LogError("Error purging task: " + err.Error())
			return false
		}
	}

	return true
}
