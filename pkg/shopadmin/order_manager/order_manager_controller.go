package order_manager

import (
	"embed"
	"encoding/json"
	"log/slog"
	"net/http"

	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/shopadmin/shared"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dracory/shopstore"
)

//go:embed *.html
//go:embed *.js
var orderFiles embed.FS

const (
	actionLoadOrders = "load-orders"
)

// == CONTROLLER ==============================================================

type orderManagerController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewOrderManagerController(registry registry.RegistryInterface) *orderManagerController {
	return &orderManagerController{registry: registry}
}

func (controller *orderManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadOrders:
		return controller.handleLoadOrders(w, r)
	default:
		return controller.renderPage(r)
	}
}

func (controller *orderManagerController) renderPage(r *http.Request) string {
	if controller.registry.GetShopStore() == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "Shop store is not initialized", links.Admin().Home(), 10)
	}

	authUser := helpers.GetAuthUser(r)
	if authUser == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "You are not logged in. Please login to continue.", links.Admin().Home(), 10)
	}

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home()},
		{Name: "Shop", URL: links.Admin().Shop(map[string]string{})},
		{Name: "Orders", URL: links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_ORDERS})},
	})

	heading := hb.Heading1().HTML("Order Manager")

	linksHelper := shared.NewLinks("/admin/shop")
	initScript := hb.Script(`
		const urlLoadOrders = '` + linksHelper.Orders(map[string]string{"action": actionLoadOrders}) + `';
	`)

	htmlContent, err := orderFiles.ReadFile("orders.html")
	if err != nil {
		slog.Error("Failed to read orders HTML template", "error", err)
		return hb.Div().HTML("Error loading orders component").ToHTML()
	}

	jsContent, err := orderFiles.ReadFile("orders.js")
	if err != nil {
		slog.Error("Failed to read orders JavaScript file", "error", err)
		return hb.Div().HTML("Error loading orders component").ToHTML()
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	vueContainer := hb.Div().
		Child(vueCDN).
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueContainer)

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Orders | Shop",
		Content: content,
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *orderManagerController) handleLoadOrders(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Shop store not available").ToString()
	}

	var reqBody struct {
		Page    int    `json:"page"`
		PerPage int    `json:"per_page"`
		SortBy  string `json:"sort_by"`
		Sort    string `json:"sort"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Invalid request body").ToString()
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
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Failed to load orders").ToString()
	}

	total, err := shopStore.OrderCount(ctx, query)
	if err != nil {
		slog.Error("Failed to count orders", "error", err)
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Failed to count orders").ToString()
	}

	orderList := []map[string]any{}
	for _, order := range orders {
		orderList = append(orderList, map[string]any{
			"id":         order.GetID(),
			"status":     order.GetStatus(),
			"created_at": order.GetCreatedAt(),
			"updated_at": order.GetUpdatedAt(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Orders loaded successfully", map[string]any{
		"orders":   orderList,
		"total":    total,
		"page":     reqBody.Page,
		"per_page": reqBody.PerPage,
	}).ToString()))
	return ""
}
