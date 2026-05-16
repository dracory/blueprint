package taskadmin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/taskstore"
)

func (a *admin) onModalQueuedTaskRequeueSubmitted(r *http.Request) string {
	taskID := req.GetStringTrimmed(r, "task_id")
	taskParameters := req.GetStringTrimmed(r, "task_parameters")

	if taskID == "" {
		return hb.NewSwal(hb.SwalOptions{Title: "Error", Text: "Task is required"}).ToHTML()
	}

	if taskParameters == "" {
		taskParameters = "{}"
	}

	if !json.Valid([]byte(taskParameters)) {
		return hb.NewSwal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Task Parameters is not valid JSON"}).ToHTML()
	}

	task, err := a.taskStore.TaskDefinitionFindByID(context.Background(), taskID)
	if task == nil {
		return hb.NewSwal(hb.SwalOptions{Title: "Error", Text: "Task not found"}).ToHTML()
	}

	var taskParametersMap map[string]any
	err = json.Unmarshal([]byte(taskParameters), &taskParametersMap)

	if err != nil {
		a.logger.Error("At adminTasks > onModalTaskEnqueueSubmitted", "error", err.Error())
		return hb.NewDiv().Class("alert alert-danger").Text("Task failed to be enqueued").ToHTML()
	}

	_, err = a.taskStore.TaskDefinitionEnqueueByAlias(context.Background(), taskstore.DefaultQueueName, task.Alias(), taskParametersMap)
	if err != nil {
		a.logger.Error("At adminTasks > onModalTaskEnqueueSubmitted", "error", err.Error())
		return hb.NewDiv().Class("alert alert-danger").Text("Task failed to be enqueued").ToHTML()
	}

	response := hb.NewSwal(hb.SwalOptions{Icon: "success", Title: "Success", Text: "Task enqueued successfully"}).ToHTML()

	response += hb.NewScript(`setTimeout(() => {window.location.href = window.location.href;}, 3000);`).ToHTML()

	return response

}
