package tagscomponent

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"project/internal/registry"

	"github.com/dracory/api"
)

// TagRequest represents the JSON request for saving tags
type TagRequest struct {
	Tags []TagItem `json:"tags"`
}

// HandleAjaxSaveTags handles AJAX requests to save tags and returns JSON string
func HandleAjaxSaveTags(registry registry.RegistryInterface, r *http.Request, productID string) string {
	if registry.GetShopStore() == nil {
		return api.ErrorWithData("Shop store not available", map[string]any{}).ToString()
	}

	var req TagRequest
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

	metas, err := product.GetMetas()
	if err != nil {
		return api.ErrorWithData("Failed to get metadata", map[string]any{}).ToString()
	}

	// Build tags string from items
	var tags []string
	for _, item := range req.Tags {
		if item.Tag != "" {
			tags = append(tags, strings.TrimSpace(item.Tag))
		}
	}

	// Update metadata
	if len(tags) > 0 {
		metas["tags"] = strings.Join(tags, ",")
	} else {
		delete(metas, "tags")
	}

	if err := product.SetMetas(metas); err != nil {
		return api.ErrorWithData("Failed to save tags: "+err.Error(), map[string]any{}).ToString()
	}

	if err := registry.GetShopStore().ProductUpdate(context.Background(), product); err != nil {
		return api.ErrorWithData("Failed to save tags: "+err.Error(), map[string]any{}).ToString()
	}

	return api.SuccessWithData("Tags saved successfully", map[string]any{
		"tags": req.Tags,
	}).ToString()
}
