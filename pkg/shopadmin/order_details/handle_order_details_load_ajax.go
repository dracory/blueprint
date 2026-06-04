package order_details

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/shopstore"
)

const (
	actionLoadOrderDetailsAjax = "load-order-details"
)

func (controller *orderDetailsController) handleOrderDetailsLoadAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}

	ctx := r.Context()

	shopStore := controller.app.GetShopStore()
	if shopStore == nil {
		api.Respond(w, r, api.Error("Shop store not available"))
		return ""
	}

	var reqBody struct {
		OrderID string `json:"order_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	if reqBody.OrderID == "" {
		api.Respond(w, r, api.Error("Order ID is required"))
		return ""
	}

	order, err := shopStore.OrderFindByID(ctx, reqBody.OrderID)
	if err != nil {
		slog.Error("Failed to load order", "error", err)
		api.Respond(w, r, api.Error("Failed to load order"))
		return ""
	}

	if order == nil {
		api.Respond(w, r, api.Error("Order not found"))
		return ""
	}

	customerName := ""
	customerEmail := ""

	if order.GetCustomerID() != "" && controller.app.GetUserStore() != nil {
		customer, err := controller.app.GetUserStore().UserFindByID(ctx, order.GetCustomerID())
		if err == nil && customer != nil {
			customerName = customer.GetFirstName() + " " + customer.GetLastName()
			customerEmail = customer.GetEmail()
		}
	}

	// Fetch order line items
	lineItems, err := shopStore.OrderLineItemList(ctx, shopstore.NewOrderLineItemQuery())
	if err != nil {
		slog.Error("Failed to load order line items", "error", err)
	}

	items := []map[string]any{}
	for _, item := range lineItems {
		if item.GetOrderID() == order.GetID() {
			productName := item.GetProductID()
			product, err := shopStore.ProductFindByID(ctx, item.GetProductID())
			if err == nil && product != nil {
				productName = product.GetTitle()
			}

			items = append(items, map[string]any{
				"id":       item.GetID(),
				"name":     productName,
				"quantity": item.GetQuantity(),
				"price":    item.GetPrice(),
				"total":    item.GetPrice(),
			})
		}
	}

	// Parse shipping address from memo
	var shippingAddress map[string]string
	if order.GetMemo() != "" {
		json.Unmarshal([]byte(order.GetMemo()), &shippingAddress)
	}

	orderData := map[string]any{
		"id":               order.GetID(),
		"status":           order.GetStatus(),
		"created_at":       order.GetCreatedAt(),
		"updated_at":       order.GetUpdatedAt(),
		"customer_id":      order.GetCustomerID(),
		"customer_name":    customerName,
		"customer_email":   customerEmail,
		"items":            items,
		"shipping_address": shippingAddress,
	}

	api.Respond(w, r, api.SuccessWithData("Order details loaded successfully", map[string]any{
		"order": orderData,
	}))
	return ""
}
