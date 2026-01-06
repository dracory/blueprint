package layouts

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/cmsstore"
	"github.com/dracory/dashboard"
	dashboardTypes "github.com/dracory/dashboard/types"
	"github.com/samber/lo"
)

func NewAdminLayout(app types.RegistryInterface, r *http.Request, options Options) dashboardTypes.DashboardInterface {
	return adminLayout(app, r, options)
}

// layout generates a dashboard based on the provided request and layout options.
//
// Parameters:
// - r: a pointer to an http.Request object representing the incoming HTTP request.
// - opts: a layoutOptions struct containing the layout options for the dashboard.
//
// Returns:
// - a dashboardTypes.DashboardInterface object representing the generated dashboard.
func adminLayout(app types.RegistryInterface, r *http.Request, options Options) dashboardTypes.DashboardInterface {
	authUser := helpers.GetAuthUser(r)

	dashboardUser := dashboardTypes.User{}
	if authUser != nil {
		firstName, lastName, err := userDisplayNames(app, r, authUser, app.GetConfig().GetVaultStoreKey())
		if err == nil {
			dashboardUser = dashboardTypes.User{
				FirstName: firstName,
				LastName:  lastName,
			}
		} else {
			dashboardUser = dashboardTypes.User{
				FirstName: "n/a",
				LastName:  "",
			}
		}
	}

	// Prepare script URLs
	scriptURLs := []string{} // prepend any if required
	scriptURLs = append(scriptURLs, options.ScriptURLs...)
	scriptURLs = lo.Uniq(scriptURLs)

	// Prepare scripts
	scripts := []string{} // prepend any if required
	scripts = append(scripts, options.Scripts...)

	styleURLs := []string{} // prepend any if required
	styleURLs = append(styleURLs, options.StyleURLs...)
	styleURLs = lo.Uniq(styleURLs)

	// Prepare styles
	styles := []string{ // prepend any if required
		`a.navbar-brand{font-size:18px;}`,
		`nav#Toolbar {border-bottom: 8px solid red;}`,
	}
	styles = append(styles, options.Styles...)

	homeLink := links.Admin().Home()

	titlePostfix := ` | ` + lo.Ternary(authUser == nil, `Guest`, `Admin`)

	if app.GetConfig().GetAppName() != "" {
		titlePostfix += ` | ` + app.GetConfig().GetAppName()
	}

	_, isPage := r.Context().Value("page").(cmsstore.PageInterface)

	if isPage {
		titlePostfix = "" // no postfix for CMS pages
	}

	redirectURL := lo.IfF(r != nil, func() string {
		if r.URL == nil {
			return ""
		}
		if r.URL.RawQuery == "" {
			return r.URL.Path
		}
		return r.URL.Path + "?" + r.URL.RawQuery
	}).ElseF(func() string {
		return ""
	})

	themeLink := links.Website().Theme(map[string]string{"redirect": redirectURL})

	template := dashboard.New()
	template.SetTemplate(dashboard.TEMPLATE_BOOTSTRAP)
	template.SetHTTPRequest(r)
	template.SetContent(options.Content.ToHTML())
	template.SetTitle(options.Title + titlePostfix)
	template.SetLoginURL(links.Auth().Login(homeLink))
	template.SetMenuMainItems(adminLayoutMainMenu(authUser))
	template.SetMenuUserItems(adminLayoutUserMenu(authUser))
	//dashboard.SetLogoImageURL("/media/user/dashboard-logo.jpg")
	//dashboard.SetNavbarBackgroundColorMode("primary")
	template.SetLogoRawHtml(adminLogoHtml())
	template.SetLogoRedirectURL(homeLink)
	template.SetUser(dashboardUser)
	template.SetThemeHandlerUrl(themeLink)
	template.SetScripts(scripts)
	template.SetScriptURLs(scriptURLs)
	template.SetStyles(styles)
	template.SetStyleURLs(styleURLs)
	template.SetFaviconURL(FaviconURL())

	return template
}
