package metadatacomponent

import (
	"context"
	"encoding/json"
	"net/http"

	"project/internal/registry"

	"github.com/dracory/api"
)

// MetadataRequest represents the JSON request for saving metadata
type MetadataRequest struct {
	Metadata []MetadataItem `json:"metadata"`
}

// HandleAjaxSaveMetadata handles AJAX requests to save metadata and returns JSON string
func HandleAjaxSaveMetadata(registry registry.RegistryInterface, r *http.Request, productID string) string {
	if registry.GetShopStore() == nil {
		return api.ErrorWithData("Shop store not available", map[string]any{}).ToString()
	}

	var req MetadataRequest
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

	productMetas := make(map[string]string)
	for _, item := range req.Metadata {
		// Skip 'tags' key as it's handled separately by the tags component
		if item.Key != "" && item.Key != "tags" {
			productMetas[item.Key] = item.Value
		}
	}

	if err := product.SetMetas(productMetas); err != nil {
		return api.ErrorWithData("Failed to save metadata: "+err.Error(), map[string]any{}).ToString()
	}

	if err := registry.GetShopStore().ProductUpdate(context.Background(), product); err != nil {
		return api.ErrorWithData("Failed to save metadata: "+err.Error(), map[string]any{}).ToString()
	}

	return api.SuccessWithData("Metadata saved successfully", map[string]any{
		"metadata": req.Metadata,
	}).ToString()
}
