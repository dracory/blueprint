package category_create

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"project/internal/app"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/pkg/shopadmin/shared"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/shopstore"
	"github.com/dracory/uid"
)

const (
	actionCreateCategory = "create-category"
)

// == CONTROLLER ==============================================================

type categoryCreateController struct {
	app app.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewCategoryCreateController(app app.AppInterface) *categoryCreateController {
	return &categoryCreateController{app: app}
}

func (controller *categoryCreateController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionCreateCategory:
		return controller.handleCreateCategory(w, r)
	default:
		return controller.renderPage(r)
	}
}

func (controller *categoryCreateController) renderPage(r *http.Request) string {
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
		{Name: "Create Category", URL: ""},
	})

	heading := hb.Heading1().HTML("Create Category")

	linksHelper := shared.NewLinks("/admin/shop")
	initScript := hb.Script(`
		const urlCreateCategory = '` + linksHelper.CategoryCreate(map[string]string{"action": actionCreateCategory}) + `';
	`)

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(initScript).
		Child(hb.Div().ID("app"))

	return layouts.NewAdminLayout(controller.app, r, layouts.Options{
		Title:   "Create Category | Shop",
		Content: content,
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *categoryCreateController) handleCreateCategory(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.app.GetShopStore()
	if shopStore == nil {
		return api.Error("Shop store not available").ToString()
	}

	var reqBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		ParentID    string `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqBody.Title == "" {
		return api.Error("Title is required").ToString()
	}

	category := shopstore.NewCategory()
	category.SetID(uid.HumanUid()[:8])
	category.SetTitle(reqBody.Title)
	category.SetDescription(reqBody.Description)
	category.SetStatus(reqBody.Status)
	category.SetParentID(reqBody.ParentID)

	if err := shopStore.CategoryCreate(ctx, category); err != nil {
		slog.Error("Failed to create category", "error", err)
		return api.Error("Failed to create category").ToString()
	}

	return api.SuccessWithData("Category created successfully", map[string]any{
		"category_id": category.GetID(),
	}).ToString()
}
