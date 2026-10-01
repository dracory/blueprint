package layout

import (
	"net/http"
	"project/internal/app"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"

	"github.com/dracory/userstore"
)

// userMenuItems generates the account dropdown items for the admin navbar.
func userMenuItems(application app.AppInterface, r *http.Request, authUser userstore.UserInterface) []layouts.MenuItem {
	userDashboardMenuItem := layouts.MenuItem{
		Title: "To User Panel",
		URL:   links.User().Home(),
	}

	logoutMenuItem := layouts.MenuItem{
		Title: "Logout",
		URL:   links.Auth().Logout(),
	}

	items := []layouts.MenuItem{}

	if authUser != nil && helpers.UserHasActiveRole(r.Context(), application, authUser, userstore.USER_ROLE_ADMINISTRATOR) {
		items = append(items, userDashboardMenuItem)
	}

	items = append(items, logoutMenuItem)

	return items
}
