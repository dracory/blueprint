package products

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/shopadmin/shared"

	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/shopstore"
)

type productUpdateController struct {
	registry       registry.RegistryInterface
	fileManagerURL string
}

func NewProductUpdateController(registry registry.RegistryInterface, fileManagerURL string) *productUpdateController {
	return &productUpdateController{registry: registry, fileManagerURL: fileManagerURL}
}

func (controller *productUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")
	productID := req.GetStringTrimmed(r, "product_id")
	view := req.GetStringTrimmedOr(r, "view", "details")

	// Handle AJAX actions
	if action != "" {
		switch action {
		case "load-media":
			return controller.handleLoadMedia(w, r, productID)
		case "save-media":
			return controller.handleSaveMedia(w, r, productID)
		case "load-metadata":
			return controller.handleLoadMetadata(w, r, productID)
		case "save-metadata":
			return controller.handleSaveMetadata(w, r, productID)
		case "load-tags":
			return controller.handleLoadTags(w, r, productID)
		case "save-tags":
			return controller.handleSaveTags(w, r, productID)
		case "load-details":
			return controller.handleLoadDetails(w, r, productID)
		case "save-details":
			return controller.handleSaveDetails(w, r, productID)
		}
	}

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

	// Handle POST requests for each view
	if r.Method == http.MethodPost {
		var component interface {
			Mount(*http.Request, shopstore.ProductInterface, string)
			Handle(*http.Request) error
			Render() hb.TagInterface
		}

		switch view {
		case "details":
			component = NewProductDetailsComponent(controller.registry)
		case "metadata":
			component = NewProductMetadataComponent(controller.registry)
		case "media":
			component = NewProductMediaComponent(controller.registry)
		case "tags":
			component = NewProductTagsComponent(controller.registry)
		default:
			component = NewProductDetailsComponent(controller.registry)
		}

		component.Mount(r, product, productID)
		component.Handle(r)
		return component.Render().ToHTML()
	}

	pageContent := controller.page(r, product, view, productID)

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Edit Product | Shop",
		Content: pageContent,
		ScriptURLs: []string{
			cdn.Jquery_3_7_1(),
			"https://cdn.jsdelivr.net/npm/summernote@0.8.18/dist/summernote-lite.min.js",
			cdn.Sweetalert2_10(),
			"https://unpkg.com/vue@3/dist/vue.global.js",
		},
		StyleURLs: []string{
			"https://cdn.jsdelivr.net/npm/summernote@0.8.18/dist/summernote-lite.min.css",
		},
	}).ToHTML()
}

func (controller *productUpdateController) page(r *http.Request, product shopstore.ProductInterface, view string, productID string) hb.TagInterface {
	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Home",
			URL:  links.Admin().Home(),
		},
		{
			Name: "Shop",
			URL:  shared.NewLinks("/admin/shop").Home(map[string]string{}),
		},
		{
			Name: "Product Manager",
			URL:  shared.NewLinks("/admin/shop").Products(map[string]string{}),
		},
		{
			Name: "Edit Product",
			URL:  shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{"product_id": productID}),
		},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.NewLinks("/admin/shop").Products(map[string]string{}))

	heading := hb.Heading1().
		HTML("Shop. Product. Edit Product").
		Child(buttonCancel)

	tabs := bs.NavTabs().
		Class("mb-3").
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "details", "active").
				Href(shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "details",
				})).
				HTML("Details"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "media", "active").
				Href(shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "media",
				})).
				HTML("Media"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "tags", "active").
				Href(shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "tags",
				})).
				HTML("Tags"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "metadata", "active").
				Href(shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "metadata",
				})).
				HTML("Metadata")))

	productTitle := hb.Heading2().
		Class("mb-3").
		Text("Product: ").
		Text(product.GetTitle())

	var body hb.TagInterface

	switch view {
	case "details":
		component := NewProductDetailsComponent(controller.registry)
		component.Mount(r, product, productID)
		body = component.Render()
	case "media":
		component := NewProductMediaComponent(controller.registry)
		component.Mount(r, product, productID)
		body = component.Render()
	case "tags":
		component := NewProductTagsComponent(controller.registry)
		component.Mount(r, product, productID)
		body = component.Render()
	case "metadata":
		component := NewProductMetadataComponent(controller.registry)
		component.Mount(r, product, productID)
		body = component.Render()
	default:
		component := NewProductDetailsComponent(controller.registry)
		component.Mount(r, product, productID)
		body = component.Render()
	}

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Child(hb.Heading4().
					HTMLIf(view == "details", "Product Details").
					HTMLIf(view == "media", "Product Media").
					HTMLIf(view == "tags", "Product Tags").
					HTMLIf(view == "metadata", "Product Metadata").
					Style("margin-bottom:0;display:inline-block;")),
		).
		Child(
			hb.Div().
				Class("card-body").
				Child(body),
		)

	return hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(productTitle).
		Child(tabs).
		Child(card)
}

