package categoryupdate

import (
	"log/slog"
	"net/http"
	"project/internal/controllers/admin/shop/categories/categoryupdate/detailscomponent"
	"project/internal/controllers/admin/shop/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/shopstore"
)

const (
	ACTION_LOAD_DETAILS    = "load-details"
	ACTION_SAVE_DETAILS    = "save-details"
	ACTION_LIST_CATEGORIES = "list-categories"
)

type categoryUpdateController struct {
	registry registry.RegistryInterface
}

func NewCategoryUpdateController(registry registry.RegistryInterface) *categoryUpdateController {
	return &categoryUpdateController{registry: registry}
}

func (controller *categoryUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	categoryID := req.GetStringTrimmed(r, "category_id")
	action := req.GetStringTrimmed(r, "action")

	if categoryID == "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Category ID is required", links.Admin().Home(), 10)
	}

	category, err := controller.registry.GetShopStore().CategoryFindByID(r.Context(), categoryID)
	if err != nil {
		slog.Error("Error. categoryUpdateController: CategoryFindByID", slog.String("error", err.Error()), slog.String("category_id", categoryID))
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Category not found", links.Admin().Home(), 10)
	}

	if category == nil {
		slog.Warn("Warning. categoryUpdateController: CategoryFindByID", slog.String("error", "Category not found"), slog.String("category_id", categoryID))
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Category not found", links.Admin().Home(), 10)
	}

	switch action {
	case ACTION_LOAD_DETAILS:
		return detailscomponent.HandleAjaxLoadDetails(controller.registry, categoryID)
	case ACTION_SAVE_DETAILS:
		return detailscomponent.HandleAjaxSaveDetails(controller.registry, r, categoryID)
	case ACTION_LIST_CATEGORIES:
		return detailscomponent.HandleAjaxListCategories(controller.registry)
	default:
		return controller.handleRenderPage(r, category, categoryID)
	}
}

func (controller *categoryUpdateController) handleRenderPage(r *http.Request, category shopstore.CategoryInterface, categoryID string) string {
	pageContent := controller.renderPageContent(r, category, categoryID)

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Edit Category | Shop",
		Content: hb.Wrap().HTML(pageContent),
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *categoryUpdateController) renderPageContent(r *http.Request, category shopstore.CategoryInterface, categoryID string) string {
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
		{
			Name: "Edit Category",
			URL:  shared.NewLinks().CategoryUpdate(map[string]string{"category_id": categoryID}),
		},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.NewLinks().Categories(map[string]string{}))

	heading := hb.Heading1().
		HTML("Shop. Category. Edit Category").
		Child(buttonCancel)

	productTitle := hb.Heading2().
		Class("mb-3").
		Text("Category: ").
		Text(category.GetTitle())

	body := detailscomponent.Render(controller.registry, categoryID)

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header d-flex justify-content-between align-items-center").
				Child(hb.Heading4().
					HTML("Category Details").
					Style("margin-bottom:0;display:inline-block;")).
				Child(
					hb.Button().
						Type("button").
						Class("btn btn-primary btn-sm").
						ID("details-save-btn-top").
						Child(hb.I().Class("bi bi-save me-1")).
						Text("Save Details")),
		).
		Child(
			hb.Div().
				Class("card-body").
				Child(body),
		).
		Child(
			hb.Div().
				Class("card-footer d-flex justify-content-end").
				Child(hb.Button().
					Type("button").
					Class("btn btn-primary").
					ID("details-save-btn-bottom").
					Child(hb.I().Class("bi bi-save me-2")).
					Text("Save Details")))

	pageContent := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(productTitle).
		Child(card)

	return pageContent.ToHTML()
}
