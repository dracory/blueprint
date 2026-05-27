package product_create

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
	"github.com/dracory/shopstore"
	"github.com/dracory/uid"
)

const (
	actionCreateProduct = "create-product"
)

// == CONTROLLER ==============================================================

type productCreateController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewProductCreateController(registry registry.RegistryInterface) *productCreateController {
	return &productCreateController{registry: registry}
}

func (controller *productCreateController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionCreateProduct:
		return controller.handleCreateProduct(w, r)
	default:
		return controller.renderPage(r)
	}
}

func (controller *productCreateController) renderPage(r *http.Request) string {
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
		{Name: "Create Product", URL: ""},
	})

	heading := hb.Heading1().HTML("Create Product")

	linksHelper := shared.NewLinks("/admin/shop")
	initScript := hb.Script(`
		const urlCreateProduct = '` + linksHelper.ProductCreate(map[string]string{"action": actionCreateProduct}) + `';
	`)

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(initScript).
		Child(hb.Div().ID("app"))

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Create Product | Shop",
		Content: content,
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *productCreateController) handleCreateProduct(w http.ResponseWriter, r *http.Request) string {
	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Shop store not available").ToString()
	}

	var reqBody struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Status      string  `json:"status"`
		Price       float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Invalid request body").ToString()
	}

	if reqBody.Title == "" {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Title is required").ToString()
	}

	product := shopstore.NewProduct()
	product.SetID(uid.HumanUid()[:8])
	product.SetTitle(reqBody.Title)
	product.SetDescription(reqBody.Description)
	product.SetStatus(reqBody.Status)
	product.SetPrice(strconv.FormatFloat(reqBody.Price, 'f', 2, 64))

	if err := shopStore.ProductCreate(r.Context(), product); err != nil {
		slog.Error("Failed to create product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Failed to create product").ToString()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Product created successfully", map[string]any{
		"product_id": product.GetID(),
	}).ToString()))
	return ""
}