func (controller *productUpdateController) handleLoadMedia(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Shop store not available"}`))
		return ""
	}

	// Load media for this product
	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(productID)
	mediaQuery.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
	medias, err := shopStore.MediaList(ctx, mediaQuery)
	if err != nil {
		slog.Error("Failed to load media", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Failed to load media"}`))
		return ""
	}

	// Convert to JSON format
	mediaItems := []map[string]any{}
	for i, media := range medias {
		mediaItems = append(mediaItems, map[string]any{
			"id":       media.GetID(),
			"fileName": media.GetTitle(),
			"url":      media.GetURL(),
			"isMain":   i == 0,
			"type":     media.GetType(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	// Build JSON response with actual media items
	jsonBytes, _ := json.Marshal(map[string]any{
		"status":  "success",
		"message": "Media loaded successfully",
		"data": map[string]any{
			"media": mediaItems,
		},
	})
	w.Write(jsonBytes)
	return ""
}

func (controller *productUpdateController) handleSaveMedia(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Shop store not available"}`))
		return ""
	}

	// Parse request body
	var reqBody struct {
		Media []map[string]any `json:"media"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	json.Unmarshal(bodyBytes, &reqBody)

	// Delete existing media for this product
	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(productID)
	existingMedias, err := shopStore.MediaList(ctx, mediaQuery)
	if err != nil {
		slog.Error("Failed to load existing media", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Failed to load existing media"}`))
		return ""
	}

	for _, existingMedia := range existingMedias {
		err := shopStore.MediaDelete(ctx, existingMedia)
		if err != nil {
			slog.Error("Failed to delete media", "error", err)
		}
	}

	// Create new media entries
	for i, item := range reqBody.Media {
		media := shopstore.NewMedia()
		media.SetEntityID(productID)
		media.SetURL(item["url"].(string))
		media.SetTitle(item["fileName"].(string))
		media.SetType(item["type"].(string))
		media.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
		media.SetSequence(i)

		err := shopStore.MediaCreate(ctx, media)
		if err != nil {
			slog.Error("Failed to create media", "error", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success","message":"Media saved successfully"}`))
	return ""
}

func (controller *productUpdateController) handleLoadMetadata(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Shop store not available"}`))
		return ""
	}

	// Load product
	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		slog.Error("Failed to load product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Failed to load product"}`))
		return ""
	}

	if product == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Product not found"}`))
		return ""
	}

	// Get metadata
	metas, _ := product.GetMetas()

	// Convert to array format
	metadataItems := []map[string]any{}
	for key, value := range metas {
		metadataItems = append(metadataItems, map[string]any{
			"id":    key,
			"key":   key,
			"value": value,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.Marshal(map[string]any{
		"status":   "success",
		"message":  "Metadata loaded successfully",
		"metadata": metadataItems,
	})
	w.Write(jsonBytes)
	return ""
}

func (controller *productUpdateController) handleSaveMetadata(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Shop store not available"}`))
		return ""
	}

	// Parse request body
	var reqBody struct {
		Metadata []map[string]any `json:"metadata"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	json.Unmarshal(bodyBytes, &reqBody)

	// Load product
	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		slog.Error("Failed to load product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Failed to load product"}`))
		return ""
	}

	if product == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Product not found"}`))
		return ""
	}

	// Convert metadata array to map
	metas := make(map[string]string)
	for _, item := range reqBody.Metadata {
		if key, ok := item["key"].(string); ok {
			if value, ok := item["value"].(string); ok {
				metas[key] = value
			}
		}
	}

	// Set metadata on product
	product.SetMetas(metas)

	// Save product
	err = shopStore.ProductUpdate(ctx, product)
	if err != nil {
		slog.Error("Failed to save product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Failed to save product"}`))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success","message":"Metadata saved successfully"}`))
	return ""
}

