package page_not_found

import (
	"net/http"

	"project/internal/app"
	"project/internal/layouts"
	"project/internal/links"

	baselayouts "github.com/dracory/base/layouts"
	"github.com/dracory/hb"
)

// == CONSTRUCTOR =============================================================

func PageNotFoundController(app app.AppInterface) *pageNotFoundController {
	return &pageNotFoundController{app: app}
}

// == CONTROLLER ==============================================================

type pageNotFoundController struct {
	app app.AppInterface
}

// PUBLIC METHODS =============================================================

func (controller *pageNotFoundController) Handler(w http.ResponseWriter, r *http.Request) string {
	w.WriteHeader(http.StatusNotFound)

	return controller.pageHTML(r)
}

// pageHTML renders the 404 page through the shared blank layout so it inherits
// the active Bootlegance theme instead of using bespoke styling.
func (controller *pageNotFoundController) pageHTML(r *http.Request) string {
	icon := hb.Div().
		Class("display-1 text-primary mb-3").
		Child(hb.I().Class("bi bi-compass"))

	code := hb.Div().
		Class("display-1 fw-bold mb-3").
		Text("404")

	title := hb.Heading1().
		Class("h2 fw-bold mb-3").
		Text("Oops! Page Not Found")

	message := hb.Paragraph().
		Class("text-secondary mb-4").
		Text("It looks like you've ventured into uncharted territory. " +
			"The page you're looking for seems to have sailed away!")

	buttons := hb.Div().
		Class("d-flex flex-wrap justify-content-center gap-2").
		Child(hb.Hyperlink().
			Href("javascript:history.back()").
			Class("btn btn-outline-primary").
			Child(hb.I().Class("bi bi-arrow-left me-2")).
			Text("Go Back")).
		Child(hb.Hyperlink().
			Href(links.Website().Home()).
			Class("btn btn-primary").
			Child(hb.I().Class("bi bi-house-door me-2")).
			Text("Back to Home"))

	card := hb.Div().
		Class("card shadow-sm mx-auto").
		Style("max-width: 600px;").
		Child(hb.Div().
			Class("card-body text-center p-5").
			Child(icon).
			Child(code).
			Child(title).
			Child(message).
			Child(buttons))

	content := hb.Div().
		Class("container").
		Style("min-height: 100vh; display: flex; align-items: center; justify-content: center;").
		Child(card)

	appName := ""
	if controller.app != nil && controller.app.GetConfig() != nil {
		appName = controller.app.GetConfig().GetAppName()
	}

	return layouts.NewBlankLayout(controller.app, r, baselayouts.Options{
		Title:   "404 - Page Not Found",
		AppName: appName,
		Content: content,
	}).ToHTML()
}
