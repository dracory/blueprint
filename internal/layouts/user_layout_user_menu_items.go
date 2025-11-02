package layouts

import (
	"project/internal/links"

	dashboardTypes "github.com/dracory/dashboard/types"
	"github.com/dracory/userstore"
)

// userLayoutUserMenu generates the user menu items for the dashboard.
//
// Parameters:
// - `authUser` (*models.User): The authenticated user.
//
// Returns:
// - `[]dashboard.MenuItem`: The user menu items.
func userLayoutUserMenuItems(authUser userstore.UserInterface) []dashboardTypes.MenuItem {
	adminDashboardMenuItem := dashboardTypes.MenuItem{
		Title: "To Admin Dashboard",
		URL:   links.Admin().Home(),
	}

	logoutMenuItem := dashboardTypes.MenuItem{
		Title: "Logout",
		URL:   links.Auth().Logout(),
	}

	profileMenuItem := dashboardTypes.MenuItem{
		Title: "My Account",
		URL:   links.User().Profile(),
	}

	items := []dashboardTypes.MenuItem{
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
