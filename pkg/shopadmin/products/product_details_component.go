package products

import (
	"embed"
	"net/http"

	"project/internal/app"
	"project/pkg/shopadmin/shared"

	"github.com/dracory/hb"
	"github.com/dracory/shopstore"
)

//go:embed details.html details.js
var detailsEmbed embed.FS

type productDetailsComponent struct {
	app  app.AppInterface
	request   *http.Request
	product   shopstore.ProductInterface
	productID string
}

func NewProductDetailsComponent(app app.AppInterface) *productDetailsComponent {
	return &productDetailsComponent{app: app}
}

func (c *productDetailsComponent) Mount(r *http.Request, product shopstore.ProductInterface, productID string) {
	c.request = r
	c.product = product
	c.productID = productID
}

func (c *productDetailsComponent) Handle(r *http.Request) error {
	// AJAX handles saving, this is a no-op
	return nil
}

func (c *productDetailsComponent) Render() hb.TagInterface {
	// Read embedded HTML and JS files
	htmlBytes, _ := detailsEmbed.ReadFile("details.html")
	jsBytes, _ := detailsEmbed.ReadFile("details.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	// Build URLs for AJAX
	urlDetailsLoad := shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{
		"action":     "load-details",
		"product_id": c.productID,
	})

	urlDetailsSave := shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{
		"action":     "save-details",
		"product_id": c.productID,
	})

	// Initialize script with URLs and product ID
	initScript := `
		const productId = "` + c.productID + `";
		const urlDetailsLoad = "` + urlDetailsLoad + `";
		const urlDetailsSave = "` + urlDetailsSave + `";
	`

	return hb.Div().
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}
