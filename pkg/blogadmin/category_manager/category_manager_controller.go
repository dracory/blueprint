package category_manager

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/blogadmin/shared"

	"github.com/dracory/api"
	"github.com/dracory/blogstore"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
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
	case "create-category":
		return controller.handleCreateCategory(w, r)
	case "update-category":
		return controller.handleUpdateCategory(w, r)
	case "delete-category":
		return controller.handleDeleteCategory(w, r)
	case "reorder-categories":
		return controller.handleReorderCategories(w, r)
	default:
		return controller.renderPage(r)
	}
}

func (controller *categoryManagerController) renderPage(r *http.Request) string {
	authUser := helpers.GetAuthUser(r)
	if authUser == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "You are not logged in. Please login to continue.", links.Admin().Blog(), 10)
	}

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home()},
		{Name: "Blog", URL: links.Admin().Blog()},
		{Name: "Categories", URL: ""},
	})

	heading := hb.Heading1().HTML("Blog. Category Manager")

	htmlContent, err := categoriesFiles.ReadFile("categories.html")
	if err != nil {
		slog.Error("Failed to read categories HTML template", "error", err)
		return hb.Div().HTML("Error loading categories component").ToHTML()
	}

	jsContent, err := categoriesFiles.ReadFile("categories.js")
	if err != nil {
		slog.Error("Failed to read categories JavaScript file", "error", err)
		return hb.Div().HTML("Error loading categories component").ToHTML()
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")
	sortableCDN := hb.Script("").Src("https://cdn.jsdelivr.net/npm/sortablejs@1.15.0/Sortable.min.js")

	initScript := hb.Script(`
		const urlCategoriesLoad = '` + shared.NewLinks("/admin/blog").CategoryManager(map[string]string{"action": "load-categories"}) + `';
		const urlCategoryCreate = '` + shared.NewLinks("/admin/blog").CategoryManager(map[string]string{"action": "create-category"}) + `';
		const urlCategoryUpdate = '` + shared.NewLinks("/admin/blog").CategoryManager(map[string]string{"action": "update-category", "category_id": "CATEGORY_ID_PLACEHOLDER"}) + `';
		const urlCategoryDelete = '` + shared.NewLinks("/admin/blog").CategoryManager(map[string]string{"action": "delete-category"}) + `';
		const urlCategoriesReorder = '` + shared.NewLinks("/admin/blog").CategoryManager(map[string]string{"action": "reorder-categories"}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	vueContainer := hb.Div().
		Child(vueCDN).
		Child(sortableCDN).
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
		Title:   "Blog | Category Manager",
		Content: content,
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *categoryManagerController) handleLoadCategories(r *http.Request) string {
	ctx := r.Context()

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	categoryTaxonomy, err := controller.ensureTaxonomy(ctx, blogStore)
	if err != nil {
		return api.Error("Failed to ensure taxonomy: " + err.Error()).ToString()
	}

	terms, err := blogStore.TermList(ctx, blogstore.TermQueryOptions{
		TaxonomyID: categoryTaxonomy.GetID(),
		OrderBy:    "sequence",
		SortOrder:  "asc",
	})
	if err != nil {
		slog.Error("Failed to load categories", "error", err)
		return api.Error("Failed to load categories").ToString()
	}

	categoryList := []map[string]any{}
	for _, term := range terms {
		categoryList = append(categoryList, map[string]any{
			"id":          term.GetID(),
			"name":        term.GetName(),
			"slug":        term.GetSlug(),
			"description": term.GetDescription(),
			"count":       term.GetCount(),
		})
	}

	return api.SuccessWithData("Categories loaded successfully", map[string]any{
		"categories": categoryList,
	}).ToString()
}

func (controller *categoryManagerController) handleCreateCategory(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	var reqData struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.Name == "" {
		return api.Error("Category name is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	categoryTaxonomy, err := controller.ensureTaxonomy(ctx, blogStore)
	if err != nil {
		return api.Error("Failed to ensure taxonomy: " + err.Error()).ToString()
	}

	slug := reqData.Slug
	if slug == "" {
		slug = shared.Slugify(reqData.Name)
	}

	term := blogstore.NewTerm()
	term.SetName(reqData.Name)
	term.SetSlug(slug)
	term.SetDescription(reqData.Description)
	term.SetTaxonomyID(categoryTaxonomy.GetID())

	if err := blogStore.TermCreate(ctx, term); err != nil {
		slog.Error("Failed to create category", "error", err)
		return api.Error("Failed to create category").ToString()
	}

	return api.SuccessWithData("Category created successfully", map[string]any{
		"id":   term.GetID(),
		"name": term.GetName(),
		"slug": term.GetSlug(),
	}).ToString()
}

func (controller *categoryManagerController) handleUpdateCategory(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	categoryID := r.URL.Query().Get("category_id")
	if categoryID == "" {
		return api.Error("Category ID is required").ToString()
	}

	var reqData struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.Name == "" {
		return api.Error("Category name is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	term, err := blogStore.TermFindByID(ctx, categoryID)
	if err != nil {
		return api.Error("Category not found").ToString()
	}

	slug := reqData.Slug
	if slug == "" {
		slug = shared.Slugify(reqData.Name)
	}

	term.SetName(reqData.Name)
	term.SetSlug(slug)
	term.SetDescription(reqData.Description)

	if err := blogStore.TermUpdate(ctx, term); err != nil {
		slog.Error("Failed to update category", "error", err)
		return api.Error("Failed to update category").ToString()
	}

	return api.SuccessWithData("Category updated successfully", map[string]any{
		"id":   term.GetID(),
		"name": term.GetName(),
		"slug": term.GetSlug(),
	}).ToString()
}

func (controller *categoryManagerController) handleReorderCategories(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	var reqData struct {
		CategoryIDs []string `json:"category_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	// Update each category's sequence based on the new order
	for seq, categoryID := range reqData.CategoryIDs {
		term, err := blogStore.TermFindByID(ctx, categoryID)
		if err != nil {
			slog.Error("Failed to find category for reorder", "category_id", categoryID, "error", err)
			continue
		}
		term.SetSequence(seq)
		if err := blogStore.TermUpdate(ctx, term); err != nil {
			slog.Error("Failed to update category sequence", "category_id", categoryID, "error", err)
			return api.Error("Failed to save category order").ToString()
		}
	}

	slog.Info("Categories reordered successfully", "count", len(reqData.CategoryIDs))
	return api.Success("Categories reordered successfully").ToString()
}

func (controller *categoryManagerController) handleDeleteCategory(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	var reqData struct {
		CategoryID string `json:"category_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		reqData.CategoryID = r.FormValue("category_id")
	}

	if reqData.CategoryID == "" {
		return api.Error("Category ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	term, err := blogStore.TermFindByID(ctx, reqData.CategoryID)
	if err != nil {
		slog.Error("Failed to find category for delete", "error", err)
		return api.Error("Category not found").ToString()
	}

	if err := blogStore.TermDelete(ctx, term); err != nil {
		slog.Error("Failed to delete category", "error", err)
		return api.Error("Failed to delete category").ToString()
	}

	return api.SuccessWithData("Category deleted successfully", map[string]any{}).ToString()
}

func (controller *categoryManagerController) ensureTaxonomy(ctx context.Context, store blogstore.StoreInterface) (blogstore.TaxonomyInterface, error) {
	categoryTaxonomy, err := store.TaxonomyFindBySlug(ctx, blogstore.TAXONOMY_CATEGORY)
	if err != nil || categoryTaxonomy == nil {
		controller.registry.GetLogger().Info("Creating category taxonomy")
		categoryTaxonomy = blogstore.NewTaxonomy()
		categoryTaxonomy.SetName("Category")
		categoryTaxonomy.SetSlug(blogstore.TAXONOMY_CATEGORY)
		categoryTaxonomy.SetDescription("Blog post categories")
		if err := store.TaxonomyCreate(ctx, categoryTaxonomy); err != nil {
			return nil, err
		}
	}

	if categoryTaxonomy == nil {
		return nil, errors.New("category taxonomy is nil after ensure")
	}

	return categoryTaxonomy, nil
}

// Deprecated: kept for backwards compatibility
type categoryManagerControllerData struct {
	page          string
	pageInt       int
	perPage       int
	taxonomyID    string
	categoryCount int64
	categoryList  []blogstore.TermInterface
}
