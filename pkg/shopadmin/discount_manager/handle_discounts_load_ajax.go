package discount_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/neat"
	"github.com/dracory/shopstore"
)

func (controller *discountManagerController) handleLoadDiscounts(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.app.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		Page    int    `json:"page"`
		PerPage int    `json:"per_page"`
		SortBy  string `json:"sort_by"`
		Sort    string `json:"sort"`
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
		reqBody.Sort = neat.SortDesc
	}

	query := shopstore.NewDiscountQuery().
		SetOffset(reqBody.Page * reqBody.PerPage).
		SetLimit(reqBody.PerPage).
		SetOrderBy(reqBody.SortBy).
		SetSortDirection(reqBody.Sort)

	discounts, err := shopStore.DiscountList(ctx, query)
	if err != nil {
		slog.Error("Failed to load discounts", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to load discounts").ToString()))
		return ""
	}

	total, err := shopStore.DiscountCount(ctx, query)
	if err != nil {
		slog.Error("Failed to count discounts", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to count discounts").ToString()))
		return ""
	}

	discountList := []map[string]any{}
	for _, discount := range discounts {
		discountList = append(discountList, map[string]any{
			"id":         discount.GetID(),
			"code":       discount.GetCode(),
			"amount":     discount.GetAmount(),
			"status":     discount.GetStatus(),
			"created_at": discount.GetCreatedAt(),
			"updated_at": discount.GetUpdatedAt(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Discounts loaded successfully", map[string]any{
		"discounts": discountList,
		"total":     total,
		"page":      reqBody.Page,
		"per_page":  reqBody.PerPage,
	}).ToString()))
	return ""
}
