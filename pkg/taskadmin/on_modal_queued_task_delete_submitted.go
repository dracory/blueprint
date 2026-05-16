package taskadmin

import (
	"context"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/req"
)

func (a *admin) onModalQueuedTaskDeleteSubmitted(r *http.Request) string {
	queueID := req.GetStringTrimmed(r, "queue_id")

	if queueID == "" {
		return hb.NewSwal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Queued task ID is required"}).ToHTML()
	}

	err := a.taskStore.TaskQueueSoftDeleteByID(context.Background(), queueID)

	if err != nil {
		a.logger.Error("At taskadmin > onModalQueuedTaskDeleteSubmitted", "error", err.Error())
		return hb.NewSwal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Queued task failed to be deleted"}).ToHTML()
	}

	return hb.NewSwal(hb.SwalOptions{Icon: "success", Title: "Success", Text: "Queued task successfully deleted"}).ToHTML() +
		hb.NewScript(`setTimeout(function(){window.location.href = window.location.href}, 3000);`).ToHTML()
}
