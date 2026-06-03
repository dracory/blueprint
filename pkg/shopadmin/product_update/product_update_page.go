package product_update

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
	//go:embed form.html
	formHTML string

	//go:embed form.js
	formJS string
)

func (controller *productUpdateController) renderPage(w http.ResponseWriter, r *http.Request) string {
	if controller.registry.GetShopStore() == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "Shop store is not initialized", links.Admin().Home(), 10)
	}

	authUser := helpers.GetAuthUser(r)
	if authUser == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "You are not logged in. Please login to continue.", links.Admin().Home(), 10)
	}

	productID := req.GetStringTrimmed(r, "product_id")
	if productID == "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "Product ID is required", links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_PRODUCTS}), 10)
	}

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home()},
		{Name: "Shop", URL: links.Admin().Shop(map[string]string{})},
		{Name: "Products", URL: links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_PRODUCTS})},
		{Name: "Update Product", URL: ""},
	})

	heading := hb.Heading1().HTML("Update Product")

	linksHelper := shared.NewLinks("/admin/shop")
	urlLoadProduct := linksHelper.ProductUpdate(map[string]string{"action": actionLoadProduct, "product_id": productID})
	urlUpdateProduct := linksHelper.ProductUpdate(map[string]string{"action": actionUpdateProduct, "product_id": productID})

	js := strings.ReplaceAll(formJS, "urlLoadProduct", "'"+urlLoadProduct+"'")
	js = strings.ReplaceAll(js, "urlUpdateProduct", "'"+urlUpdateProduct+"'")

	appDiv := hb.Div().ID("app-product-update").HTML(formHTML)

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(appDiv)

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Update Product | Shop",
		Content: content,
		ScriptURLs: []string{
			cdn.VueJs_3_5_32(),
			cdn.Notiflix_3_2_8(),
		},
		StyleURLs: []string{
			cdn.Notiflix_3_2_8_CSS(),
		},
		Scripts: []string{js},
	}).ToHTML()
}
