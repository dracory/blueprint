package layouts

import (
	"project/internal/links"

	dashboardTypes "github.com/dracory/dashboard/types"
	"github.com/dracory/hb"
	"github.com/dracory/userstore"
)

// userLayoutMainMenu generates the main menu items for the user dashboard.
//
// Parameters:
// - `user` (*models.User): The authenticated user.
//
// Returns:
// - `[]dashboard.MenuItem`: The main menu items.
func userLayoutMainMenuItems(user userstore.UserInterface) []dashboardTypes.MenuItem {
	websiteHomeLink := links.Website().Home()
	dashboardLink := links.User().Home(map[string]string{})
	loginLink := links.Auth().Login(dashboardLink)
	logoutLink := links.Auth().Logout()

	homeMenuItem := dashboardTypes.MenuItem{
		Icon:  hb.I().Class("bi bi-house").Style("margin-right:10px;").ToHTML(),
		Title: "Home",
		URL:   websiteHomeLink,
	}

	profileMenuItem := dashboardTypes.MenuItem{
		Icon:  hb.I().Class("bi bi-person").Style("margin-right:10px;").ToHTML(),
		Title: "My Account",
		URL:   links.User().Profile(),
	}

	loginMenuItem := dashboardTypes.MenuItem{
		Icon:  hb.I().Class("bi bi-door-open").Style("margin-right:10px;").ToHTML(),
		Title: "Login",
		URL:   loginLink,
	}

	// shopMenuItem := dashboard.MenuItem{
	// 	Icon:  hb.I().Class("bi bi-shop").Style("margin-right:10px;").ToHTML(),
	// 	Title: "Your Shop",
	// 	URL:   links.NewUserLinks().Shop(map[string]string{}),
	// }

	// inviteFriendMenuItem := dashboard.MenuItem{
	// 	Icon:  hb.I().Class("bi bi-people").Style("margin-right:10px;").ToHTML(),
	// 	Title: "Invite a Friend",
	// 	URL:   links.NewUserLinks().InviteFriend(),
	// }

	websiteMenuItem := dashboardTypes.MenuItem{
		Icon:   hb.I().Class("bi bi-globe").Style("margin-right:10px;").ToHTML(),
		Title:  "To Website",
		URL:    websiteHomeLink,
		Target: "_blank",
	}

	logoutMenuItem := dashboardTypes.MenuItem{
		Icon:  hb.I().Class("bi bi-arrow-right").Style("margin-right:10px;").ToHTML(),
		Title: "Logout",
		URL:   logoutLink,
	}

	dashboardMenuItem := dashboardTypes.MenuItem{
		Icon:  hb.I().Class("bi bi-speedometer").Style("margin-right:10px;").ToHTML(),
		Title: "Dashboard",
		URL:   dashboardLink,
	}

	menuItems := []dashboardTypes.MenuItem{}

	if user != nil {
		menuItems = append(menuItems, dashboardMenuItem)
		// menuItems = append(menuItems, shopMenuItem)
		menuItems = append(menuItems, profileMenuItem)
		// menuItems = append(menuItems, inviteFriendMenuItem)
		menuItems = append(menuItems, websiteMenuItem)
		menuItems = append(menuItems, logoutMenuItem)
	} else {
		menuItems = append(menuItems, homeMenuItem)
		menuItems = append(menuItems, loginMenuItem)
	}

	return menuItems
}
