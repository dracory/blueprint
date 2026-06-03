package user_update

import (
	_ "embed"
	"log/slog"
	"net/http"
	"strings"

	"project/internal/ext"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/pkg/useradmin/shared"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

var (
	//go:embed form.html
	formHTML string

	//go:embed form.js
	formJS string
)

func (controller *userUpdateController) renderPage(w http.ResponseWriter, r *http.Request) string {
	userID := req.GetStringTrimmed(r, "user_id")

	if controller.registry.GetUserStore() == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "User store is not configured", shared.NewLinks("/admin/users").UserManager(), 10)
	}

	if userID == "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "User ID is required", shared.NewLinks("/admin/users").UserManager(), 10)
	}

	user, err := controller.registry.GetUserStore().UserFindByID(r.Context(), userID)
	if err != nil {
		if controller.registry.GetLogger() != nil {
			controller.registry.GetLogger().Error("renderPage UserFindByID", slog.String("user_id", userID), slog.String("error", err.Error()))
		}
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Error loading user", shared.NewLinks("/admin/users").UserManager(), 10)
	}
	if user == nil {
		if controller.registry.GetLogger() != nil {
			controller.registry.GetLogger().Error("renderPage user not found", slog.String("user_id", userID))
		}
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "User not found", shared.NewLinks("/admin/users").UserManager(), 10)
	}

	firstName := user.GetFirstName()
	lastName := user.GetLastName()
	if controller.registry.GetConfig().GetVaultStoreUsed() && controller.registry.GetVaultStore() != nil {
		firstName, lastName, _, _, _, err = ext.UserUntokenize(r.Context(), controller.registry, controller.registry.GetConfig().GetVaultStoreKey(), user)
		if err != nil {
			if controller.registry.GetLogger() != nil {
				controller.registry.GetLogger().Error("At userUpdateController > renderPage", slog.String("error", err.Error()))
			}
		}
	}

	displayName := strings.TrimSpace(firstName + " " + lastName)
	if displayName == "" {
		displayName = user.GetID()
	}

	returnURL := shared.NewLinks("/admin/users").UserManager()
	urlGetUser := shared.NewLinks("/admin/users").UserUpdate(map[string]string{"action": actionUserFetch, "user_id": userID})
	urlGetTimezones := shared.NewLinks("/admin/users").UserUpdate(map[string]string{"action": actionGetTimezones})
	urlUpdateUser := shared.NewLinks("/admin/users").UserUpdate(map[string]string{"action": actionUserUpdate})

	html := strings.ReplaceAll(formHTML, "USER_ID_PLACEHOLDER", "'"+userID+"'")
	html = strings.ReplaceAll(html, "RETURN_URL_PLACEHOLDER", "'"+returnURL+"'")
	js := strings.ReplaceAll(formJS, "USER_ID_PLACEHOLDER", "'"+userID+"'")
	js = strings.ReplaceAll(js, "RETURN_URL_PLACEHOLDER", "'"+returnURL+"'")
	js = strings.ReplaceAll(js, "urlGetUser", "'"+urlGetUser+"'")
	js = strings.ReplaceAll(js, "urlGetTimezones", "'"+urlGetTimezones+"'")
	js = strings.ReplaceAll(js, "urlUpdateUser", "'"+urlUpdateUser+"'")

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")
	appDiv := hb.Div().ID("app-user-update").Class("mt-3").HTML(html)

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home()},
		{Name: "User Manager", URL: shared.NewLinks("/admin/users").UserManager()},
		{Name: "Edit User", URL: shared.NewLinks("/admin/users").UserUpdate(map[string]string{"user_id": userID})},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.NewLinks("/admin/users").UserManager())

	heading := hb.Heading1().HTML("Edit User").Child(buttonCancel)

	userTitle := hb.Heading2().Class("mb-3").Text("User: ").Text(displayName)

	card := hb.Div().Class("card").Child(
		hb.Div().Class("card-header").Style("display:flex;justify-content:space-between;align-items:center;").
			Child(hb.Heading4().HTML("User Details").Style("margin-bottom:0;display:inline-block;")),
	).Child(
		hb.Div().Class("card-body").Child(vueCDN).Child(appDiv),
	)

	content := layouts.AdminPage(
		breadcrumbs,
		hb.HR(),
		heading,
		userTitle,
		card,
	)

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Edit User | Users",
		Content: content,
		ScriptURLs: []string{
			cdn.Notiflix_3_2_8(),
		},
		StyleURLs: []string{
			cdn.Notiflix_3_2_8_CSS(),
		},
		Scripts: []string{js},
	}).ToHTML()
}
