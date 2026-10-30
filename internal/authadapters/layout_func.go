package authadapters

import (
	"net/http"

	"project/internal/app"
	"project/internal/layouts"

	"github.com/dracory/auth/types"
	baselayouts "github.com/dracory/base/layouts"
	"github.com/dracory/hb"
)

// NewAuthLayoutFunc renders the auth package's built-in pages inside the
// app's auth layout — the package supplies the card body and required
// assets via options; the app owns the document (title, favicon,
// Bootstrap + active theme + theme switcher). The signature matches the
// shared dracory FuncLayout convention (blogadmin/shopadmin/auth).
func NewAuthLayoutFunc(application app.AppInterface) func(w http.ResponseWriter, r *http.Request, title string, body string, options types.LayoutOptions) string {
	return func(w http.ResponseWriter, r *http.Request, title string, body string, options types.LayoutOptions) string {
		return layouts.NewAuthLayout(application, r, baselayouts.Options{
			Title:      title,
			Content:    hb.Raw(body),
			Styles:     options.Styles,
			StyleURLs:  options.StyleURLs,
			Scripts:    options.Scripts,
			ScriptURLs: options.ScriptURLs,
		}).ToHTML()
	}
}
