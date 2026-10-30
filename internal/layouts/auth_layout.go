package layouts

import (
	"net/http"
	"project/internal/app"

	baselayouts "github.com/dracory/base/layouts"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

// == TYPE ===================================================================

// authLayout is the blank layout plus a Bootlegance theme switcher fixed
// in the top-right corner.
//
// It is used for the auth package's built-in pages (login, register,
// code verify, password restore/reset) so guests can preview and pick a
// theme before signing in. The selection is stored in the same theme
// cookie the rest of the app reads, so the choice persists after login.
type authLayout struct {
	app        app.AppInterface
	r          *http.Request
	title      string
	content    hb.TagInterface
	scriptURLs []string
	scripts    []string
	styleURLs  []string
	styles     []string
}

// == CONSTRUCTOR =============================================================

// NewAuthLayout creates a new auth layout with a theme switcher.
func NewAuthLayout(app app.AppInterface, r *http.Request, options baselayouts.Options) *authLayout {
	layout := &authLayout{}
	layout.app = app
	layout.r = r
	layout.title = options.Title + " | " + options.AppName
	layout.content = options.Content
	layout.scriptURLs = options.ScriptURLs
	layout.scripts = options.Scripts
	layout.styleURLs = options.StyleURLs
	layout.styles = options.Styles
	return layout
}

// == PUBLIC METHODS ==========================================================

// ToHTML generates the HTML for the auth layout
func (layout *authLayout) ToHTML() string {
	layout.styleURLs = append([]string{
		cdn.BootstrapCss_5_3_3(),
		cdn.BootstrapIconsCss_1_10_2(),
		ThemeStyleURL(ThemeName(layout.r)),
	}, layout.styleURLs...)
	// jQuery + Bootstrap JS first so the theme dropdown works and the
	// package's page scripts (jQuery-based form handlers) can rely on them.
	layout.scriptURLs = append([]string{
		cdn.Jquery_3_7_1(),
		cdn.BootstrapJs_5_3_3(),
	}, layout.scriptURLs...)

	themeSwitcher := hb.Div().
		Style("position:fixed;top:1rem;right:1rem;z-index:1050").
		Child(ThemeDropdown(layout.app, layout.r))

	webpage := hb.Webpage().
		SetTitle(layout.title).
		SetFavicon(FaviconURL()).
		AddStyles(layout.styles).
		AddStyleURLs(layout.styleURLs).
		AddScripts(layout.scripts).
		AddScriptURLs(layout.scriptURLs).
		AddChild(themeSwitcher).
		AddChild(layout.content)
	return webpage.ToHTML()
}
