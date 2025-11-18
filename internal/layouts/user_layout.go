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

func NewUserLayout(app types.AppInterface, r *http.Request, options Options) dashboardTypes.DashboardInterface {
	return userLayout(app, r, options)
}

// layout generates a dashboard based on the provided request and layout options.
//
// Parameters:
// - r: a pointer to an http.Request object representing the incoming HTTP request.
// - opts: a layoutOptions struct containing the layout options for the dashboard.
//
// Returns:
// - a pointer to a dashboard.Dashboard object representing the generated dashboard.
func userLayout(app types.AppInterface, r *http.Request, options Options) dashboardTypes.DashboardInterface {
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

	// googleTagScriptURL := "https://www.googletagmanager.com/gtag/js?id=G-247NHE839P"
	// googleTagScript := `window.dataLayer = window.dataLayer || []; function gtag(){dataLayer.push(arguments);} gtag('js', new Date()); gtag('config', 'G-247NHE839P');`
	// googleAdsScriptURL := "https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=ca-pub-8821108004642146"
	// statcounterScript := `<script type="text/javascript">var sc_project=12939246;var sc_invisible=1;var sc_security="2c1cdc75";</script><script	src="https://www.statcounter.com/counter/counter.js" async></script><noscript><img	src="https://c.statcounter.com/12939246/0/2c1cdc75/1/" alt=""></noscript>`

	// Prepare script URLs
	scriptURLs := []string{} // prepend any if required
	scriptURLs = append(scriptURLs, options.ScriptURLs...)
	scriptURLs = lo.Uniq(scriptURLs)

	// Prepare scripts
	scripts := []string{} // prepend any if required
	scripts = append(scripts, options.Scripts...)

	// Prepare styles
	styles := []string{ // prepend any if required
		`a.navbar-brand{font-size:18px;}`,
		`nav#Toolbar {border-bottom: 8px solid blue;}`,
	}
	styles = append(styles, options.Styles...)

	homeLink := links.User().Home()

	titlePostfix := ` | ` + lo.Ternary(authUser == nil, `Guest`, `User`)

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

	dashboard := dashboard.New()
	dashboard.SetHTTPRequest(r)
	dashboard.SetContent(options.Content.ToHTML())
	dashboard.SetTitle(options.Title + titlePostfix)
	dashboard.SetFaviconURL(FaviconURL())
	dashboard.SetLoginURL(links.Auth().Login(homeLink))
	dashboard.SetMenuMainItems(userLayoutMainMenuItems(authUser))
	dashboard.SetMenuUserItems(userLayoutUserMenuItems(authUser))
	// dashboard.SetMenuQuickAccessItems(userLayoutQuickAccessMenuItems(authUser))
	dashboard.SetNavbarBackgroundColorMode("primary")
	dashboard.SetLogoRawHtml(userLogoHtml())
	dashboard.SetLogoRedirectURL(homeLink)
	dashboard.SetUser(dashboardUser)
	dashboard.SetThemeHandlerUrl(themeLink)
	dashboard.SetScripts(scripts)
	dashboard.SetScriptURLs(scriptURLs)
	dashboard.SetStyles(styles)
	dashboard.SetStyleURLs(options.StyleURLs)

	return dashboard
}
