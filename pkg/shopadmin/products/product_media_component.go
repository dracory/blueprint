package products

import (
	"context"
	"embed"
	"log/slog"
	"net/http"

	"project/internal/registry"
	"project/pkg/shopadmin/shared"

	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/shopstore"
	"github.com/spf13/cast"
)

//go:embed media.html media.js
var mediaEmbed embed.FS

type productMediaComponent struct {
	registry       registry.RegistryInterface
	request        *http.Request
	product        shopstore.ProductInterface
	productID      string
	fileManagerURL string

	formMedias []shopstore.MediaInterface

	formErrorMessage   string
	formSuccessMessage string
}

func NewProductMediaComponent(registry registry.RegistryInterface) *productMediaComponent {
	return &productMediaComponent{registry: registry}
}

func (c *productMediaComponent) Mount(r *http.Request, product shopstore.ProductInterface, productID string) {
	c.request = r
	c.product = product
	c.productID = productID
	c.fileManagerURL = "/admin/files" // Default file manager URL

	// Load media for this product
	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(productID)
	mediaQuery.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
	medias, _ := c.registry.GetShopStore().MediaList(context.Background(), mediaQuery)
	c.formMedias = medias
}

func (c *productMediaComponent) Handle(r *http.Request) error {
	// Get media URLs and types from form
	mediaURLs := req.GetAll(r)["media_url"]
	mediaTypes := req.GetAll(r)["media_type"]

	if len(mediaURLs) != len(mediaTypes) {
		c.formErrorMessage = "Media URLs and types count mismatch"
		return nil
	}

	// Delete existing media for this product
	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(c.productID)
	existingMedias, err := c.registry.GetShopStore().MediaList(context.Background(), mediaQuery)
	if err != nil {
		slog.Error("At productMediaComponent > Handle", slog.String("error", err.Error()))
		c.formErrorMessage = "System error. Loading existing media failed"
		return nil
	}

	for _, existingMedia := range existingMedias {
		err := c.registry.GetShopStore().MediaDelete(context.Background(), existingMedia)
		if err != nil {
			slog.Error("At productMediaComponent > Handle", slog.String("error", err.Error()))
		}
	}

	// Create new media entries
	for i, mediaURL := range mediaURLs {
		if mediaURL == "" {
			continue
		}

		if len(mediaTypes) <= i || mediaTypes[i] == "" {
			continue
		}

		media := shopstore.NewMedia()
		media.SetEntityID(c.productID)
		media.SetURL(mediaURL)
		media.SetType(mediaTypes[i])
		media.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
		media.SetSequence(i)

		err := c.registry.GetShopStore().MediaCreate(context.Background(), media)
		if err != nil {
			slog.Error("At productMediaComponent > Handle", slog.String("error", err.Error()))
			c.formErrorMessage = "System error. Creating media failed"
			return nil
		}
	}

	c.formSuccessMessage = "Media saved successfully"
	return nil
}

func (c *productMediaComponent) Render() hb.TagInterface {
	// Read embedded HTML and JS files
	htmlBytes, _ := mediaEmbed.ReadFile("media.html")
	jsBytes, _ := mediaEmbed.ReadFile("media.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	// Build URLs for AJAX
	urlMediaLoad := shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{
		"action":     "load-media",
		"product_id": c.productID,
	})

	urlMediaSave := shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{
		"action":     "save-media",
		"product_id": c.productID,
	})

	// Initialize script with URLs and product ID
	initScript := `
		const productId = "` + c.productID + `";
		const urlMediaLoad = "` + urlMediaLoad + `";
		const urlMediaSave = "` + urlMediaSave + `";
	`

	return hb.Div().
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}

func (c *productMediaComponent) createMediaItem(media shopstore.MediaInterface, index int) hb.TagInterface {
	mediaType := media.GetType()
	mediaURL := media.GetURL()
	mediaTitle := media.GetTitle()

	icon := hb.I().Class("bi bi-image fs-4 text-info")
	if mediaType == shopstore.MEDIA_TYPE_VIDEO_MP4 {
		icon = hb.I().Class("bi bi-play-circle fs-4 text-primary")
	}

	title := mediaTitle
	if title == "" {
		title = "Media " + cast.ToString(index+1)
	}

	typeBadge := hb.Span().Class("badge bg-info").HTML("Image")
	if mediaType == shopstore.MEDIA_TYPE_VIDEO_MP4 {
		typeBadge = hb.Span().Class("badge bg-primary").HTML("Video")
	}

	titleDiv := hb.Div().Class("fw-bold").HTML(title)
	urlDiv := hb.Div().Class("text-muted small").HTML(mediaURL)

	mediaInfo := hb.Div().
		Class("d-flex align-items-center gap-3").
		Child(icon).
		Child(hb.Div().Child(titleDiv).Child(urlDiv))

	urlInput := hb.Input().
		Type("hidden").
		Name("media_url[" + cast.ToString(index) + "]").
		Value(mediaURL)

	typeInput := hb.Input().
		Type("hidden").
		Name("media_type[" + cast.ToString(index) + "]").
		Value(mediaType)

	deleteButton := hb.Button().
		Type("button").
		Class("btn btn-sm btn-outline-danger w-100").
		OnClick("removeMediaItem(this)").
		Child(hb.I().Class("bi bi-trash"))

	row := hb.Div().
		Class("row align-items-center g-3").
		Child(hb.Div().Class("col-md-8").Child(mediaInfo)).
		Child(hb.Div().Class("col-md-3").Child(typeBadge)).
		Child(hb.Div().Class("col-md-1").Child(deleteButton))

	cardBody := hb.Div().
		Class("card-body p-3").
		Child(row).
		Child(urlInput).
		Child(typeInput)

	card := hb.Div().
		Class("card border-0 shadow-sm").
		Child(cardBody)

	return hb.Div().
		Class("media-item mb-3").
		Child(card)
}
