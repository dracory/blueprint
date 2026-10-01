package home

import (
	"net/http"
	"project/internal/app"
	"project/internal/layouts"
	"project/internal/links"

	baselayouts "github.com/dracory/base/layouts"

	"github.com/dracory/hb"
)

// == CONSTRUCTOR ==============================================================

func NewHomeController(app app.AppInterface) *homeController {
	return &homeController{
		app: app,
	}
}

// == CONTROLLER ===============================================================

type homeController struct {
	app app.AppInterface
}

// == PUBLIC METHODS ===========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	appName := "Dracory Blueprint"
	if controller != nil && controller.app != nil && controller.app.GetConfig() != nil {
		if controller.app.GetConfig().GetAppName() != "" {
			appName = controller.app.GetConfig().GetAppName()
		}
	}

	actions := hb.Div().
		Class("d-flex flex-wrap gap-2 mt-4").
		Child(hb.Hyperlink().Class("btn btn-primary").Href("/swagger").Text("API")).
		Child(hb.Hyperlink().Class("btn btn-outline-secondary").Href(links.Website().Blog()).Text("Blog")).
		Child(hb.Hyperlink().Class("btn btn-outline-secondary").Href(links.User().Home()).Text("Dashboard"))

	itemRoutes := hb.Div().
		Class("col").
		Child(hb.Div().
			Class("card h-100").
			Child(hb.Div().
				Class("card-body").
				Child(hb.Heading3().Class("h6").Text("Routes")).
				Child(hb.Paragraph().
					Class("card-text small text-secondary mb-0").
					HTML("Website routes are defined under <code>internal/controllers/website</code>."))))

	itemConfig := hb.Div().
		Class("col").
		Child(hb.Div().
			Class("card h-100").
			Child(hb.Div().
				Class("card-body").
				Child(hb.Heading3().Class("h6").Text("Config")).
				Child(hb.Paragraph().
					Class("card-text small text-secondary mb-0").
					HTML("Use environment variables / config to set <code>AppName</code>, stores, and integrations."))))

	itemNext := hb.Div().
		Class("col").
		Child(hb.Div().
			Class("card h-100").
			Child(hb.Div().
				Class("card-body").
				Child(hb.Heading3().Class("h6").Text("Next")).
				Child(hb.Paragraph().
					Class("card-text small text-secondary mb-0").
					HTML("Replace this page with your real landing page or enable CMS pages."))))

	grid := hb.Div().
		Class("row row-cols-1 row-cols-md-3 g-3 mt-4").
		Child(itemRoutes).
		Child(itemConfig).
		Child(itemNext)

	page := hb.Div().
		Class("py-5").
		Child(hb.Heading1().Text("Welcome to " + appName)).
		Child(hb.Paragraph().
			Class("lead text-secondary").
			Text("Your application is running. This starter includes routing, controllers, and optional modules you can enable as you build.")).
		Child(actions).
		Child(grid)

	options := baselayouts.Options{
		Title:   "Home",
		AppName: appName,
		Content: page,
	}

	if controller.app != nil && controller.app.GetConfig() != nil && controller.app.GetConfig().GetCmsStoreUsed() {
		return layouts.NewCmsLayout(controller.app, r, options).ToHTML()
	}

	return layouts.NewPageLayout(controller.app, r, options).ToHTML()
}
