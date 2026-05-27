package category_update

import (
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
)

const (
	actionLoadCategory   = "load-category"
	actionUpdateCategory = "update-category"
)

// == CONTROLLER ==============================================================

type categoryUpdateController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewCategoryUpdateController(registry registry.RegistryInterface) *categoryUpdateController {
	return &categoryUpdateController{registry: registry}
}

func (controller *categoryUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadCategory:
		return controller.handleLoadCategory(w, r)
	case actionUpdateCategory:
		return controller.handleUpdateCategory(w, r)
	default:
		return controller.renderPage(r)
	}
}

func (controller *categoryUpdateController) renderPage(r *http.Request) string {
	if controller.registry.GetShopStore() == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "Shop store is not initialized", links.Admin().Home(), 10)
	}

	authUser := helpers.GetAuthUser(r)
	if authUser == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "You are not logged in. Please login to continue.", links.Admin().Home(), 10)
	}

	categoryID := req.GetStringTrimmed(r, "category_id")
	if categoryID == "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "Category ID is required", links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_CATEGORIES}), 10)
	}

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home()},
		{Name: "Shop", URL: links.Admin().Shop(map[string]string{})},
		{Name: "Categories", URL: links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_CATEGORIES})},
		{Name: "Update Category", URL: ""},
	})

	heading := hb.Heading1().HTML("Update Category")

	linksHelper := shared.NewLinks("/admin/shop")
	initScript := hb.Script(`
		const urlLoadCategory = '` + linksHelper.CategoryUpdate(map[string]string{"action": actionLoadCategory, "category_id": categoryID}) + `';
		const urlUpdateCategory = '` + linksHelper.CategoryUpdate(map[string]string{"action": actionUpdateCategory, "category_id": categoryID}) + `';
		const categoryID = '` + categoryID + `';
	`)

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(initScript).
		Child(hb.Div().ID("app"))

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Update Category | Shop",
		Content: content,
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *categoryUpdateController) handleLoadCategory(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		return api.Error("Shop store not available").ToString()
	}

	categoryID := req.GetStringTrimmed(r, "category_id")
	if categoryID == "" {
		return api.Error("Category ID is required").ToString()
	}

	category, err := shopStore.CategoryFindByID(ctx, categoryID)
	if err != nil || category == nil {
		slog.Error("Failed to load category", "error", err)
		return api.Error("Category not found").ToString()
	}

	return api.SuccessWithData("Category loaded successfully", map[string]any{
		"category": map[string]any{
			"id":          category.GetID(),
			"title":       category.GetTitle(),
			"description": category.GetDescription(),
			"status":      category.GetStatus(),
			"parent_id":   category.GetParentID(),
		},
	}).ToString()
}

func (controller *categoryUpdateController) handleUpdateCategory(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		return api.Error("Shop store not available").ToString()
	}

	var reqBody struct {
		CategoryID  string `json:"category_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		ParentID    string `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqBody.CategoryID == "" {
		return api.Error("Category ID is required").ToString()
	}

	category, err := shopStore.CategoryFindByID(ctx, reqBody.CategoryID)
	if err != nil || category == nil {
		return api.Error("Category not found").ToString()
	}

	category.SetTitle(reqBody.Title)
	category.SetDescription(reqBody.Description)
	category.SetStatus(reqBody.Status)
	category.SetParentID(reqBody.ParentID)

	if err := shopStore.CategoryUpdate(ctx, category); err != nil {
		slog.Error("Failed to update category", "error", err)
		return api.Error("Failed to update category").ToString()
	}

	return api.Success("Category updated successfully").ToString()
}
