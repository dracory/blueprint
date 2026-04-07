package productupdate

import (
	"log/slog"
	"net/http"
	"project/internal/controllers/admin/shop/products/productupdate/detailscomponent"
	"project/internal/controllers/admin/shop/products/productupdate/mediacomponent"
	metadatacomponent "project/internal/controllers/admin/shop/products/productupdate/metadatacomponent"
	"project/internal/controllers/admin/shop/products/productupdate/tagscomponent"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/req"
)

// Action constants for product update operations
const (
	ACTION_LOAD_DETAILS  = "load-details"
	ACTION_SAVE_DETAILS  = "save-details"
	ACTION_LOAD_METADATA = "load-metadata"
	ACTION_SAVE_METADATA = "save-metadata"
	ACTION_LOAD_TAGS     = "load-tags"
	ACTION_SAVE_TAGS     = "save-tags"
	ACTION_LOAD_MEDIA    = "load-media"
	ACTION_SAVE_MEDIA    = "save-media"
)

type productUpdateController struct {
	registry registry.RegistryInterface
}

func NewProductUpdateController(registry registry.RegistryInterface) *productUpdateController {
	return &productUpdateController{registry: registry}
}

func (controller *productUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	productID := req.GetStringTrimmed(r, "product_id")
	view := req.GetStringTrimmedOr(r, "view", "details")
	action := req.GetStringTrimmed(r, "action")

	if productID == "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Product ID is required", links.Admin().Home(), 10)
	}

	product, err := controller.registry.GetShopStore().ProductFindByID(r.Context(), productID)
	if err != nil {
		slog.Error("Error. productUpdateController: ProductFindByID", slog.String("error", err.Error()), slog.String("product_id", productID))
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Product not found", links.Admin().Home(), 10)
	}

	if product == nil {
		slog.Warn("Warning. productUpdateController: ProductFindByID", slog.String("error", "Product not found"), slog.String("product_id", productID))
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Product not found", links.Admin().Home(), 10)
	}

	// Dispatch based on action
	switch action {
	case ACTION_LOAD_DETAILS:
		return detailscomponent.HandleAjaxLoadDetails(controller.registry, productID)
	case ACTION_SAVE_DETAILS:
		return detailscomponent.HandleAjaxSaveDetails(controller.registry, r, productID)
	case ACTION_LOAD_METADATA:
		return metadatacomponent.HandleAjaxLoadMetadata(controller.registry, productID)
	case ACTION_SAVE_METADATA:
		return metadatacomponent.HandleAjaxSaveMetadata(controller.registry, r, productID)
	case ACTION_LOAD_TAGS:
		return tagscomponent.HandleAjaxLoadTags(controller.registry, productID)
	case ACTION_SAVE_TAGS:
		return tagscomponent.HandleAjaxSaveTags(controller.registry, r, productID)
	case ACTION_LOAD_MEDIA:
		return mediacomponent.HandleAjaxLoadMedia(controller.registry, productID)
	case ACTION_SAVE_MEDIA:
		return mediacomponent.HandleAjaxSaveMedia(controller.registry, r, productID)
	default:
		// No action specified - render page
		return controller.handleRenderPage(r, product, view, productID)
	}
}
