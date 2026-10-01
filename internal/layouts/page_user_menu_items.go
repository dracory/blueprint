package layouts

import (
	"context"
	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"

	"github.com/dracory/userstore"
)

// pageUserMenuItems generates the account dropdown items for the page navbar.
func pageUserMenuItems(application app.AppInterface, ctx context.Context, authUser userstore.UserInterface) []MenuItem {
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
		if helpers.UserHasAnyActiveRole(ctx, application, authUser, userstore.USER_ROLE_ADMINISTRATOR, userstore.USER_ROLE_SUPERUSER) {
			items = append(items, adminDashboardMenuItem)
		}
	}

	items = append(items, logoutMenuItem)

	return items
}
