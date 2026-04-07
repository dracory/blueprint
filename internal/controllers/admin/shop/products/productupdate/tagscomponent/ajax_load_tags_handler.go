package tagscomponent

import (
	"context"
	"sort"
	"strings"

	"project/internal/registry"

	"github.com/dracory/api"
)

// TagItem represents a single tag entry
type TagItem struct {
	ID  string `json:"id"`
	Tag string `json:"tag"`
}

// HandleAjaxLoadTags handles AJAX requests to load tags and returns JSON string
func HandleAjaxLoadTags(registry registry.RegistryInterface, productID string) string {
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
	tagItems := []TagItem{}

	// Parse tags from metadata
	if tagsMeta, exists := metas["tags"]; exists && tagsMeta != "" {
		tags := strings.Split(tagsMeta, ",")
		for i, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tagItems = append(tagItems, TagItem{
					ID:  string(rune('a' + i)),
					Tag: tag,
				})
			}
		}
	}

	// Sort tags alphabetically
	sort.Slice(tagItems, func(i, j int) bool {
		return tagItems[i].Tag < tagItems[j].Tag
	})

	return api.SuccessWithData("", map[string]any{
		"tags": tagItems,
	}).ToString()
}
