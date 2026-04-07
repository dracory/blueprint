package categorymanager

import (
	"embed"
	"log/slog"
	"net/http"
	"project/internal/controllers/admin/shop/shared"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/shopstore"
)

//go:embed *.html
//go:embed *.js
var categoriesFiles embed.FS

type categoryManagerController struct {
	registry registry.RegistryInterface
}

func NewCategoryManagerController(registry registry.RegistryInterface) *categoryManagerController {
	return &categoryManagerController{registry: registry}
}

func (controller *categoryManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := r.URL.Query().Get("action")

	switch action {
	case "load-categories":
		return controller.handleLoadCategories(r)
	case "delete-category":
		return controller.handleDeleteCategory(w, r)
	default:
		return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
			Title:   "Categories | Shop",
			Content: hb.Wrap().HTML(controller.renderPage(r)),
			ScriptURLs: []string{
				cdn.Htmx_1_9_4(),
				cdn.Sweetalert2_10(),
			},
			Styles: []string{},
		}).ToHTML()
	}
}

func (controller *categoryManagerController) renderPage(r *http.Request) string {
	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Home",
			URL:  links.Admin().Home(),
		},
		{
			Name: "Shop",
			URL:  shared.NewLinks().Home(map[string]string{}),
		},
		{
			Name: "Category Manager",
			URL:  shared.NewLinks().Categories(map[string]string{}),
		},
	})

	heading := hb.Heading1().
		HTML("Shop. Category Manager")

	// Load HTML template
	htmlContent, err := categoriesFiles.ReadFile("categories.html")
	if err != nil {
		slog.Error("Failed to read categories HTML template", "error", err)
		return hb.Div().HTML("Error loading categories component").ToHTML()
	}

	// Load JavaScript
	jsContent, err := categoriesFiles.ReadFile("categories.js")
	if err != nil {
		slog.Error("Failed to read categories JavaScript file", "error", err)
		return hb.Div().HTML("Error loading categories component").ToHTML()
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const urlCategoriesLoad = '` + shared.NewLinks().Categories(map[string]string{"action": "load-categories"}) + `&X-Requested-With=XMLHttpRequest';
		const urlCategoryDelete = '` + shared.NewLinks().Categories(map[string]string{"action": "delete-category"}) + `&X-Requested-With=XMLHttpRequest';
		const urlCategoryCreate = '` + shared.NewLinks().CategoryCreate(map[string]string{}) + `';
		const urlCategoryUpdate = '` + shared.NewLinks().CategoryUpdate(map[string]string{"category_id": ""}) + `';
	`)

	// HTML template must come before the script that mounts to it
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

	return content.ToHTML()
}

func (controller *categoryManagerController) handleLoadCategories(r *http.Request) string {
	if controller.registry.GetShopStore() == nil {
		return api.ErrorWithData("Shop store not available", map[string]any{}).ToString()
	}

	query := shopstore.NewCategoryQuery()
	categories, err := controller.registry.GetShopStore().CategoryList(r.Context(), query)
	if err != nil {
		return api.ErrorWithData("Failed to load categories", map[string]any{}).ToString()
	}

	categoryList := []map[string]string{}
	for _, cat := range categories {
		categoryList = append(categoryList, map[string]string{
			"id":          cat.GetID(),
			"title":       cat.GetTitle(),
			"description": cat.GetDescription(),
			"status":      cat.GetStatus(),
			"parent_id":   cat.GetParentID(),
		})
	}

	return api.SuccessWithData("Categories loaded successfully", map[string]any{
		"categories": categoryList,
	}).ToString()
}

func (controller *categoryManagerController) handleDeleteCategory(w http.ResponseWriter, r *http.Request) string {
	categoryID := r.FormValue("category_id")

	if categoryID == "" {
		return api.ErrorWithData("Category ID is required", map[string]any{}).ToString()
	}

	if controller.registry.GetShopStore() == nil {
		return api.ErrorWithData("Shop store not available", map[string]any{}).ToString()
	}

	if err := controller.registry.GetShopStore().CategoryDeleteByID(r.Context(), categoryID); err != nil {
		slog.Error("Failed to delete category", slog.String("error", err.Error()))
		return api.ErrorWithData("Failed to delete category", map[string]any{}).ToString()
	}

	return api.SuccessWithData("Category deleted successfully", map[string]any{}).ToString()
}
