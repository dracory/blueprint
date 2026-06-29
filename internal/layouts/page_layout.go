package layouts

import (
	_ "embed"
	"net/http"
	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

//go:embed brutalski_theme.css
var brutalskiThemeCSS string

// pageLayout is the unified layout for all public and user-facing pages.
// It replaces the separate blank, user, and CMS layouts with a single
// neo-brutalist page wrapper that adapts its navbar based on auth state.
type pageLayout struct {
	app             app.AppInterface
	request         *http.Request
	title           string
	metaDescription string
	content         hb.TagInterface
	scriptURLs      []string
	scripts         []string
	styleURLs       []string
	styles          []string
	disableNavbar   bool
}

// NewPageLayout creates a unified page layout with neo-brutalist styling.
// The navbar is shown by default and adapts to the user's auth state.
// Set Options.DisableNavbar to true for pages that provide their own navigation.
func NewPageLayout(app app.AppInterface, r *http.Request, options Options) LayoutInterface {
	authUser := helpers.GetAuthUser(r)

	titlePostfix := " | " + lo.Ternary(authUser == nil, "Guest", "User")
	if app.GetConfig().GetAppName() != "" {
		titlePostfix += " | " + app.GetConfig().GetAppName()
	}

	if r != nil {
		_, isPage := r.Context().Value("page").(struct{})
		if isPage {
			titlePostfix = ""
		}
	}

	return &pageLayout{
		app:             app,
		request:         r,
		title:           options.Title + titlePostfix,
		metaDescription: options.MetaDescription,
		content:         options.Content,
		scriptURLs:      options.ScriptURLs,
		scripts:         options.Scripts,
		styleURLs:       options.StyleURLs,
		styles:          options.Styles,
		disableNavbar:   options.DisableNavbar,
	}
}

// ToHTML generates the complete HTML page.
func (layout *pageLayout) ToHTML() string {
	styleURLs := append([]string{
		cdn.BootstrapCss_5_3_3(),
		cdn.BootstrapIconsCss_1_13_1(),
		"https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800;900&family=JetBrains+Mono:wght@700;800&display=swap",
	}, layout.styleURLs...)

	scriptURLs := append([]string{
		cdn.BootstrapJs_5_3_3(),
	}, layout.scriptURLs...)

	themeCSS := brutalskiThemeCSS

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
		AddStyles([]string{themeCSS}).
		AddScriptURLs(scriptURLs).
		AddScripts(layout.scripts).
		AddScripts([]string{themeToggleJS})

	if layout.app != nil && layout.app.GetConfig() != nil && layout.app.GetConfig().IsEnvProduction() {
		ga4ID := "G-4ZHR10WN2P"
		ga4Script := `(function(){
	var ga4Id = '` + ga4ID + `';
	var consent = localStorage.getItem('cf-cookie-consent');
	function loadGA4() {
		var s = document.createElement('script');
		s.async = true;
		s.src = 'https://www.googletagmanager.com/gtag/js?id=' + ga4Id;
		document.head.appendChild(s);
		window.dataLayer = window.dataLayer || [];
		function gtag(){dataLayer.push(arguments);}
		gtag('js', new Date());
		gtag('config', ga4Id);
		window.gtag = gtag;
		var p = new URLSearchParams(window.location.search);
		if (p.get('ga_purchase')) {
			gtag('event', 'purchase', {
				value: parseFloat(p.get('ga_purchase')),
				currency: 'GBP',
				transaction_id: p.get('ga_order') || ''
			});
		}
	}
	if (consent === 'accepted') {
		loadGA4();
	} else if (consent === null) {
		var banner = document.createElement('div');
		banner.id = 'cf-consent-banner';
		banner.style.cssText = 'position:fixed;bottom:0;left:0;right:0;z-index:9999;background:#1e293b;color:#fff;padding:16px 24px;display:flex;align-items:center;justify-content:space-between;gap:16px;font-family:Inter,sans-serif;font-size:14px;box-shadow:0 -4px 6px rgba(0,0,0,0.1);';
		banner.innerHTML = '<span>We use cookies to analyze traffic and improve your experience. By accepting, you agree to our use of analytics cookies.</span>';
		var btns = document.createElement('div');
		btns.style.cssText = 'display:flex;gap:8px;flex-shrink:0;';
		var accept = document.createElement('button');
		accept.textContent = 'Accept';
		accept.style.cssText = 'background:#4f46e5;color:#fff;border:2px solid #1e293b;padding:8px 20px;font-weight:800;cursor:pointer;border-radius:0;text-transform:uppercase;letter-spacing:0.05em;';
		accept.onclick = function() {
			localStorage.setItem('cf-cookie-consent', 'accepted');
			banner.remove();
			loadGA4();
		};
		var decline = document.createElement('button');
		decline.textContent = 'Decline';
		decline.style.cssText = 'background:transparent;color:#fff;border:2px solid #fff;padding:8px 20px;font-weight:800;cursor:pointer;border-radius:0;text-transform:uppercase;letter-spacing:0.05em;';
		decline.onclick = function() {
			localStorage.setItem('cf-cookie-consent', 'declined');
			banner.remove();
		};
		btns.appendChild(accept);
		btns.appendChild(decline);
		banner.appendChild(btns);
		document.addEventListener('DOMContentLoaded', function() {
			document.body.appendChild(banner);
		});
	}
})();`
		webpage.AddScripts([]string{ga4Script})
	}

	if layout.metaDescription != "" {
		webpage.AddChild(hb.Meta().Attr("name", "description").Attr("content", layout.metaDescription))
	}

	if !layout.disableNavbar {
		webpage.AddChild(userLayoutNavbar(layout.app, layout.request))
	}

	contentArea := hb.Main().
		Class("cf-main").
		Child(hb.Div().
			Class("container").
			Style("max-width: 1200px;").
			Child(layout.content))

	footer := hb.Footer().
		Class("cf-footer py-4 mt-auto").
		Child(hb.Div().
			Class("container text-center").
			Child(hb.Div().
				Class("d-flex justify-content-center gap-3 mb-2").
				Child(hb.Hyperlink().
					Href(links.Website().Home()).
					Class("text-secondary fw-bold text-uppercase tracking-wider small text-decoration-none").
					Text("Home")).
				Child(hb.Hyperlink().
					Href(links.Website().Shop()).
					Class("text-secondary fw-bold text-uppercase tracking-wider small text-decoration-none").
					Text("Courses")).
				Child(hb.Hyperlink().
					Href(links.Website().Blog()).
					Class("text-secondary fw-bold text-uppercase tracking-wider small text-decoration-none").
					Text("Blog")).
				Child(hb.Hyperlink().
					Href(links.Website().Contact()).
					Class("text-secondary fw-bold text-uppercase tracking-wider small text-decoration-none").
					Text("Contact"))).
			Child(hb.Div().
				Class("d-flex justify-content-center gap-3 mb-3").
				Child(hb.Hyperlink().
					Href("https://x.com/CourseThreadHQ").
					Attr("target", "_blank").
					Attr("aria-label", "X").
					Class("text-secondary fs-5 text-decoration-none").
					Child(hb.I().Class("bi bi-twitter-x"))).
				Child(hb.Hyperlink().
					Href("https://www.linkedin.com/company/coursethread").
					Attr("target", "_blank").
					Attr("aria-label", "LinkedIn").
					Class("text-secondary fs-5 text-decoration-none").
					Child(hb.I().Class("bi bi-linkedin"))).
				Child(hb.Hyperlink().
					Href("https://www.reddit.com/user/CourseThreadHQ").
					Attr("target", "_blank").
					Attr("aria-label", "Reddit").
					Class("text-secondary fs-5 text-decoration-none").
					Child(hb.I().Class("bi bi-reddit"))).
				Child(hb.Hyperlink().
					Href("https://www.producthunt.com/@coursethreadhq").
					Attr("target", "_blank").
					Attr("aria-label", "Product Hunt").
					Class("text-secondary fw-black text-decoration-none tracking-wider small border border-2 border-secondary rounded px-1").
					Text("PH"))).
			Child(hb.Small().
				Class("text-secondary fw-bold text-uppercase tracking-wider d-block mb-1").
				Text("CourseThread - Interactive Email Learning")).
			Child(hb.Small().
				Class("text-secondary text-uppercase tracking-wider").
				Text("All rights reserved 2026")))

	webpage.
		AddChild(contentArea).
		AddChild(footer)

	if layout.disableNavbar {
		floatingToggle := hb.Div().
			Style("position: fixed; bottom: 1.5rem; right: 1.5rem; z-index: 1050;").
			Child(hb.Button().
				Type("button").
				ID("themeToggle").
				Class("btn btn-outline-dark rounded-4 px-3 py-2 fw-black").
				Attr("title", "Toggle theme").
				Style("box-shadow: var(--cf-shadow);").
				Child(hb.I().Class("bi bi-moon").ID("themeIcon")))
		webpage.AddChild(floatingToggle)
	}

	return webpage.ToHTML()
}
