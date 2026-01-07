package aititlegenerator

import (
	"fmt"
	"net/http"
	"project/internal/controllers/admin/blog/shared"

	"github.com/dracory/hb"
	"github.com/dracory/req"
)

func (c *AiTitleGeneratorController) onDeleteTitle(r *http.Request) string {
	titleID := req.GetStringTrimmed(r, "record_post_id")
	if titleID == "" {
		return shared.ErrorPopup("Title ID is required").ToHTML()
	}

	err := c.registry.GetCustomStore().RecordDeleteByID(titleID)
	if err != nil {
		return shared.ErrorPopup(fmt.Sprintf("Error deleting title: %s", err.Error())).ToHTML()
	}

	return hb.Swal(hb.SwalOptions{
		Title:            "Success",
		Text:             "Title deleted successfully! Reloading page...",
		Icon:             "success",
		Timer:            3000,
		TimerProgressBar: true,
		RedirectURL:      shared.NewLinks().AiTitleGenerator(),
		RedirectSeconds:  3,
	}).ToHTML()
}
