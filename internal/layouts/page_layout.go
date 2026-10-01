package layouts

import (
	"net/http"
	"project/internal/app"
	"project/internal/links"

	baselayouts "github.com/dracory/base/layouts"

	basesession "github.com/dracory/base/session"

	"github.com/dracory/cdn"
	"github.com/dracory/cmsstore"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

// NavbarFunc renders a navbar for the page layout. It receives the app
// instance and the current request so it can adapt to the auth state.
type NavbarFunc func(app app.AppInterface, r *http.Request) hb.TagInterface

// PageLayoutConfig customizes the shared page scaffold for section-specific
// layouts (e.g. the admin layout supplies its own navbar).
type PageLayoutConfig struct {
	// NavbarFn renders the navbar. When nil, the default auth-aware
	// pageNavbar is used. Ignored when Options.DisableNavbar is set.
	NavbarFn NavbarFunc

	// IsAdmin marks the page as an admin page. It suppresses the public
	// site footer and adjusts the title postfix.
	IsAdmin bool
}

// pageLayout is the unified layout for all public and user-facing pages.
// It renders the complete HTML page via hb.Webpage and adapts its navbar
// based on auth state. Section layouts (e.g. admin) reuse it through
// NewPageLayoutWithConfig.
type pageLayout struct {
	app             app.AppInterface
	request         *http.Request
	title           string
	metaDescription string
	metaKeywords    string
	canonicalURL    string
	imageURL        string
	content         hb.TagInterface
	scriptURLs      []string
	scripts         []string
	styleURLs       []string
	styles          []string
	disableNavbar   bool
	isAdmin         bool
	navbarFn        NavbarFunc
	theme           string
}

// NewPageLayout creates a unified page layout. The navbar is shown by
// default and adapts to the user's auth state.
// Set options.DisableNavbar to true for pages that provide their own
// navigation.
func NewPageLayout(app app.AppInterface, r *http.Request, options baselayouts.Options) baselayouts.LayoutInterface {
	return NewPageLayoutWithConfig(app, r, options, PageLayoutConfig{})
}

// NewPageLayoutWithConfig creates a page layout with section-specific
// overrides (custom navbar, admin chrome). Section layout packages such as
// internal/controllers/admin/layout use this to wrap the shared scaffold.
func NewPageLayoutWithConfig(
	app app.AppInterface,
	r *http.Request,
	options baselayouts.Options,
	config PageLayoutConfig,
) baselayouts.LayoutInterface {
	authUser := basesession.GetAuthUser(r)

	section := lo.Ternary(config.IsAdmin, "Admin", lo.Ternary(authUser == nil, "Guest", "User"))
	titlePostfix := " | " + section
	if app.GetConfig().GetAppName() != "" {
		titlePostfix += " | " + app.GetConfig().GetAppName()
	}

	if r != nil {
		// CMS pages carry a cmsstore.PageInterface in the context and
		// supply their own full title
		if page, ok := r.Context().Value("page").(cmsstore.PageInterface); ok && page != nil {
			titlePostfix = ""
		}
	}

	navbarFn := config.NavbarFn
	if navbarFn == nil {
		navbarFn = pageNavbar
	}

	return &pageLayout{
		app:             app,
		request:         r,
		title:           options.Title + titlePostfix,
		metaDescription: options.MetaDescription,
		metaKeywords:    options.MetaKeywords,
		canonicalURL:    options.CanonicalURL,
		imageURL:        options.ImageURL,
		content:         options.Content,
		scriptURLs:      options.ScriptURLs,
		scripts:         options.Scripts,
		styleURLs:       options.StyleURLs,
		styles:          options.Styles,
		disableNavbar:   options.DisableNavbar,
		isAdmin:         config.IsAdmin,
		navbarFn:        navbarFn,
		theme:           ThemeName(r),
	}
}

// ToHTML generates the complete HTML page.
func (layout *pageLayout) ToHTML() string {
	styleURLs := append([]string{
		cdn.BootstrapCss_5_3_3(),
		cdn.BootstrapIconsCss_1_13_1(),
		"https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800;900&family=JetBrains+Mono:wght@700;800&display=swap",
		ThemeStyleURL(layout.theme),
	}, layout.styleURLs...)

	scriptURLs := append([]string{
		cdn.BootstrapJs_5_3_3(),
	}, layout.scriptURLs...)

	themeToggleJS := `document.addEventListener('DOMContentLoaded', function() {
		var html = document.documentElement;
		var btn = document.getElementById('themeToggle');
		var icon = document.getElementById('themeIcon');
		var stored = localStorage.getItem('cf-theme');
		function apply(theme) {
			html.setAttribute('data-bs-theme', theme);
			if (icon) icon.className = theme === 'dark' ? 'bi bi-sun' : 'bi bi-moon';
		}
		if (stored) apply(stored);
		if (btn) {
			btn.addEventListener('click', function() {
				var current = html.getAttribute('data-bs-theme') || 'light';
				var next = current === 'dark' ? 'light' : 'dark';
				apply(next);
				localStorage.setItem('cf-theme', next);
			});
		}
	});`

	webpage := hb.Webpage().
		SetTitle(layout.title).
		SetFavicon(FaviconURL()).
		AddStyleURLs(styleURLs).
		AddStyles(layout.styles).
		AddScriptURLs(scriptURLs).
		AddScripts(layout.scripts).
		AddScripts([]string{themeToggleJS})

	if layout.metaDescription != "" {
		webpage.AddMeta(hb.Meta().Attr("name", "description").Attr("content", layout.metaDescription))
	}

	if layout.metaKeywords != "" {
		webpage.AddMeta(hb.Meta().Attr("name", "keywords").Attr("content", layout.metaKeywords))
	}

	// Canonical link tag to prevent duplicate-content issues for crawlers.
	if layout.canonicalURL != "" {
		webpage.Head().AddChild(hb.NewTag("link").Attr("rel", "canonical").Attr("href", layout.canonicalURL))
	}

	// Open Graph tags for social media link previews.
	if layout.canonicalURL != "" {
		webpage.Head().AddChild(hb.NewTag("meta").Attr("property", "og:url").Attr("content", layout.canonicalURL))
	}
	if layout.title != "" {
		webpage.Head().AddChild(hb.NewTag("meta").Attr("property", "og:title").Attr("content", layout.title))
	}
	if layout.metaDescription != "" {
		webpage.Head().AddChild(hb.NewTag("meta").Attr("property", "og:description").Attr("content", layout.metaDescription))
	}
	if layout.imageURL != "" {
		webpage.Head().AddChild(hb.NewTag("meta").Attr("property", "og:image").Attr("content", layout.imageURL))
	}
	if layout.app != nil && layout.app.GetConfig() != nil && layout.app.GetConfig().GetAppName() != "" {
		webpage.Head().AddChild(hb.NewTag("meta").Attr("property", "og:site_name").Attr("content", layout.app.GetConfig().GetAppName()))
	}
	webpage.Head().AddChild(hb.NewTag("meta").Attr("property", "og:type").Attr("content", "website"))

	if !layout.disableNavbar && layout.navbarFn != nil {
		webpage.AddChild(layout.navbarFn(layout.app, layout.request))
	}

	contentArea := hb.Main().
		Class("cf-main").
		Child(hb.Div().
			Class("container").
			Style("max-width: 1200px;").
			Child(layout.content))

	webpage.AddChild(contentArea)

	// Admin pages manage their own chrome; the public footer is skipped.
	if !layout.isAdmin {
		webpage.AddChild(layout.footer())
	}

	if layout.disableNavbar {
		floatingToggle := hb.Div().
			Style("position: fixed; bottom: 1.5rem; right: 1.5rem; z-index: 1050;").
			Child(hb.Button().
				Type("button").
				ID("themeToggle").
				Class("btn btn-outline-secondary rounded-4 px-3 py-2 fw-black").
				Attr("title", "Toggle theme").
				Style("box-shadow: var(--cf-shadow);").
				Child(hb.I().Class("bi bi-moon").ID("themeIcon")))
		webpage.AddChild(floatingToggle)
	}

	return webpage.ToHTML()
}

// footer renders the generic public site footer.
func (layout *pageLayout) footer() hb.TagInterface {
	appName := ""
	if layout.app != nil && layout.app.GetConfig() != nil {
		appName = layout.app.GetConfig().GetAppName()
	}

	return hb.Footer().
		Class("cf-footer py-4 mt-auto").
		Child(hb.Div().
			Class("container text-center").
			Child(hb.Div().
				Class("d-flex justify-content-center flex-wrap gap-3 mb-2").
				Child(footerLink(links.Website().Home(), "Home")).
				Child(footerLink(links.Website().Blog(), "Blog")).
				Child(footerLink(links.Website().Shop(), "Shop")).
				Child(footerLink(links.Website().Contact(), "Contact"))).
			Child(hb.Small().
				Class("text-secondary fw-bold text-uppercase tracking-wider d-block mb-1").
				Text(appName)).
			Child(hb.Small().
				Class("text-secondary text-uppercase tracking-wider").
				Text("All rights reserved")))
}

// footerLink builds a single footer navigation link.
func footerLink(url string, title string) hb.TagInterface {
	return hb.Hyperlink().
		Href(url).
		Class("text-secondary fw-bold text-uppercase tracking-wider small text-decoration-none").
		Text(title)
}
