package product_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/sb"
	"github.com/dracory/shopstore"
)

func (controller *productManagerController) handleLoadProducts(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		Page        int    `json:"page"`
		PerPage     int    `json:"per_page"`
		SortBy      string `json:"sort_by"`
		Sort        string `json:"sort"`
		Status      string `json:"status"`
		CreatedFrom string `json:"created_from"`
		CreatedTo   string `json:"created_to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
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

	query := shopstore.NewProductQuery().
		SetOffset(reqBody.Page * reqBody.PerPage).
		SetLimit(reqBody.PerPage).
		SetOrderBy(reqBody.SortBy).
		SetSortDirection(reqBody.Sort)

	if reqBody.Status != "" {
		query.SetStatus(reqBody.Status)
	}

	if reqBody.CreatedFrom != "" {
		query.SetCreatedAtGte(reqBody.CreatedFrom + " 00:00:00")
	}

	if reqBody.CreatedTo != "" {
		query.SetCreatedAtLte(reqBody.CreatedTo + " 23:59:59")
	}

	products, err := shopStore.ProductList(ctx, query)
	if err != nil {
		slog.Error("Failed to load products", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to load products").ToString()))
		return ""
	}

	total, err := shopStore.ProductCount(ctx, query)
	if err != nil {
		slog.Error("Failed to count products", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to count products").ToString()))
		return ""
	}

	productList := []map[string]any{}
	for _, product := range products {
		productList = append(productList, map[string]any{
			"id":         product.GetID(),
			"title":      product.GetTitle(),
			"status":     product.GetStatus(),
			"price":      product.GetPrice(),
			"created_at": product.GetCreatedAt(),
			"updated_at": product.GetUpdatedAt(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Products loaded successfully", map[string]any{
		"products": productList,
		"total":    total,
		"page":     reqBody.Page,
		"per_page": reqBody.PerPage,
	}).ToString()))
	return ""
}
