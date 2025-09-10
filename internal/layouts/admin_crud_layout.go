package layouts

import (
	"net/http"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

func NewAdminCrudLayout(app types.AppInterface, r *http.Request, title string, content string, styleURLs []string, style string, jsURLs []string, js string) string {
	return adminCrudLayout(app, r, title, content, styleURLs, style, jsURLs, js)
}

func adminCrudLayout(app types.AppInterface, r *http.Request, title string, content string, styleURLs []string, style string, jsURLs []string, js string) string {
	jsURLs = append([]string{
		cdn.Jquery_3_7_1(),
		cdn.JqueryUiJs_1_13_1(),
		links.URL("/resources/blockarea_v0200.js", map[string]string{}),
	}, jsURLs...)
	styleURLs = append([]string{
		// "https://cdn.datatables.net/1.13.4/css/jquery.dataTables.min.css",
		cdn.JqueryUiCss_1_13_1(),
	}, styleURLs...)
	// cfmt.Infoln(styleURLs)
	dashboard := NewAdminLayout(app, r, Options{
		Title:      title,
		Content:    hb.Raw(content),
		Scripts:    []string{js},
		ScriptURLs: jsURLs,
		StyleURLs:  styleURLs,
		Styles:     []string{style},
	})
	return dashboard.ToHTML()
}
