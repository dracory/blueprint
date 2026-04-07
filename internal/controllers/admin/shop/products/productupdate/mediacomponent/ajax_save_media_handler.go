package mediacomponent

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"project/internal/registry"

	"github.com/dracory/api"
	"github.com/dracory/shopstore"
)

// MediaRequest represents the JSON request for saving media
type MediaRequest struct {
	Media []MediaItem `json:"media"`
}

// HandleAjaxSaveMedia handles AJAX requests to save media and returns JSON string
func HandleAjaxSaveMedia(registry registry.RegistryInterface, r *http.Request, productID string) string {
	if registry.GetShopStore() == nil {
		return api.ErrorWithData("Shop store not available", map[string]any{}).ToString()
	}

	var req MediaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.ErrorWithData("Invalid request body", map[string]any{}).ToString()
	}

	// Delete existing media for this product
	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(productID)
	existingMedias, err := registry.GetShopStore().MediaList(context.Background(), mediaQuery)
	if err != nil {
		return api.ErrorWithData("Failed to load existing media", map[string]any{}).ToString()
	}

	for _, existingMedia := range existingMedias {
		err := registry.GetShopStore().MediaDelete(context.Background(), existingMedia)
		if err != nil {
			return api.ErrorWithData("Failed to delete existing media", map[string]any{}).ToString()
		}
	}

	// Create new media entries
	for i, item := range req.Media {
		if item.URL == "" {
			continue
		}

		media := shopstore.NewMedia()
		media.SetEntityID(productID)
		media.SetURL(item.URL)
		media.SetType(determineMediaType(item.URL))
		media.SetTitle(item.FileName)
		media.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
		media.SetSequence(i)

		err := registry.GetShopStore().MediaCreate(context.Background(), media)
		if err != nil {
			slog.Error("Failed to create media", slog.String("error", err.Error()), slog.String("url", item.URL), slog.String("product_id", productID))
			return api.ErrorWithData("Failed to save media: "+err.Error(), map[string]any{}).ToString()
		}
	}

	return api.SuccessWithData("Media saved successfully", map[string]any{
		"media": req.Media,
	}).ToString()
}

// determineMediaType determines the media type based on file extension
func determineMediaType(url string) string {
	ext := strings.ToLower(filepath.Ext(url))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp", ".ico":
		return "image"
	case ".mp4", ".webm", ".ogv", ".mov", ".avi":
		return "video"
	case ".mp3", ".wav", ".ogg", ".oga", ".aac":
		return "audio"
	default:
		return "file"
	}
}
