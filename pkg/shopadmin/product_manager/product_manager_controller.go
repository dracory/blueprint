package product_manager

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
var productFiles embed.FS

const (
	actionLoadProducts          = "load-products"
	actionProductDelete         = "delete-product"
	actionProductDeleteSelected = "delete-selected"
)

// == CONTROLLER ==============================================================

type productManagerController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewProductManagerController(registry registry.RegistryInterface) *productManagerController {
	return &productManagerController{registry: registry}
}

func (controller *productManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadProducts:
		return controller.handleLoadProducts(w, r)
	case actionProductDelete:
		return controller.handleProductDelete(w, r)
	case actionProductDeleteSelected:
		return controller.handleProductDeleteSelected(w, r)
	default:
		return controller.renderPage(r)
	}
}

func (controller *productManagerController) renderPage(r *http.Request) string {
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
		{Name: "Products", URL: links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_PRODUCTS})},
	})

	heading := hb.Heading1().HTML("Product Manager")

	linksHelper := shared.NewLinks("/admin/shop")
	urlLoadProducts := linksHelper.Products(map[string]string{"action": actionLoadProducts})
	urlProductDelete := linksHelper.Products(map[string]string{"action": actionProductDelete})
	urlProductDeleteSelected := linksHelper.Products(map[string]string{"action": actionProductDeleteSelected})
	urlUpdateProduct := linksHelper.ProductUpdate(map[string]string{})

	initScript := hb.Script(`
		const urlLoadProducts = '` + urlLoadProducts + `';
		const urlProductDelete = '` + urlProductDelete + `';
		const urlProductDeleteSelected = '` + urlProductDeleteSelected + `';
		const urlUpdateProduct = '` + urlUpdateProduct + `';
	`)

	htmlContent, err := productFiles.ReadFile("products.html")
	if err != nil {
		slog.Error("Failed to read products HTML template", "error", err)
		return hb.Div().HTML("Error loading products component").ToHTML()
	}

	jsContent, err := productFiles.ReadFile("products.js")
	if err != nil {
		slog.Error("Failed to read products JavaScript file", "error", err)
		return hb.Div().HTML("Error loading products component").ToHTML()
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
		Title:   "Products | Shop",
		Content: content,
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *productManagerController) handleLoadProducts(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		Page    int    `json:"page"`
		PerPage int    `json:"per_page"`
		SortBy  string `json:"sort_by"`
		Sort    string `json:"sort"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
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

	query := shopstore.NewProductQuery().
		SetOffset(reqBody.Page * reqBody.PerPage).
		SetLimit(reqBody.PerPage).
		SetOrderBy(reqBody.SortBy).
		SetSortDirection(reqBody.Sort)

	products, err := shopStore.ProductList(ctx, query)
	if err != nil {
		slog.Error("Failed to load products", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to load products").ToString()))
		return ""
	}

	total, err := shopStore.ProductCount(ctx, query)
	if err != nil {
		slog.Error("Failed to count products", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to count products").ToString()))
		return ""
	}

	productList := []map[string]any{}
	for _, product := range products {
		productList = append(productList, map[string]any{
			"id":         product.GetID(),
			"title":      product.GetTitle(),
			"status":     product.GetStatus(),
			"price":      product.GetPrice(),
			"created_at": product.GetCreatedAt(),
			"updated_at": product.GetUpdatedAt(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Products loaded successfully", map[string]any{
		"products": productList,
		"total":    total,
		"page":     reqBody.Page,
		"per_page": reqBody.PerPage,
	}).ToString()))
	return ""
}

func (controller *productManagerController) handleProductDelete(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Shop store not available").ToString()
	}

	var reqBody struct {
		ProductID string `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Invalid request body").ToString()
	}

	if reqBody.ProductID == "" {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Product ID is required").ToString()
	}

	if err := shopStore.ProductDeleteByID(ctx, reqBody.ProductID); err != nil {
		slog.Error("Failed to delete product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Failed to delete product").ToString()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Product deleted successfully").ToString()))
	return ""
}

func (controller *productManagerController) handleProductDeleteSelected(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Shop store not available").ToString()
	}

	var reqBody struct {
		BulkProductIDs []string `json:"bulk_product_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Invalid request body").ToString()
	}

	if len(reqBody.BulkProductIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("No product IDs provided").ToString()
	}

	for _, productID := range reqBody.BulkProductIDs {
		if err := shopStore.ProductDeleteByID(ctx, productID); err != nil {
			slog.Error("Failed to delete product", "error", err, "product_id", productID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Products deleted successfully").ToString()))
	return ""
}
