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

func (c *AiTitleGeneratorController) onApproveTitle(r *http.Request) string {
	titleID := req.GetStringTrimmed(r, "record_post_id")
	if titleID == "" {
		return shared.ErrorPopup("Title ID is required").ToHTML()
	}

	customStore := c.registry.GetCustomStore()
	if customStore == nil {
		return shared.ErrorPopup("Custom store not configured").ToHTML()
	}

	record, err := customStore.RecordFindByID(titleID)
	if err != nil {
		return shared.ErrorPopup(fmt.Sprintf("Error finding title: %s", err.Error())).ToHTML()
	}

	if record == nil {
		return shared.ErrorPopup("Title not found").ToHTML()
	}

	if err := record.SetPayloadMapKey("status", blogai.POST_STATUS_APPROVED); err != nil {
		return shared.ErrorPopup(fmt.Sprintf("Error updating title status: %s", err.Error())).ToHTML()
	}
	if err := record.SetPayloadMapKey("updated_at", carbon.Now().ToDateTimeString(carbon.UTC)); err != nil {
		return shared.ErrorPopup(fmt.Sprintf("Error updating title timestamp: %s", err.Error())).ToHTML()
	}

	if err := customStore.RecordUpdate(record); err != nil {
		return shared.ErrorPopup(fmt.Sprintf("Error updating title: %s", err.Error())).ToHTML()
	}

	return hb.Swal(hb.SwalOptions{
		Title:            "Success",
		Text:             "Title approved successfully! Reloading page...",
		Icon:             "success",
		Timer:            3000,
		TimerProgressBar: true,
		RedirectURL:      shared.NewLinks().AiTitleGenerator(),
		RedirectSeconds:  3,
	}).ToHTML()
}
