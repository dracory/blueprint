package product_update

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dracory/api"
)

func (controller *productUpdateController) handleUpdateProduct(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		return api.Error("Shop store not available").ToString()
	}

	var reqBody struct {
		ProductID   string  `json:"product_id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Status      string  `json:"status"`
		Price       float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqBody.ProductID == "" {
		return api.Error("Product ID is required").ToString()
	}

	product, err := shopStore.ProductFindByID(ctx, reqBody.ProductID)
	if err != nil || product == nil {
		return api.Error("Product not found").ToString()
	}

	product.SetTitle(reqBody.Title)
	product.SetDescription(reqBody.Description)
	product.SetStatus(reqBody.Status)
	product.SetPrice(strconv.FormatFloat(reqBody.Price, 'f', 2, 64))

	if err := shopStore.ProductUpdate(ctx, product); err != nil {
		slog.Error("Failed to update product", "error", err)
		return api.Error("Failed to update product").ToString()
	}

	return api.Success("Product updated successfully").ToString()
}
