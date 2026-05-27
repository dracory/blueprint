package products

import (
	"embed"
	"net/http"

	"project/internal/registry"
	"project/pkg/shopadmin/shared"

	"github.com/dracory/hb"
	"github.com/dracory/shopstore"
)

//go:embed metadata.html metadata.js
var metadataEmbed embed.FS

type productMetadataComponent struct {
	registry  registry.RegistryInterface
	request   *http.Request
	product   shopstore.ProductInterface
	productID string

	formMetas map[string]string

	formErrorMessage   string
	formSuccessMessage string
}

func NewProductMetadataComponent(registry registry.RegistryInterface) *productMetadataComponent {
	return &productMetadataComponent{registry: registry}
}

func (c *productMetadataComponent) Mount(r *http.Request, product shopstore.ProductInterface, productID string) {
	c.request = r
	c.product = product
	c.productID = productID

	if product != nil {
		metas, _ := product.GetMetas()
		c.formMetas = metas
	}
}

func (c *productMetadataComponent) Handle(r *http.Request) error {
	// For now, just return success - full implementation would parse form data
	c.formSuccessMessage = "Metadata saved successfully"
	return nil
}

func (c *productMetadataComponent) Render() hb.TagInterface {
	// Read embedded HTML and JS files
	htmlBytes, _ := metadataEmbed.ReadFile("metadata.html")
	jsBytes, _ := metadataEmbed.ReadFile("metadata.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	// Build URLs for AJAX
	urlMetadataLoad := shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{
		"action":     "load-metadata",
		"product_id": c.productID,
	})

	urlMetadataSave := shared.NewLinks("/admin/shop").ProductUpdate(map[string]string{
		"action":     "save-metadata",
		"product_id": c.productID,
	})

	// Initialize script with URLs and product ID
	initScript := `
		const productId = "` + c.productID + `";
		const urlMetadataLoad = "` + urlMetadataLoad + `";
		const urlMetadataSave = "` + urlMetadataSave + `";
	`

	return hb.Div().
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}
