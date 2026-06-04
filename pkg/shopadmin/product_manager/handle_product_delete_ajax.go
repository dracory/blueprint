package product_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
)

func (controller *productManagerController) handleProductDelete(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.app.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Shop store not available").ToString()
	}

	var reqBody struct {
		ProductID string `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Invalid request body").ToString()
	}

	if reqBody.ProductID == "" {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Product ID is required").ToString()
	}

	if err := shopStore.ProductDeleteByID(ctx, reqBody.ProductID); err != nil {
		slog.Error("Failed to delete product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Failed to delete product").ToString()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Product deleted successfully").ToString()))
	return ""
}
