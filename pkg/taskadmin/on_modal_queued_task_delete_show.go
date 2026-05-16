package taskadmin

import (
	"context"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/req"
)

func (a *admin) onModalQueuedTaskDeleteShow(r *http.Request) string {
	queueID := req.GetStringTrimmed(r, "queue_id")

	if queueID == "" {
		return hb.NewSwal(hb.SwalOptions{Title: "Error", Text: "Queued task ID is required"}).ToHTML()
	}

	queue, err := a.taskStore.TaskQueueFindByID(context.Background(), queueID)

	if err != nil {
		a.logger.Error("At taskadmin > onModalQueuedTaskDeleteShow", "error", err.Error())
		return hb.NewSwal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Error retrieving queued task"}).ToHTML()
	}

	if queue == nil {
		return hb.NewSwal(hb.SwalOptions{Title: "Error", Text: "Queued task not found"}).ToHTML()
	}

	return a.modalQueuedTaskDelete(r, queueID).ToHTML()
}
