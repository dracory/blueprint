package category_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/sb"
	"github.com/dracory/shopstore"
)

func (controller *categoryManagerController) handleLoadCategories(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		json.NewEncoder(w).Encode(api.Error("Shop store not available"))
		return ""
	}

	var reqBody struct {
		Page    int    `json:"page"`
		PerPage int    `json:"per_page"`
		SortBy  string `json:"sort_by"`
		Sort    string `json:"sort"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		json.NewEncoder(w).Encode(api.Error("Invalid request body"))
		return ""
	}

	if reqBody.Page < 0 {
		reqBody.Page = 0
	}
	if reqBody.PerPage <= 0 {
		reqBody.PerPage = 10
	}
	if reqBody.SortBy == "" {
		reqBody.SortBy = shopstore.COLUMN_CREATED_AT
	}
	if reqBody.Sort == "" {
		reqBody.Sort = sb.DESC
	}

	query := shopstore.NewCategoryQuery().
		SetOffset(reqBody.Page * reqBody.PerPage).
		SetLimit(reqBody.PerPage).
		SetOrderBy(reqBody.SortBy).
		SetSortDirection(reqBody.Sort)

	categories, err := shopStore.CategoryList(ctx, query)
	if err != nil {
		slog.Error("Failed to load categories", "error", err)
		json.NewEncoder(w).Encode(api.Error("Failed to load categories"))
		return ""
	}

	total, err := shopStore.CategoryCount(ctx, query)
	if err != nil {
		slog.Error("Failed to count categories", "error", err)
		json.NewEncoder(w).Encode(api.Error("Failed to count categories"))
		return ""
	}

	categoryList := []map[string]any{}
	for _, category := range categories {
		categoryList = append(categoryList, map[string]any{
			"id":          category.GetID(),
			"title":       category.GetTitle(),
			"description": category.GetDescription(),
			"status":      category.GetStatus(),
			"parent_id":   category.GetParentID(),
			"created_at":  category.GetCreatedAt(),
			"updated_at":  category.GetUpdatedAt(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Categories loaded successfully", map[string]any{
		"categories": categoryList,
		"total":      total,
		"page":       reqBody.Page,
		"per_page":   reqBody.PerPage,
	}).ToString()))
	return ""
}
