package detailscomponent

import (
	"context"
	"project/internal/registry"

	"github.com/dracory/api"
	"github.com/dracory/shopstore"
)

func HandleAjaxListCategories(registry registry.RegistryInterface) string {
	if registry.GetShopStore() == nil {
		return api.ErrorWithData("Shop store not available", map[string]any{}).ToString()
	}

	query := shopstore.NewCategoryQuery()
	categories, err := registry.GetShopStore().CategoryList(context.Background(), query)
	if err != nil {
		return api.ErrorWithData("Failed to load categories", map[string]any{}).ToString()
	}

	categoryList := []map[string]string{}
	for _, cat := range categories {
		categoryList = append(categoryList, map[string]string{
			"id":    cat.GetID(),
			"title": cat.GetTitle(),
		})
	}

	return api.SuccessWithData("Categories loaded successfully", map[string]any{
		"categories": categoryList,
	}).ToString()
}
