package home

import (
	"embed"
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
	"github.com/dracory/shopstore"
)

//go:embed *.html
//go:embed *.js
var homeFiles embed.FS

const (
	actionLoadStats = "load-stats"
)

// == CONT ==============================================================

type homeController struct {
	registry       registry.RegistryInterface
	fileManagerURL string
}

// == CONSTRUCTOR ==============================================================

func NewHomeController(registry registry.RegistryInterface, fileManagerURL string) *homeController {
	return &homeController{registry: registry, fileManagerURL: fileManagerURL}
}

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := r.URL.Query().Get("action")

	switch action {
	case actionLoadStats:
		return controller.handleLoadStats(w, r)
	default:
		return controller.renderPage(r)
	}
}

func (controller *homeController) renderPage(r *http.Request) string {
	authUser := helpers.GetAuthUser(r)
	if authUser == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "You are not logged in. Please login to continue.", links.Admin().Home(), 10)
	}

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home()},
		{Name: "Shop", URL: links.Admin().Shop(map[string]string{})},
	})

	heading := hb.Heading1().HTML("Shop Dashboard")

	htmlContent, err := homeFiles.ReadFile("home.html")
	if err != nil {
		slog.Error("Failed to read home HTML template", "error", err)
		return hb.Div().HTML("Error loading home component").ToHTML()
	}

	jsContent, err := homeFiles.ReadFile("home.js")
	if err != nil {
		slog.Error("Failed to read home JavaScript file", "error", err)
		return hb.Div().HTML("Error loading home component").ToHTML()
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	linksHelper := shared.NewLinks("/admin/shop")
	initScript := hb.Script(`
		const urlLoadStats = '` + linksHelper.Home(map[string]string{"action": actionLoadStats}) + `';
		const urlProducts = '` + linksHelper.Products(map[string]string{}) + `';
		const urlCategories = '` + linksHelper.Categories(map[string]string{}) + `';
		const urlDiscounts = '` + linksHelper.Discounts(map[string]string{}) + `';
		const urlOrders = '` + linksHelper.Orders(map[string]string{}) + `';
	`)

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
		Title:   "Shop | Dashboard",
		Content: content,
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *homeController) handleLoadStats(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		return api.Error("Shop store not available").ToString()
	}

	// Get stats
	productCount, err := shopStore.ProductCount(ctx, shopstore.NewProductQuery())
	if err != nil {
		slog.Error("Failed to count products", "error", err)
		productCount = 0
	}

	categoryCount, err := shopStore.CategoryCount(ctx, shopstore.NewCategoryQuery())
	if err != nil {
		slog.Error("Failed to count categories", "error", err)
		categoryCount = 0
	}

	orderCount, err := shopStore.OrderCount(ctx, shopstore.NewOrderQuery())
	if err != nil {
		slog.Error("Failed to count orders", "error", err)
		orderCount = 0
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Stats loaded successfully", map[string]any{
		"product_count":  productCount,
		"category_count": categoryCount,
		"order_count":    orderCount,
	}).ToString()))
	return ""
}
