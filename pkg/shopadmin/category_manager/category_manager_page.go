package category_manager

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
	//go:embed categories.html
	categoriesHTML string

	//go:embed categories.js
	categoriesJS string
)

func (controller *categoryManagerController) renderPage(w http.ResponseWriter, r *http.Request) string {
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
		{Name: "Categories", URL: links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_CATEGORIES})},
	})

	heading := hb.Heading1().HTML("Category Manager")

	linksHelper := shared.NewLinks("/admin/shop")
	urlLoadCategories := linksHelper.Categories(map[string]string{"action": actionLoadCategories})
	urlCategoryDelete := linksHelper.Categories(map[string]string{"action": actionCategoryDelete})
	urlCategoryDeleteSelected := linksHelper.Categories(map[string]string{"action": actionCategoryDeleteSelected})

	html := strings.ReplaceAll(categoriesHTML, "urlLoadCategories", "'"+urlLoadCategories+"'")
	html = strings.ReplaceAll(html, "urlCategoryDelete", "'"+urlCategoryDelete+"'")
	html = strings.ReplaceAll(html, "urlCategoryDeleteSelected", "'"+urlCategoryDeleteSelected+"'")

	js := strings.ReplaceAll(categoriesJS, "urlLoadCategories", "'"+urlLoadCategories+"'")
	js = strings.ReplaceAll(js, "urlCategoryDelete", "'"+urlCategoryDelete+"'")
	js = strings.ReplaceAll(js, "urlCategoryDeleteSelected", "'"+urlCategoryDeleteSelected+"'")

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
		Title:   "Categories | Shop",
		Content: content,
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}
