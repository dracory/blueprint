package product_update

import (
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

func (controller *productUpdateController) handleLoadProduct(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		return api.Error("Shop store not available").ToString()
	}

	productID := req.GetStringTrimmed(r, "product_id")
	if productID == "" {
		return api.Error("Product ID is required").ToString()
	}

	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil || product == nil {
		slog.Error("Failed to load product", "error", err)
		return api.Error("Product not found").ToString()
	}

	return api.SuccessWithData("Product loaded successfully", map[string]any{
		"product": map[string]any{
			"id":          product.GetID(),
			"title":       product.GetTitle(),
			"description": product.GetDescription(),
			"status":      product.GetStatus(),
			"price":       product.GetPrice(),
		},
	}).ToString()
}
