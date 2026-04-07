package mediacomponent

import (
	"context"

	"project/internal/registry"

	"github.com/dracory/api"
	"github.com/dracory/shopstore"
)

// MediaItem represents a single media item
type MediaItem struct {
	ID       string `json:"id"`
	FileName string `json:"fileName"`
	URL      string `json:"url"`
	IsMain   bool   `json:"isMain"`
}

// HandleAjaxLoadMedia handles AJAX requests to load media and returns JSON string
func HandleAjaxLoadMedia(registry registry.RegistryInterface, productID string) string {
	if registry.GetShopStore() == nil {
		return api.ErrorWithData("Shop store not available", map[string]any{}).ToString()
	}

	// Load media for this product
	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(productID)
	mediaQuery.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
	medias, err := registry.GetShopStore().MediaList(context.Background(), mediaQuery)
	if err != nil {
		return api.ErrorWithData("Failed to load media", map[string]any{}).ToString()
	}

	// Initialize as empty slice to avoid null in JSON
	mediaItems := []MediaItem{}

	for i, media := range medias {
		mediaItems = append(mediaItems, MediaItem{
			ID:       media.GetID(),
			FileName: media.GetTitle(),
			URL:      media.GetURL(),
			IsMain:   i == 0, // First image is main
		})
	}

	return api.SuccessWithData("", map[string]any{
		"media": mediaItems,
	}).ToString()
}
