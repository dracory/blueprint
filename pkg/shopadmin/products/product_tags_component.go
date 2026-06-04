package products

import (
	"embed"
	"net/http"

	"project/internal/app"
	"project/pkg/shopadmin/shared"

	"github.com/dracory/hb"
	"github.com/dracory/shopstore"
)

//go:embed tags.html tags.js
var tagsEmbed embed.FS

type productTagsComponent struct {
	app  app.AppInterface
	request   *http.Request
	product   shopstore.ProductInterface
	productID string
}

func NewProductTagsComponent(app app.AppInterface) *productTagsComponent {
	return &productTagsComponent{app: app}
}

func (c *productTagsComponent) Mount(r *http.Request, product shopstore.ProductInterface, productID string) {
	c.request = r
	c.product = product
	c.productID = productID
}

func (c *productTagsComponent) Handle(r *http.Request) error {
	// Tags are now handled via AJAX, no form handling needed
	return nil
}

func (c *productTagsComponent) Render() hb.TagInterface {
	// Read embedded HTML and JS files
	htmlBytes, _ := tagsEmbed.ReadFile("tags.html")
	jsBytes, _ := tagsEmbed.ReadFile("tags.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	// Build URLs for AJAX
	urlTagsLoad := shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{
		"action":     "load-tags",
		"product_id": c.productID,
	})

	urlTagsSave := shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{
		"action":     "save-tags",
		"product_id": c.productID,
	})

	// Initialize script with URLs and product ID
	initScript := `
		const productId = "` + c.productID + `";
		const urlTagsLoad = "` + urlTagsLoad + `";
		const urlTagsSave = "` + urlTagsSave + `";
	`

	return hb.Div().
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}
