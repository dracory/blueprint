package user_manager

import (
	_ "embed"
	"net/http"
	"strings"

	"project/internal/layouts"
	"project/internal/links"
	"project/pkg/useradmin/shared"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

var (
	//go:embed users.html
	usersHTML string

	//go:embed users.js
	usersJS string
)

func (controller *userManagerController) renderPage(w http.ResponseWriter, r *http.Request) string {
	if controller.app == nil {
		http.Error(w, "app not initialized", http.StatusInternalServerError)
		return ""
	}

	urlUsersLoad := shared.NewLinks("/admin/users").UserManager(map[string]string{"action": actionLoadUsers})
	urlUserDelete := shared.NewLinks("/admin/users").UserManager(map[string]string{"action": actionDeleteUser})
	urlUserCreate := shared.NewLinks("/admin/users").UserManager(map[string]string{"action": actionCreateUser})
	urlUserUpdate := shared.NewLinks("/admin/users").UserUpdate(map[string]string{"user_id": "USER_ID_PLACEHOLDER"})
	urlUserImpersonate := shared.NewLinks("/admin/users").UserImpersonate(map[string]string{"user_id": "USER_ID_PLACEHOLDER"})

	html := strings.ReplaceAll(usersHTML, "urlUsersLoad", "'"+urlUsersLoad+"'")
	html = strings.ReplaceAll(html, "urlUserUpdate", "'"+urlUserUpdate+"'")
	html = strings.ReplaceAll(html, "urlUserImpersonate", "'"+urlUserImpersonate+"'")
	js := strings.ReplaceAll(usersJS, "urlUsersLoad", "'"+urlUsersLoad+"'")
	js = strings.ReplaceAll(js, "urlUserDelete", "'"+urlUserDelete+"'")
	js = strings.ReplaceAll(js, "urlUserCreate", "'"+urlUserCreate+"'")
	js = strings.ReplaceAll(js, "urlUserUpdate", "'"+urlUserUpdate+"'")
	js = strings.ReplaceAll(js, "urlUserImpersonate", "'"+urlUserImpersonate+"'")

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home(map[string]string{})},
		{Name: "User Manager", URL: shared.NewLinks("/admin/users").UserManager()},
	})

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	content := hb.Div().
		Child(vueCDN).
		Child(hb.Raw(html)).
		Child(hb.Script(js))

	return layouts.NewAdminLayout(controller.app, r, layouts.Options{
		Title:   "Users | User Manager",
		Content: layouts.AdminPage(breadcrumbs, content),
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}
