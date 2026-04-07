package metadatacomponent

import (
	"embed"
	"log/slog"

	"project/internal/controllers/admin/shop/shared"
	"project/internal/registry"

	"github.com/dracory/hb"
)

//go:embed *.html
//go:embed *.js
var metadataFiles embed.FS

// Render renders the metadata component HTML for the given product
func Render(registry registry.RegistryInterface, productID string) hb.TagInterface {
	htmlContent, err := metadataFiles.ReadFile("metadata.html")
	if err != nil {
		slog.Error("Failed to read metadata HTML template", "error", err)
		return hb.Div().HTML("Error loading metadata component")
	}

	jsContent, err := metadataFiles.ReadFile("metadata.js")
	if err != nil {
		slog.Error("Failed to read metadata JavaScript file", "error", err)
		return hb.Div().HTML("Error loading metadata component")
	}

	vueContainer := hb.Div().ID("metadata-wrapper")

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const productId = '` + productID + `';
		const urlMetasLoad = '` + shared.NewLinks().ProductUpdate(map[string]string{"product_id": productID, "action": "load-metadata"}) + `&X-Requested-With=XMLHttpRequest';
		const urlMetasSave = '` + shared.NewLinks().ProductUpdate(map[string]string{"product_id": productID, "action": "save-metadata"}) + `&X-Requested-With=XMLHttpRequest';
	`)

	componentScript := hb.Script(string(jsContent))

	htmlTemplate := hb.Div().HTML(string(htmlContent))

	vueContainer.
		Child(vueCDN).
		Child(initScript).
		Child(componentScript).
		Child(htmlTemplate)

	formWrapper := hb.Div().
		ID("FormProductMetadataUpdate").
		Child(vueContainer)

	return formWrapper
}
