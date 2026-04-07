package detailscomponent

import (
	"context"

	"project/internal/registry"

	"github.com/dracory/api"
)

// ProductDetails represents the product details structure
type ProductDetails struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
	Status      string `json:"status"`
}

// HandleAjaxLoadDetails handles AJAX requests to load product details and returns JSON string
func HandleAjaxLoadDetails(registry registry.RegistryInterface, productID string) string {
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

	details := ProductDetails{
		ID:          product.GetID(),
		Title:       product.GetTitle(),
		Description: product.GetDescription(),
		Price:       product.GetPrice(),
		Quantity:    product.GetQuantity(),
		Status:      product.GetStatus(),
	}

	return api.SuccessWithData("", map[string]any{
		"details": details,
	}).ToString()
}
