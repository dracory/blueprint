package mediacomponent

import (
	"embed"
	"log/slog"

	"project/internal/controllers/admin/shop/shared"
	"project/internal/registry"

	"github.com/dracory/hb"
)

//go:embed *.html
//go:embed *.js
var mediaFiles embed.FS

// Render renders the media component HTML for the given product
func Render(registry registry.RegistryInterface, productID string) hb.TagInterface {
	htmlContent, err := mediaFiles.ReadFile("media.html")
	if err != nil {
		slog.Error("Failed to read media HTML template", "error", err)
		return hb.Div().HTML("Error loading media component")
	}

	jsContent, err := mediaFiles.ReadFile("media.js")
	if err != nil {
		slog.Error("Failed to read media JavaScript file", "error", err)
		return hb.Div().HTML("Error loading media component")
	}

	vueContainer := hb.Div().ID("media-wrapper")

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const productId = '` + productID + `';
		const urlMediaLoad = '` + shared.NewLinks().ProductUpdate(map[string]string{"product_id": productID, "action": "load-media"}) + `&X-Requested-With=XMLHttpRequest';
		const urlMediaSave = '` + shared.NewLinks().ProductUpdate(map[string]string{"product_id": productID, "action": "save-media"}) + `&X-Requested-With=XMLHttpRequest';
	`)

	componentScript := hb.Script(string(jsContent))

	htmlTemplate := hb.Div().HTML(string(htmlContent))

	vueContainer.
		Child(vueCDN).
		Child(initScript).
		Child(componentScript).
		Child(htmlTemplate)

	formWrapper := hb.Div().
		ID("FormProductMediaUpdate").
		Child(vueContainer)

	return formWrapper
}
