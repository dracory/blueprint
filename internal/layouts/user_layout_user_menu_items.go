package layouts

import (
	"context"
	"project/internal/app"
	"project/internal/helpers"
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
func userLayoutUserMenuItems(application app.AppInterface, ctx context.Context, authUser userstore.UserInterface) []dashboardTypes.MenuItem {
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
		if helpers.UserHasAnyActiveRole(ctx, application, authUser, userstore.USER_ROLE_ADMINISTRATOR, userstore.USER_ROLE_SUPERUSER) {
			items = append(items, adminDashboardMenuItem)
		}
	}

	items = append(items, logoutMenuItem)

	return items
}
