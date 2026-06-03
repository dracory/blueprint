package discount_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
)

func (controller *discountManagerController) handleDiscountDelete(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		DiscountID string `json:"discount_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.DiscountID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Discount ID is required").ToString()))
		return ""
	}

	if err := shopStore.DiscountDeleteByID(ctx, reqBody.DiscountID); err != nil {
		slog.Error("Failed to delete discount", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to delete discount").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Discount deleted successfully").ToString()))
	return ""
}
