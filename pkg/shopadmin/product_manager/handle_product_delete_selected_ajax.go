package product_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
)

func (controller *productManagerController) handleProductDeleteSelected(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Shop store not available").ToString()
	}

	var reqBody struct {
		BulkProductIDs []string `json:"bulk_product_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Invalid request body").ToString()
	}

	if len(reqBody.BulkProductIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("No product IDs provided").ToString()
	}

	for _, productID := range reqBody.BulkProductIDs {
		if err := shopStore.ProductDeleteByID(ctx, productID); err != nil {
			slog.Error("Failed to delete product", "error", err, "product_id", productID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Products deleted successfully").ToString()))
	return ""
}
