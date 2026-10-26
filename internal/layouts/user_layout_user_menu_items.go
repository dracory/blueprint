package layouts

import (
	"project/internal/links"

	"github.com/dracory/userstore"
)

// userLayoutUserMenu generates the user menu items for the dashboard.
//
// Parameters:
// - `authUser` (*models.User): The authenticated user.
//
// Returns:
// - `[]MenuItem`: The user menu items.
func userLayoutUserMenuItems(authUser userstore.UserInterface) []MenuItem {
	adminDashboardMenuItem := MenuItem{
		Title: "To Admin Dashboard",
		URL:   links.Admin().Home(),
	}

	logoutMenuItem := MenuItem{
		Title: "Logout",
		URL:   links.Auth().Logout(),
	}

	profileMenuItem := MenuItem{
		Title: "My Account",
		URL:   links.User().Profile(),
	}

	items := []MenuItem{
		profileMenuItem,
	}

	if authUser != nil {
		if authUser.IsAdministrator() || authUser.IsSuperuser() {
			items = append(items, adminDashboardMenuItem)
		}
	}

	items = append(items, logoutMenuItem)

	return items
}
