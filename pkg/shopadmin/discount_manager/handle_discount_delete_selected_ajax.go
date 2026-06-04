package discount_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
)

func (controller *discountManagerController) handleDiscountDeleteSelected(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.app.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		BulkDiscountIDs []string `json:"bulk_discount_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if len(reqBody.BulkDiscountIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("No discount IDs provided").ToString()))
		return ""
	}

	for _, discountID := range reqBody.BulkDiscountIDs {
		if err := shopStore.DiscountDeleteByID(ctx, discountID); err != nil {
			slog.Error("Failed to delete discount", "error", err, "discount_id", discountID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Discounts deleted successfully").ToString()))
	return ""
}
