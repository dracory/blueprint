package cmds

import (
	"context"
	"log"
	"project/internal/types"

	"github.com/dracory/cmd"
	"github.com/dracory/taskstore"
	"github.com/samber/lo"
)

// ExecuteJob executes a job (queued task) with the given arguments.
//
// Example:
// go run . job run --task_id=20231008040147830106 --force=yes
//
// Args: an array of strings representing the arguments for the job.
// Return type: None.
func ExecuteJob(app types.RegistryInterface, args []string) {
	name := "No name"
	argumentsMap := cmd.ArgsToMap(args)
	log.Println("Executing job: ", name, " with arguments: ", argumentsMap, " ...")

	queuedTaskID := lo.ValueOr(argumentsMap, "task_id", "")
	force := lo.ValueOr(argumentsMap, "force", "")

	if queuedTaskID == "" {
		log.Println("Task ID is required and must be the first argument")
		return
	}

	if app.GetTaskStore() == nil {
		log.Println("TaskStore is nil")
		return
	}

	queuedTask, err := app.GetTaskStore().TaskQueueFindByID(context.Background(), queuedTaskID)

	if err != nil {
		log.Println("Task not found: ", queuedTaskID)
		return
	}

	if queuedTask == nil {
		log.Println("Task not found: ", queuedTaskID)
		return
	}

	if queuedTask.Status() == taskstore.TaskQueueStatusRunning {
		log.Println("Task is currently running: ", queuedTaskID, "Aborted")
		return
	}

	if force != "yes" && queuedTask.Status() != taskstore.TaskQueueStatusQueued {
		log.Println("Task is not queued: ", queuedTaskID, " . You can use the --force=yes option to force the execution of the job. Aborted")
		return
	}

	isOK, err := app.GetTaskStore().TaskQueueProcessTask(context.Background(), queuedTask)

	if err != nil {
		log.Println("Error processing task: ", queuedTaskID, " ", err.Error())
		return
	}

	if isOK {
		log.Println("Job: ", queuedTaskID, " run OK")
	} else {
		log.Println("Job: ", queuedTaskID, " run failed")
	}
}
