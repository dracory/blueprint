package layouts

import (
	"net/http"
	"project/internal/app"
	"project/internal/links"

	"github.com/dracory/hb"
)

// ThemeDropdown renders a Bootlegance theme picker for navbars. Selecting a
// theme hits the /theme route, which stores the choice in a cookie and
// redirects back to the current page.
func ThemeDropdown(app app.AppInterface, r *http.Request) hb.TagInterface {
	current := ThemeName(r)

	redirectURL := ""
	if r != nil && r.URL != nil {
		redirectURL = r.URL.Path
		if r.URL.RawQuery != "" {
			redirectURL += "?" + r.URL.RawQuery
		}
	}

	var items []hb.TagInterface
	for _, theme := range BootleganceThemes {
		link := hb.Hyperlink().
			Href(links.Website().Theme(map[string]string{
				"theme":    theme,
				"redirect": redirectURL,
			})).
			Class("dropdown-item fw-bold py-2").
			ClassIf(theme == current, "active").
			Text(theme)
		items = append(items, hb.LI().Child(link))
	}

	menu := hb.UL().
		Class("dropdown-menu dropdown-menu-end rounded-4 mt-2").
		Children(items)

	toggle := hb.Button().
		Type("button").
		Class("btn btn-outline-secondary rounded-4 px-3 py-2 fw-black dropdown-toggle").
		Attr("data-bs-toggle", "dropdown").
		Attr("aria-expanded", "false").
		Attr("title", "Choose theme").
		Child(hb.I().Class("bi bi-palette"))

	return hb.Div().
		Class("dropdown").
		Child(toggle).
		Child(menu)
}
