package detailscomponent

import (
	"embed"
	"log/slog"

	"project/internal/controllers/admin/shop/shared"
	"project/internal/registry"

	"github.com/dracory/hb"
)

//go:embed *.html
//go:embed *.js
var detailsFiles embed.FS

// Render renders the details component HTML for the given product
func Render(registry registry.RegistryInterface, productID string) hb.TagInterface {
	htmlContent, err := detailsFiles.ReadFile("details.html")
	if err != nil {
		slog.Error("Failed to read details HTML template", "error", err)
		return hb.Div().HTML("Error loading details component")
	}

	jsContent, err := detailsFiles.ReadFile("details.js")
	if err != nil {
		slog.Error("Failed to read details JavaScript file", "error", err)
		return hb.Div().HTML("Error loading details component")
	}

	vueContainer := hb.Div().ID("details-wrapper")

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	quillCSS := hb.Style("").HTML(`
		@import url('https://cdn.jsdelivr.net/npm/quill@1.3.7/dist/quill.snow.css');
		.quill-editor-container .ql-editor {
			min-height: 200px;
		}
	`)

	quillJS := hb.Script("").Src("https://cdn.jsdelivr.net/npm/quill@1.3.7/dist/quill.min.js")

	vueQuillJS := hb.Script("").Src("https://unpkg.com/@vueup/vue-quill@latest/dist/vue-quill.global.prod.js")

	initScript := hb.Script(`
		const productId = '` + productID + `';
		const urlDetailsLoad = '` + shared.NewLinks().ProductUpdate(map[string]string{"product_id": productID, "action": "load-details"}) + `&X-Requested-With=XMLHttpRequest';
		const urlDetailsSave = '` + shared.NewLinks().ProductUpdate(map[string]string{"product_id": productID, "action": "save-details"}) + `&X-Requested-With=XMLHttpRequest';
		const { QuillEditor } = window.VueQuill;
	`)

	componentScript := hb.Script(string(jsContent))

	htmlTemplate := hb.Div().HTML(string(htmlContent))

	vueContainer.
		Child(quillCSS).
		Child(vueCDN).
		Child(quillJS).
		Child(vueQuillJS).
		Child(initScript).
		Child(componentScript).
		Child(htmlTemplate)

	formWrapper := hb.Div().
		ID("FormProductDetailsUpdate").
		Child(vueContainer)

	return formWrapper
}
