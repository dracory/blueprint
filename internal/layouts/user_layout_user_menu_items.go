package layouts

import (
	"project/internal/links"

	"github.com/dracory/userstore"
	"github.com/gouniverse/dashboard"
)

// userLayoutUserMenu generates the user menu items for the dashboard.
//
// Parameters:
// - `authUser` (*models.User): The authenticated user.
//
// Returns:
// - `[]dashboard.MenuItem`: The user menu items.
func userLayoutUserMenuItems(authUser userstore.UserInterface) []dashboard.MenuItem {
	adminDashboardMenuItem := dashboard.MenuItem{
		Title: "To Admin Dashboard",
		URL:   links.Admin().Home(),
	}

	logoutMenuItem := dashboard.MenuItem{
		Title: "Logout",
		URL:   links.Auth().Logout(),
	}

	profileMenuItem := dashboard.MenuItem{
		Title: "My Account",
		URL:   links.NewUserLinks().Profile(map[string]string{}),
	}

	items := []dashboard.MenuItem{profileMenuItem}

	if authUser != nil {
		if authUser.IsAdministrator() || authUser.IsSuperuser() {
			items = append(items, adminDashboardMenuItem)
		}
	}

	items = append(items, logoutMenuItem)

	return items
}
