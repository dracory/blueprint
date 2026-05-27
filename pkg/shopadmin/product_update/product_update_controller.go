package product_update

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

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
	actionLoadProduct   = "load-product"
	actionUpdateProduct = "update-product"
)

// == CONTROLLER ==============================================================

type productUpdateController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewProductUpdateController(registry registry.RegistryInterface) *productUpdateController {
	return &productUpdateController{registry: registry}
}

func (controller *productUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadProduct:
		return controller.handleLoadProduct(w, r)
	case actionUpdateProduct:
		return controller.handleUpdateProduct(w, r)
	default:
		return controller.renderPage(r)
	}
}

func (controller *productUpdateController) renderPage(r *http.Request) string {
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
	initScript := hb.Script(`
		const urlLoadProduct = '` + linksHelper.ProductUpdate(map[string]string{"action": actionLoadProduct, "product_id": productID}) + `';
		const urlUpdateProduct = '` + linksHelper.ProductUpdate(map[string]string{"action": actionUpdateProduct, "product_id": productID}) + `';
		const productID = '` + productID + `';
	`)

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(initScript).
		Child(hb.Div().ID("app"))

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Update Product | Shop",
		Content: content,
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *productUpdateController) handleLoadProduct(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		return api.Error("Shop store not available").ToString()
	}

	productID := req.GetStringTrimmed(r, "product_id")
	if productID == "" {
		return api.Error("Product ID is required").ToString()
	}

	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil || product == nil {
		slog.Error("Failed to load product", "error", err)
		return api.Error("Product not found").ToString()
	}

	return api.SuccessWithData("Product loaded successfully", map[string]any{
		"product": map[string]any{
			"id":          product.GetID(),
			"title":       product.GetTitle(),
			"description": product.GetDescription(),
			"status":      product.GetStatus(),
			"price":       product.GetPrice(),
		},
	}).ToString()
}

func (controller *productUpdateController) handleUpdateProduct(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		return api.Error("Shop store not available").ToString()
	}

	var reqBody struct {
		ProductID   string  `json:"product_id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Status      string  `json:"status"`
		Price       float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqBody.ProductID == "" {
		return api.Error("Product ID is required").ToString()
	}

	product, err := shopStore.ProductFindByID(ctx, reqBody.ProductID)
	if err != nil || product == nil {
		return api.Error("Product not found").ToString()
	}

	product.SetTitle(reqBody.Title)
	product.SetDescription(reqBody.Description)
	product.SetStatus(reqBody.Status)
	product.SetPrice(strconv.FormatFloat(reqBody.Price, 'f', 2, 64))

	if err := shopStore.ProductUpdate(ctx, product); err != nil {
		slog.Error("Failed to update product", "error", err)
		return api.Error("Failed to update product").ToString()
	}

	return api.Success("Product updated successfully").ToString()
}
