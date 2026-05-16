package taskadmin

import (
	"context"

	"github.com/dracory/hb"
)

func (a *admin) onModalQueuedTaskDetailsShow(queueID string) string {
	if queueID == "" {
		return hb.NewDiv().Class("alert alert-danger").Text("queue id is required").ToHTML()
	}

	queue, err := a.taskStore.TaskQueueFindByID(context.Background(), queueID)

	if err != nil {
		a.logger.Error("At taskadmin > onModalQueuedTaskDetailsShow", "error", err.Error())
		return hb.NewDiv().Class("alert alert-danger").Text("Error retrieving queued task").ToHTML()
	}

	if queue == nil {
		return hb.NewDiv().Class("alert alert-danger").Text("Queue task not found").ToHTML()
	}

	return a.modalQueuedTaskDetails(queue.Details()).ToHTML()
}
