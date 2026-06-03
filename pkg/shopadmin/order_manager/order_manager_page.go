package order_manager

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
)

var (
	//go:embed orders.html
	ordersHTML string

	//go:embed orders.js
	ordersJS string
)

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
	urlLoadOrders := linksHelper.Orders(map[string]string{"action": actionLoadOrdersAjax})
	urlOrderDetails := linksHelper.OrderDetails(map[string]string{"order_id": "ORDER_ID"})

	html := strings.ReplaceAll(ordersHTML, "urlLoadOrders", "'"+urlLoadOrders+"'")
	html = strings.ReplaceAll(html, "urlOrderDetails", "'"+urlOrderDetails+"'")
	js := strings.ReplaceAll(ordersJS, "urlLoadOrders", "'"+urlLoadOrders+"'")
	js = strings.ReplaceAll(js, "urlOrderDetails", "'"+urlOrderDetails+"'")

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueCDN).
		Child(hb.Raw(html)).
		Child(hb.Script(js))

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
