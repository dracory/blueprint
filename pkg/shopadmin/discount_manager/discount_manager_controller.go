package discount_manager

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
var discountFiles embed.FS

const (
	actionLoadDiscounts          = "load-discounts"
	actionDiscountDelete         = "delete-discount"
	actionDiscountDeleteSelected = "delete-selected"
)

// == CONTROLLER ==============================================================

type discountManagerController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewDiscountManagerController(registry registry.RegistryInterface) *discountManagerController {
	return &discountManagerController{registry: registry}
}

func (controller *discountManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadDiscounts:
		return controller.handleLoadDiscounts(w, r)
	case actionDiscountDelete:
		return controller.handleDiscountDelete(w, r)
	case actionDiscountDeleteSelected:
		return controller.handleDiscountDeleteSelected(w, r)
	default:
		return controller.renderPage(r)
	}
}

func (controller *discountManagerController) renderPage(r *http.Request) string {
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
		{Name: "Discounts", URL: links.Admin().Shop(map[string]string{"controller": shared.CONTROLLER_DISCOUNTS})},
	})

	heading := hb.Heading1().HTML("Discount Manager")

	linksHelper := shared.NewLinks("/admin/shop")
	initScript := hb.Script(`
		const urlLoadDiscounts = '` + linksHelper.Discounts(map[string]string{"action": actionLoadDiscounts}) + `';
		const urlDiscountDelete = '` + linksHelper.Discounts(map[string]string{"action": actionDiscountDelete}) + `';
		const urlDiscountDeleteSelected = '` + linksHelper.Discounts(map[string]string{"action": actionDiscountDeleteSelected}) + `';
	`)

	htmlContent, err := discountFiles.ReadFile("discounts.html")
	if err != nil {
		slog.Error("Failed to read discounts HTML template", "error", err)
		return hb.Div().HTML("Error loading discounts component").ToHTML()
	}

	jsContent, err := discountFiles.ReadFile("discounts.js")
	if err != nil {
		slog.Error("Failed to read discounts JavaScript file", "error", err)
		return hb.Div().HTML("Error loading discounts component").ToHTML()
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
		Title:   "Discounts | Shop",
		Content: content,
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *discountManagerController) handleLoadDiscounts(w http.ResponseWriter, r *http.Request) string {
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

	query := shopstore.NewDiscountQuery().
		SetOffset(reqBody.Page * reqBody.PerPage).
		SetLimit(reqBody.PerPage).
		SetOrderBy(reqBody.SortBy).
		SetSortDirection(reqBody.Sort)

	discounts, err := shopStore.DiscountList(ctx, query)
	if err != nil {
		slog.Error("Failed to load discounts", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to load discounts").ToString()))
		return ""
	}

	total, err := shopStore.DiscountCount(ctx, query)
	if err != nil {
		slog.Error("Failed to count discounts", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to count discounts").ToString()))
		return ""
	}

	discountList := []map[string]any{}
	for _, discount := range discounts {
		discountList = append(discountList, map[string]any{
			"id":         discount.GetID(),
			"code":       discount.GetCode(),
			"amount":     discount.GetAmount(),
			"status":     discount.GetStatus(),
			"created_at": discount.GetCreatedAt(),
			"updated_at": discount.GetUpdatedAt(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Discounts loaded successfully", map[string]any{
		"discounts": discountList,
		"total":     total,
		"page":      reqBody.Page,
		"per_page":  reqBody.PerPage,
	}).ToString()))
	return ""
}

func (controller *discountManagerController) handleDiscountDelete(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		DiscountID string `json:"discount_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.DiscountID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Discount ID is required").ToString()))
		return ""
	}

	if err := shopStore.DiscountDeleteByID(ctx, reqBody.DiscountID); err != nil {
		slog.Error("Failed to delete discount", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to delete discount").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Discount deleted successfully").ToString()))
	return ""
}

func (controller *discountManagerController) handleDiscountDeleteSelected(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		BulkDiscountIDs []string `json:"bulk_discount_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if len(reqBody.BulkDiscountIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("No discount IDs provided").ToString()))
		return ""
	}

	for _, discountID := range reqBody.BulkDiscountIDs {
		if err := shopStore.DiscountDeleteByID(ctx, discountID); err != nil {
			slog.Error("Failed to delete discount", "error", err, "discount_id", discountID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Discounts deleted successfully").ToString()))
	return ""
}
