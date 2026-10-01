package layout

import (
	"net/http"
	"project/internal/app"
	"project/internal/links"

	baselayouts "github.com/dracory/base/layouts"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

// NewCrud renders an admin page preloaded with jQuery, jQuery UI and the
// blockarea script used by the legacy CRUD screens.
func NewCrud(app app.AppInterface, r *http.Request, title string, content string, styleURLs []string, style string, jsURLs []string, js string) string {
	jsURLs = append([]string{
		cdn.Jquery_3_7_1(),
		cdn.JqueryUiJs_1_13_1(),
		links.URL("/resources/blockarea_v0200.js", map[string]string{}),
	}, jsURLs...)
	styleURLs = append([]string{
		cdn.JqueryUiCss_1_13_1(),
	}, styleURLs...)
	return New(app, r, baselayouts.Options{
		Title:      title,
		Content:    hb.Raw(content),
		Scripts:    []string{js},
		ScriptURLs: jsURLs,
		StyleURLs:  styleURLs,
		Styles:     []string{style},
	}).ToHTML()
}
