package detailscomponent

import (
	"context"
	"encoding/json"
	"net/http"

	"project/internal/registry"

	"github.com/dracory/api"
	"github.com/dracory/shopstore"
)

// DetailsRequest represents the JSON request for saving product details
type DetailsRequest struct {
	Details ProductDetails `json:"details"`
}

// HandleAjaxSaveDetails handles AJAX requests to save product details and returns JSON string
func HandleAjaxSaveDetails(registry registry.RegistryInterface, r *http.Request, productID string) string {
	if registry.GetShopStore() == nil {
		return api.ErrorWithData("Shop store not available", map[string]any{}).ToString()
	}

	var req DetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.ErrorWithData("Invalid request body", map[string]any{}).ToString()
	}

	product, err := registry.GetShopStore().ProductFindByID(context.Background(), productID)
	if err != nil {
		return api.ErrorWithData("Product not found", map[string]any{}).ToString()
	}

	if product == nil {
		return api.ErrorWithData("Product not found", map[string]any{}).ToString()
	}

	// Update product fields
	product.SetTitle(req.Details.Title)
	product.SetDescription(req.Details.Description)
	product.SetPrice(req.Details.Price)
	product.SetQuantity(req.Details.Quantity)

	// Only update status if it's a valid value
	if req.Details.Status == shopstore.PRODUCT_STATUS_ACTIVE ||
		req.Details.Status == shopstore.PRODUCT_STATUS_DISABLED ||
		req.Details.Status == shopstore.PRODUCT_STATUS_DRAFT {
		product.SetStatus(req.Details.Status)
	}

	if err := registry.GetShopStore().ProductUpdate(context.Background(), product); err != nil {
		return api.ErrorWithData("Failed to save product details: "+err.Error(), map[string]any{}).ToString()
	}

	return api.SuccessWithData("Product details saved successfully", map[string]any{
		"details": req.Details,
	}).ToString()
}
