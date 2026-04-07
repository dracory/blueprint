package detailscomponent

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"project/internal/registry"

	"github.com/dracory/api"
)

type CategoryDetailsRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ParentID    string `json:"parent_id"`
}

func HandleAjaxSaveDetails(registry registry.RegistryInterface, r *http.Request, categoryID string) string {
	if registry.GetShopStore() == nil {
		return api.ErrorWithData("Shop store not available", map[string]any{}).ToString()
	}

	var req CategoryDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.ErrorWithData("Invalid request body", map[string]any{}).ToString()
	}

	category, err := registry.GetShopStore().CategoryFindByID(context.Background(), categoryID)
	if err != nil {
		return api.ErrorWithData("Failed to load category", map[string]any{}).ToString()
	}

	if category == nil {
		return api.ErrorWithData("Category not found", map[string]any{}).ToString()
	}

	category.SetTitle(req.Title)
	category.SetDescription(req.Description)
	category.SetStatus(req.Status)
	category.SetParentID(req.ParentID)

	if err := registry.GetShopStore().CategoryUpdate(context.Background(), category); err != nil {
		slog.Error("Failed to update category", slog.String("error", err.Error()))
		return api.ErrorWithData("Failed to save category", map[string]any{}).ToString()
	}

	return api.SuccessWithData("Category saved successfully", map[string]any{
		"category_id": categoryID,
	}).ToString()
}