func (controller *productUpdateController) handleLoadTags(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Shop store not available"}`))
		return ""
	}

	// Load product
	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		slog.Error("Failed to load product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Failed to load product"}`))
		return ""
	}

	if product == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Product not found"}`))
		return ""
	}

	// Get metadata
	metas, _ := product.GetMetas()

	// Parse tags from metadata
	var tags []string
	if tagsMeta, exists := metas["tags"]; exists && tagsMeta != "" {
		// Split by comma and trim whitespace
		tagParts := strings.Split(tagsMeta, ",")
		for _, tag := range tagParts {
			trimmed := strings.TrimSpace(tag)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.Marshal(map[string]any{
		"status":  "success",
		"message": "Tags loaded successfully",
		"tags":    tags,
	})
	w.Write(jsonBytes)
	return ""
}

func (controller *productUpdateController) handleSaveTags(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Shop store not available"}`))
		return ""
	}

	// Parse request body
	var reqBody struct {
		Tags []string `json:"tags"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	json.Unmarshal(bodyBytes, &reqBody)

	// Load product
	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		slog.Error("Failed to load product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Failed to load product"}`))
		return ""
	}

	if product == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Product not found"}`))
		return ""
	}

	// Get existing metadata
	metas, _ := product.GetMetas()
	if metas == nil {
		metas = make(map[string]string)
	}

	// Convert tags array to comma-separated string
	tagsString := strings.Join(reqBody.Tags, ",")

	// Update tags in metadata
	if tagsString != "" {
		metas["tags"] = tagsString
	} else {
		delete(metas, "tags")
	}

	// Set metadata on product
	product.SetMetas(metas)

	// Save product
	err = shopStore.ProductUpdate(ctx, product)
	if err != nil {
		slog.Error("Failed to save product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Failed to save product"}`))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success","message":"Tags saved successfully"}`))
	return ""
}

func (controller *productUpdateController) handleLoadDetails(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Shop store not available"}`))
		return ""
	}

	// Load product
	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		slog.Error("Failed to load product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Failed to load product"}`))
		return ""
	}

	if product == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Product not found"}`))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.Marshal(map[string]any{
		"status":  "success",
		"message": "Details loaded successfully",
		"data": map[string]any{
			"status":      product.GetStatus(),
			"title":       product.GetTitle(),
			"description": product.GetDescription(),
			"price":       product.GetPrice(),
			"quantity":    product.GetQuantity(),
			"memo":        product.GetMemo(),
		},
	})
	w.Write(jsonBytes)
	return ""
}

func (controller *productUpdateController) handleSaveDetails(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := controller.registry.GetShopStore()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Shop store not available"}`))
		return ""
	}

	// Parse request body
	var reqBody struct {
		Status      string `json:"status"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Price       string `json:"price"`
		Quantity    string `json:"quantity"`
		Memo        string `json:"memo"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	json.Unmarshal(bodyBytes, &reqBody)

	// Load product
	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		slog.Error("Failed to load product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Failed to load product"}`))
		return ""
	}

	if product == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Product not found"}`))
		return ""
	}

	// Update product fields
	product.SetStatus(reqBody.Status)
	product.SetTitle(reqBody.Title)
	product.SetDescription(reqBody.Description)
	product.SetPrice(reqBody.Price)
	product.SetQuantity(reqBody.Quantity)
	product.SetMemo(reqBody.Memo)

	// Save product
	err = shopStore.ProductUpdate(ctx, product)
	if err != nil {
		slog.Error("Failed to save product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"error","message":"Failed to save product"}`))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success","message":"Details saved successfully"}`))
	return ""
}
