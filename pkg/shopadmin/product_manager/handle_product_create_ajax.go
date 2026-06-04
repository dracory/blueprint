package product_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/shopstore"
)

func (controller *productManagerController) handleProductCreateAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}

	if controller.app.GetShopStore() == nil {
		api.Respond(w, r, api.Error("Shop store not configured"))
		return ""
	}

	var reqBody struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	if strings.TrimSpace(reqBody.Title) == "" {
		api.Respond(w, r, api.Error("Title is required"))
		return ""
	}

	product := shopstore.NewProduct()
	product.SetTitle(strings.TrimSpace(reqBody.Title))

	if err := controller.app.GetShopStore().ProductCreate(r.Context(), product); err != nil {
		controller.app.GetLogger().Error("productManagerController.handleProductCreateAjax", slog.String("error", err.Error()))
		api.Respond(w, r, api.Error("Failed to create product"))
		return ""
	}

	api.Respond(w, r, api.SuccessWithData("Product created successfully", map[string]interface{}{FieldProductID: product.GetID()}))
	return ""
}
