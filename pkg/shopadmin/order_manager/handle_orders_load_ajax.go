package order_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/sb"
	"github.com/dracory/shopstore"
	"github.com/dracory/userstore"
)

func (controller *orderManagerController) handleOrdersLoadAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}

	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		api.Respond(w, r, api.Error("Shop store not available"))
		return ""
	}

	var reqBody struct {
		Page          int    `json:"page"`
		PerPage       int    `json:"per_page"`
		SortBy        string `json:"sort_by"`
		Sort          string `json:"sort"`
		Status        string `json:"status"`
		CustomerName  string `json:"customer_name"`
		CustomerEmail string `json:"customer_email"`
		OrderID       string `json:"order_id"`
		CreatedFrom   string `json:"created_from"`
		CreatedTo     string `json:"created_to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
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

	query := shopstore.NewOrderQuery().
		SetOffset(reqBody.Page * reqBody.PerPage).
		SetLimit(reqBody.PerPage).
		SetOrderBy(reqBody.SortBy).
		SetSortDirection(reqBody.Sort)

	orders, err := shopStore.OrderList(ctx, query)
	if err != nil {
		slog.Error("Failed to load orders", "error", err)
		api.Respond(w, r, api.Error("Failed to load orders"))
		return ""
	}

	filteredOrders := orders

	// Filter by status
	if reqBody.Status != "" {
		filteredOrders = []shopstore.OrderInterface{}
		for _, order := range orders {
			if order.GetStatus() == reqBody.Status {
				filteredOrders = append(filteredOrders, order)
			}
		}
	}

	// Filter by order ID
	if reqBody.OrderID != "" {
		tempOrders := filteredOrders
		filteredOrders = []shopstore.OrderInterface{}
		for _, order := range tempOrders {
			if strings.Contains(strings.ToLower(order.GetID()), strings.ToLower(reqBody.OrderID)) {
				filteredOrders = append(filteredOrders, order)
			}
		}
	}

	// Filter by customer name/email
	if (reqBody.CustomerName != "" || reqBody.CustomerEmail != "") && controller.registry.GetUserStore() != nil {
		matchingCustomerIDs := []string{}

		users, err := controller.registry.GetUserStore().UserList(ctx, userstore.NewUserQuery())
		if err == nil {
			for _, user := range users {
				fullName := strings.ToLower(user.GetFirstName() + " " + user.GetLastName())
				email := strings.ToLower(user.GetEmail())

				match := false
				if reqBody.CustomerName != "" && strings.Contains(fullName, strings.ToLower(reqBody.CustomerName)) {
					match = true
				}
				if reqBody.CustomerEmail != "" && strings.Contains(email, strings.ToLower(reqBody.CustomerEmail)) {
					match = true
				}

				if match {
					matchingCustomerIDs = append(matchingCustomerIDs, user.GetID())
				}
			}
		}

		if len(matchingCustomerIDs) > 0 {
			tempOrders := filteredOrders
			filteredOrders = []shopstore.OrderInterface{}
			for _, order := range tempOrders {
				for _, customerID := range matchingCustomerIDs {
					if order.GetCustomerID() == customerID {
						filteredOrders = append(filteredOrders, order)
						break
					}
				}
			}
		} else {
			filteredOrders = []shopstore.OrderInterface{}
		}
	}

	// Filter by date range
	if reqBody.CreatedFrom != "" || reqBody.CreatedTo != "" {
		tempOrders := filteredOrders
		filteredOrders = []shopstore.OrderInterface{}
		for _, order := range tempOrders {
			createdAt := order.GetCreatedAt()
			if createdAt == "" {
				continue
			}

			match := true
			if reqBody.CreatedFrom != "" && createdAt < reqBody.CreatedFrom {
				match = false
			}
			if reqBody.CreatedTo != "" && createdAt > reqBody.CreatedTo {
				match = false
			}

			if match {
				filteredOrders = append(filteredOrders, order)
			}
		}
	}

	total, err := shopStore.OrderCount(ctx, query)
	if err != nil {
		slog.Error("Failed to count orders", "error", err)
		api.Respond(w, r, api.Error("Failed to count orders"))
		return ""
	}

	orderList := []map[string]any{}
	for _, order := range filteredOrders {
		customerName := ""
		customerEmail := ""

		if order.GetCustomerID() != "" && controller.registry.GetUserStore() != nil {
			customer, err := controller.registry.GetUserStore().UserFindByID(ctx, order.GetCustomerID())
			if err == nil && customer != nil {
				customerName = customer.GetFirstName() + " " + customer.GetLastName()
				customerEmail = customer.GetEmail()
			}
		}

		orderList = append(orderList, map[string]any{
			FieldID:            order.GetID(),
			FieldStatus:        order.GetStatus(),
			FieldCreatedAt:     order.GetCreatedAt(),
			FieldUpdatedAt:     order.GetUpdatedAt(),
			FieldCustomerID:    order.GetCustomerID(),
			FieldCustomerName:  customerName,
			FieldCustomerEmail: customerEmail,
		})
	}

	api.Respond(w, r, api.SuccessWithData("Orders loaded successfully", map[string]any{
		FieldOrders: orderList,
		FieldTotal:  total,
		"page":      reqBody.Page,
		"per_page":  reqBody.PerPage,
	}))
	return ""
}
