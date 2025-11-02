package aititlegenerator

import (
	"fmt"
	"net/http"
	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogai"

	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dromara/carbon/v2"
)

func (c *AiTitleGeneratorController) onRejectTitle(r *http.Request) string {
	titleID := req.GetStringTrimmed(r, "record_post_id")
	if titleID == "" {
		return shared.ErrorPopup("Title ID is required").ToHTML()
	}

	record, err := c.app.GetCustomStore().RecordFindByID(titleID)
	if err != nil {
		return shared.ErrorPopup(fmt.Sprintf("Error finding title: %s", err.Error())).ToHTML()
	}

	record.SetPayloadMapKey("status", blogai.POST_STATUS_REJECTED)
	record.SetPayloadMapKey("updated_at", carbon.Now().ToDateTimeString(carbon.UTC))

	if err := c.app.GetCustomStore().RecordUpdate(record); err != nil {
		return shared.ErrorPopup(fmt.Sprintf("Error updating title: %s", err.Error())).ToHTML()
	}

	return hb.Swal(hb.SwalOptions{
		Title:            "Success",
		Text:             "Title rejected successfully! Reloading page...",
		Icon:             "success",
		Timer:            3000,
		TimerProgressBar: true,
		RedirectURL:      shared.NewLinks().AiTitleGenerator(),
		RedirectSeconds:  3,
	}).ToHTML()
}
