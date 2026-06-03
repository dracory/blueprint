package category_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
)

func (controller *categoryManagerController) handleCategoryDeleteSelected(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		BulkCategoryIDs []string `json:"bulk_category_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if len(reqBody.BulkCategoryIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("No category IDs provided").ToString()))
		return ""
	}

	for _, categoryID := range reqBody.BulkCategoryIDs {
		if err := shopStore.CategoryDeleteByID(ctx, categoryID); err != nil {
			slog.Error("Failed to delete category", "error", err, "category_id", categoryID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Categories deleted successfully").ToString()))
	return ""
}
