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

func Render(registry registry.RegistryInterface, categoryID string) hb.TagInterface {
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

	initScript := hb.Script(`
		const categoryId = '` + categoryID + `';
		const urlDetailsLoad = '` + shared.NewLinks().CategoryUpdate(map[string]string{"category_id": categoryID, "action": "load-details"}) + `&X-Requested-With=XMLHttpRequest';
		const urlDetailsSave = '` + shared.NewLinks().CategoryUpdate(map[string]string{"category_id": categoryID, "action": "save-details"}) + `&X-Requested-With=XMLHttpRequest';
		const urlCategoriesList = '` + shared.NewLinks().CategoryUpdate(map[string]string{"category_id": categoryID, "action": "list-categories"}) + `&X-Requested-With=XMLHttpRequest';
	`)

	componentScript := hb.Script(string(jsContent))

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))

	vueContainer.
		Child(vueCDN).
		Child(initScript).
		Child(componentScript).
		Child(htmlTemplate)

	formWrapper := hb.Div().
		ID("FormCategoryDetailsUpdate").
		Child(vueContainer)

	return formWrapper
}
