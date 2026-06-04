package order_details

import (
	_ "embed"
	"net/http"
	"strings"

	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/pkg/shopadmin/shared"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

var (
	//go:embed order_details.html
	orderDetailsHTML string

	//go:embed order_details.js
	orderDetailsJS string
)

func (controller *orderDetailsController) renderPage(r *http.Request) string {
	if controller.app.GetShopStore() == nil {
		return helpers.ToFlashError(controller.app.GetCacheStore(), nil, r, "Shop store is not initialized", links.Admin().Home(), 10)
	}

	authUser := helpers.GetAuthUser(r)
	if authUser == nil {
		return helpers.ToFlashError(controller.app.GetCacheStore(), nil, r, "You are not logged in. Please login to continue.", links.Admin().Home(), 10)
	}

	orderID := req.GetStringTrimmed(r, "order_id")
	if orderID == "" {
		return helpers.ToFlashError(controller.app.GetCacheStore(), nil, r, "Order ID is required", links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_ORDERS}), 10)
	}

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home()},
		{Name: "Shop", URL: links.Admin().Shop(map[string]string{})},
		{Name: "Orders", URL: links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_ORDERS})},
		{Name: "Order Details", URL: links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_ORDER_DETAILS, "order_id": orderID})},
	})

	heading := hb.Heading1().HTML("Order Details: " + orderID)

	linksHelper := shared.NewLinks("/admin/shop")
	urlOrders := linksHelper.Orders(map[string]string{})
	urlLoadOrderDetails := linksHelper.OrderDetails(map[string]string{"action": "load-order-details", "order_id": orderID})

	html := strings.ReplaceAll(orderDetailsHTML, "urlOrders", "'"+urlOrders+"'")
	html = strings.ReplaceAll(html, "urlLoadOrderDetails", "'"+urlLoadOrderDetails+"'")
	html = strings.ReplaceAll(html, "ORDER_ID", orderID)
	js := strings.ReplaceAll(orderDetailsJS, "urlOrders", "'"+urlOrders+"'")
	js = strings.ReplaceAll(js, "urlLoadOrderDetails", "'"+urlLoadOrderDetails+"'")
	js = strings.ReplaceAll(js, "ORDER_ID", "'"+orderID+"'")

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueCDN).
		Child(hb.Raw(html)).
		Child(hb.Script(js))

	return layouts.NewAdminLayout(controller.app, r, layouts.Options{
		Title:   "Order Details | Shop",
		Content: content,
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}
