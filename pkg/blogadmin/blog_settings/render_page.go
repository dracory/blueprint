package blog_settings

import (
	"log/slog"
	"net/http"

	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/pkg/blogadmin/shared"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

func (controller *blogSettingsController) renderPage(w http.ResponseWriter, r *http.Request) string {
	authUser := helpers.GetAuthUser(r)
	if authUser == nil {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, "You are not logged in. Please login to continue.", shared.NewLinks("/admin/blog").Home(), 10)
	}

	htmlContent, err := settingsFiles.ReadFile("settings.html")
	if err != nil {
		slog.Error("Failed to read settings HTML", "error", err)
		return hb.Div().HTML("Error loading settings component").ToHTML()
	}

	jsContent, err := settingsFiles.ReadFile("settings.js")
	if err != nil {
		slog.Error("Failed to read settings JS", "error", err)
		return hb.Div().HTML("Error loading settings component").ToHTML()
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		window.blogSettingsReturnUrl = '` + shared.NewLinks("/admin/blog").PostManager() + `';
		const urlBlogSettingsFetchData = '` + shared.NewLinks("/admin/blog").BlogSettings() + `?action=` + actionFetchData + `';
		const urlBlogSettingsSubmit = '` + shared.NewLinks("/admin/blog").BlogSettings() + `?action=` + actionSubmit + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	vueContainer := hb.Div().
		Child(vueCDN).
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Dashboard", URL: links.Admin().Home()},
		{Name: "Blog", URL: links.Admin().Blog()},
		{Name: "Settings", URL: shared.NewLinks("/admin/blog").BlogSettings()},
	})

	heading := hb.Heading1().HTML("Blog Settings")

	buttonBack := hb.Hyperlink().
		Class("btn btn-secondary ms-3").
		HTML("Back to Blog").
		Href(shared.NewLinks("/admin/blog").Home())

	cardBody := hb.Div().
		Class("card-body").
		Child(vueContainer)

	card := hb.Div().
		Class("card shadow-sm").
		Child(hb.Div().
			Class("card-header d-flex justify-content-between align-items-center").
			Child(hb.Heading4().Class("mb-0").HTML("General Settings"))).
		Child(cardBody)

	page := hb.Div().
		Class("container py-4 min-vh-100").
		Child(breadcrumbs).
		Child(hb.Div().
			Class("d-flex align-items-center mb-4").
			Child(heading).
			Child(buttonBack)).
		Child(card)

	return layouts.NewAdminLayout(controller.app, r, layouts.Options{
		Title:   "Settings | Blog",
		Content: page,
		ScriptURLs: []string{
			cdn.Sweetalert2_11(),
		},
	}).ToHTML()
}
