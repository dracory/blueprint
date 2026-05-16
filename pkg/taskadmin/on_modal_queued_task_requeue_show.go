package taskadmin

import (
	"context"

	"github.com/dracory/hb"
)

func (a *admin) onModalQueuedTaskRequeueShow(queueID string) string {
	queue, err := a.taskStore.TaskQueueFindByID(context.Background(), queueID)

	if err != nil {
		a.logger.Error("At taskadmin > onModalQueuedTaskRequeueShow", "error", err.Error())
		return hb.NewDiv().Class("alert alert-danger").Text("Error retrieving queued task").ToHTML()
	}

	if queue == nil {
		return hb.NewDiv().Class("alert alert-danger").Text("Queued task not found").ToHTML()
	}

	return a.modalQueuedTaskRequeue(queue.Parameters()).ToHTML()
}
