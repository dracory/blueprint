package tagscomponent

import (
	"embed"
	"log/slog"

	"project/internal/controllers/admin/shop/shared"
	"project/internal/registry"

	"github.com/dracory/hb"
)

//go:embed *.html
//go:embed *.js
var tagsFiles embed.FS

// Render renders the tags component HTML for the given product
func Render(registry registry.RegistryInterface, productID string) hb.TagInterface {
	htmlContent, err := tagsFiles.ReadFile("tags.html")
	if err != nil {
		slog.Error("Failed to read tags HTML template", "error", err)
		return hb.Div().HTML("Error loading tags component")
	}

	jsContent, err := tagsFiles.ReadFile("tags.js")
	if err != nil {
		slog.Error("Failed to read tags JavaScript file", "error", err)
		return hb.Div().HTML("Error loading tags component")
	}

	vueContainer := hb.Div().ID("tags-wrapper")

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const productId = '` + productID + `';
		const urlTagsLoad = '` + shared.NewLinks().ProductUpdate(map[string]string{"product_id": productID, "action": "load-tags"}) + `&X-Requested-With=XMLHttpRequest';
		const urlTagsSave = '` + shared.NewLinks().ProductUpdate(map[string]string{"product_id": productID, "action": "save-tags"}) + `&X-Requested-With=XMLHttpRequest';
	`)

	componentScript := hb.Script(string(jsContent))

	htmlTemplate := hb.Div().HTML(string(htmlContent))

	vueContainer.
		Child(vueCDN).
		Child(initScript).
		Child(componentScript).
		Child(htmlTemplate)

	formWrapper := hb.Div().
		ID("FormProductTagsUpdate").
		Child(vueContainer)

	return formWrapper
}
