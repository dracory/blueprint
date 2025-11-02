package layouts

import (
	"project/internal/links"

	"github.com/dracory/dashboard/types"
	"github.com/dracory/userstore"
)

// userDashboardUserMenu generates the user menu items for the dashboard.
//
// Parameters:
// - `authUser` (*models.User): The authenticated user.
//
// Returns:
// - `[]types.MenuItem`: The user menu items.
func adminLayoutUserMenu(authUser userstore.UserInterface) []types.MenuItem {
	userDashboardMenuItem := types.MenuItem{
		Title: "To User Panel",
		URL:   links.User().Home(),
	}

	logoutMenuItem := types.MenuItem{
		Title: "Logout",
		URL:   links.Auth().Logout(),
	}

	items := []types.MenuItem{}

	if authUser != nil && authUser.IsAdministrator() {
		items = append(items, userDashboardMenuItem)
	}

	items = append(items, logoutMenuItem)

	return items
}
