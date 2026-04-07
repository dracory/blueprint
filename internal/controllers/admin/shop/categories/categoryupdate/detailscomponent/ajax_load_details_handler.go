package detailscomponent

import (
	"context"
	"project/internal/registry"

	"github.com/dracory/api"
)

func HandleAjaxLoadDetails(registry registry.RegistryInterface, categoryID string) string {
	if registry.GetShopStore() == nil {
		return api.ErrorWithData("Shop store not available", map[string]any{}).ToString()
	}

	category, err := registry.GetShopStore().CategoryFindByID(context.Background(), categoryID)
	if err != nil {
		return api.ErrorWithData("Failed to load category", map[string]any{}).ToString()
	}

	if category == nil {
		return api.ErrorWithData("Category not found", map[string]any{}).ToString()
	}

	details := map[string]string{
		"id":          category.GetID(),
		"title":       category.GetTitle(),
		"description": category.GetDescription(),
		"status":      category.GetStatus(),
		"parent_id":   category.GetParentID(),
	}

	return api.SuccessWithData("Category loaded successfully", map[string]any{
		"details": details,
	}).ToString()
}
