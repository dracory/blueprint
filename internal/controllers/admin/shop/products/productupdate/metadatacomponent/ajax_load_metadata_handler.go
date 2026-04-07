package metadatacomponent

import (
	"context"
	"sort"

	"project/internal/registry"

	"github.com/dracory/api"
	"github.com/spf13/cast"
)

// MetadataItem represents a single metadata entry
type MetadataItem struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// HandleAjaxLoadMetadata handles AJAX requests to load metadata and returns JSON string
func HandleAjaxLoadMetadata(registry registry.RegistryInterface, productID string) string {
	if registry.GetShopStore() == nil {
		return api.ErrorWithData("Shop store not available", map[string]any{}).ToString()
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

	// Initialize as empty slice to avoid null in JSON
	metadataItems := []MetadataItem{}

	// Sort keys alphabetically, excluding 'tags' which is handled separately
	keys := make([]string, 0, len(metas))
	for key := range metas {
		if key != "tags" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	for i, key := range keys {
		metadataItems = append(metadataItems, MetadataItem{
			ID:    cast.ToString(i),
			Key:   key,
			Value: metas[key],
		})
	}

	return api.SuccessWithData("", map[string]any{
		"metadata": metadataItems,
	}).ToString()
}
