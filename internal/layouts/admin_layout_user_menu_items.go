package layouts

import (
	"project/internal/links"

	"github.com/dracory/userstore"
	"github.com/gouniverse/dashboard"
)

// userDashboardUserMenu generates the user menu items for the dashboard.
//
// Parameters:
// - `authUser` (*models.User): The authenticated user.
//
// Returns:
// - `[]dashboard.MenuItem`: The user menu items.
func adminLayoutUserMenu(authUser userstore.UserInterface) []dashboard.MenuItem {
	userDashboardMenuItem := dashboard.MenuItem{
		Title: "To User Panel",
		URL:   links.User().Home(),
	}

	logoutMenuItem := dashboard.MenuItem{
		Title: "Logout",
		URL:   links.Auth().Logout(),
	}

	items := []dashboard.MenuItem{}

	if authUser != nil && authUser.IsAdministrator() {
		items = append(items, userDashboardMenuItem)
	}

	items = append(items, logoutMenuItem)

	return items
}
