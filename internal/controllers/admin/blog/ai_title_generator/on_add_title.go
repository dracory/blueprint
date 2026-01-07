package aititlegenerator

import (
	"fmt"
	"net/http"
	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogai"
	"time"

	"github.com/dracory/customstore"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

func (c *AiTitleGeneratorController) onAddTitle(r *http.Request) string {
	customTitle := req.GetStringTrimmed(r, "custom_title")
	if customTitle == "" {
		return shared.ErrorPopup("Title is required").ToHTML()
	}

	record := customstore.NewRecord(blogai.POST_RECORD_TYPE)
	now := time.Now().UTC().Format(time.RFC3339)
	if err := record.SetPayloadMap(map[string]any{
		"id":         record.ID(),
		"title":      customTitle,
		"status":     blogai.POST_STATUS_PENDING,
		"created_at": now,
		"updated_at": now,
	}); err != nil {
		return shared.ErrorPopup(fmt.Sprintf("Error preparing title payload: %s", err.Error())).ToHTML()
	}

	if err := c.registry.GetCustomStore().RecordCreate(record); err != nil {
		return shared.ErrorPopup(fmt.Sprintf("Error saving title: %s", err.Error())).ToHTML()
	}

	cleanupScript := fmt.Sprintf(`(function(){var modal=document.getElementById('%s');if(modal){modal.remove();}var backdrops=document.getElementsByClassName('%s');while(backdrops.length){backdrops[0].remove();}})();`, modalAddTitleID, modalAddTitleBackdropClass)

	return hb.Wrap().
		Child(hb.Script(cleanupScript)).
		Child(hb.Swal(hb.SwalOptions{
			Title:            "Success",
			Text:             "Custom title added successfully! Reloading page...",
			Icon:             "success",
			Timer:            3000,
			TimerProgressBar: true,
			RedirectURL:      shared.NewLinks().AiTitleGenerator(),
			RedirectSeconds:  3,
		})).
		ToHTML()
}
