package category_manager

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
var categoryFiles embed.FS

const (
	actionLoadCategories         = "load-categories"
	actionCategoryDelete         = "delete-category"
	actionCategoryDeleteSelected = "delete-selected"
)

// == CONTROLLER ==============================================================

type categoryManagerController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewCategoryManagerController(registry registry.RegistryInterface) *categoryManagerController {
	return &categoryManagerController{registry: registry}
}

func (controller *categoryManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadCategories:
		return controller.handleLoadCategories(w, r)
	case actionCategoryDelete:
		return controller.handleCategoryDelete(w, r)
	case actionCategoryDeleteSelected:
		return controller.handleCategoryDeleteSelected(w, r)
	default:
		return controller.renderPage(r)
	}
}

func (controller *categoryManagerController) renderPage(r *http.Request) string {
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
		{Name: "Categories", URL: links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_CATEGORIES})},
	})

	heading := hb.Heading1().HTML("Category Manager")

	linksHelper := shared.NewLinks("/admin/shop")
	initScript := hb.Script(`
		const urlLoadCategories = '` + linksHelper.Categories(map[string]string{"action": actionLoadCategories}) + `';
		const urlCategoryDelete = '` + linksHelper.Categories(map[string]string{"action": actionCategoryDelete}) + `';
		const urlCategoryDeleteSelected = '` + linksHelper.Categories(map[string]string{"action": actionCategoryDeleteSelected}) + `';
	`)

	htmlContent, err := categoryFiles.ReadFile("categories.html")
	if err != nil {
		slog.Error("Failed to read categories HTML template", "error", err)
		return hb.Div().HTML("Error loading categories component").ToHTML()
	}

	jsContent, err := categoryFiles.ReadFile("categories.js")
	if err != nil {
		slog.Error("Failed to read categories JavaScript file", "error", err)
		return hb.Div().HTML("Error loading categories component").ToHTML()
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
		Title:   "Categories | Shop",
		Content: content,
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *categoryManagerController) handleLoadCategories(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		json.NewEncoder(w).Encode(api.Error("Shop store not available"))
		return ""
	}

	var reqBody struct {
		Page    int    `json:"page"`
		PerPage int    `json:"per_page"`
		SortBy  string `json:"sort_by"`
		Sort    string `json:"sort"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		json.NewEncoder(w).Encode(api.Error("Invalid request body"))
		return ""
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

	query := shopstore.NewCategoryQuery().
		SetOffset(reqBody.Page * reqBody.PerPage).
		SetLimit(reqBody.PerPage).
		SetOrderBy(reqBody.SortBy).
		SetSortDirection(reqBody.Sort)

	categories, err := shopStore.CategoryList(ctx, query)
	if err != nil {
		slog.Error("Failed to load categories", "error", err)
		json.NewEncoder(w).Encode(api.Error("Failed to load categories"))
		return ""
	}

	total, err := shopStore.CategoryCount(ctx, query)
	if err != nil {
		slog.Error("Failed to count categories", "error", err)
		json.NewEncoder(w).Encode(api.Error("Failed to count categories"))
		return ""
	}

	categoryList := []map[string]any{}
	for _, category := range categories {
		categoryList = append(categoryList, map[string]any{
			"id":          category.GetID(),
			"title":       category.GetTitle(),
			"description": category.GetDescription(),
			"status":      category.GetStatus(),
			"parent_id":   category.GetParentID(),
			"created_at":  category.GetCreatedAt(),
			"updated_at":  category.GetUpdatedAt(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Categories loaded successfully", map[string]any{
		"categories": categoryList,
		"total":      total,
		"page":       reqBody.Page,
		"per_page":   reqBody.PerPage,
	}).ToString()))
	return ""
}

func (controller *categoryManagerController) handleCategoryDelete(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		CategoryID string `json:"category_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.CategoryID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Category ID is required").ToString()))
		return ""
	}

	if err := shopStore.CategoryDeleteByID(ctx, reqBody.CategoryID); err != nil {
		slog.Error("Failed to delete category", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to delete category").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Category deleted successfully").ToString()))
	return ""
}

func (controller *categoryManagerController) handleCategoryDeleteSelected(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		BulkCategoryIDs []string `json:"bulk_category_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if len(reqBody.BulkCategoryIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("No category IDs provided").ToString()))
		return ""
	}

	for _, categoryID := range reqBody.BulkCategoryIDs {
		if err := shopStore.CategoryDeleteByID(ctx, categoryID); err != nil {
			slog.Error("Failed to delete category", "error", err, "category_id", categoryID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Categories deleted successfully").ToString()))
	return ""
}
