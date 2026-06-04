package discount_manager

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
	//go:embed discounts.html
	discountsHTML string

	//go:embed discounts.js
	discountsJS string
)

func (controller *discountManagerController) renderPage(w http.ResponseWriter, r *http.Request) string {
	if controller.app.GetShopStore() == nil {
		return helpers.ToFlashError(controller.app.GetCacheStore(), nil, r, "Shop store is not initialized", links.Admin().Home(), 10)
	}

	authUser := helpers.GetAuthUser(r)
	if authUser == nil {
		return helpers.ToFlashError(controller.app.GetCacheStore(), nil, r, "You are not logged in. Please login to continue.", links.Admin().Home(), 10)
	}

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home()},
		{Name: "Shop", URL: links.Admin().Shop(map[string]string{})},
		{Name: "Discounts", URL: links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_DISCOUNTS})},
	})

	heading := hb.Heading1().HTML("Discount Manager")

	linksHelper := shared.NewLinks("/admin/shop")
	urlLoadDiscounts := linksHelper.Discounts(map[string]string{"action": actionLoadDiscounts})
	urlDiscountDelete := linksHelper.Discounts(map[string]string{"action": actionDiscountDelete})
	urlDiscountDeleteSelected := linksHelper.Discounts(map[string]string{"action": actionDiscountDeleteSelected})

	html := strings.ReplaceAll(discountsHTML, "urlLoadDiscounts", "'"+urlLoadDiscounts+"'")
	html = strings.ReplaceAll(html, "urlDiscountDelete", "'"+urlDiscountDelete+"'")
	html = strings.ReplaceAll(html, "urlDiscountDeleteSelected", "'"+urlDiscountDeleteSelected+"'")

	js := strings.ReplaceAll(discountsJS, "urlLoadDiscounts", "'"+urlLoadDiscounts+"'")
	js = strings.ReplaceAll(js, "urlDiscountDelete", "'"+urlDiscountDelete+"'")
	js = strings.ReplaceAll(js, "urlDiscountDeleteSelected", "'"+urlDiscountDeleteSelected+"'")

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
		Title:   "Discounts | Shop",
		Content: content,
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}
