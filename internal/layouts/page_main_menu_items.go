package layouts

import (
	"project/internal/links"

	"github.com/dracory/hb"
	"github.com/dracory/userstore"
)

// pageMainMenuItems generates the main navigation items for the page navbar.
// The set differs for guests and authenticated users.
func pageMainMenuItems(user userstore.UserInterface) []MenuItem {
	websiteHomeLink := links.Website().Home()
	dashboardLink := links.User().Home(map[string]string{})
	loginLink := links.Auth().Login(dashboardLink)
	logoutLink := links.Auth().Logout()

	homeMenuItem := MenuItem{
		Icon:  hb.I().Class("bi bi-house").Style("margin-right:10px;").ToHTML(),
		Title: "Home",
		URL:   websiteHomeLink,
	}

	loginMenuItem := MenuItem{
		Icon:  hb.I().Class("bi bi-door-open").Style("margin-right:10px;").ToHTML(),
		Title: "Login",
		URL:   loginLink,
	}

	profileMenuItem := MenuItem{
		Icon:  hb.I().Class("bi bi-person").Style("margin-right:10px;").ToHTML(),
		Title: "My Account",
		URL:   links.User().Profile(),
	}

	// shopMenuItem := MenuItem{
	// 	Icon:  hb.I().Class("bi bi-shop").Style("margin-right:10px;").ToHTML(),
	// 	Title: "Your Shop",
	// 	URL:   links.NewUserLinks().Shop(map[string]string{}),
	// }

	// inviteFriendMenuItem := MenuItem{
	// 	Icon:  hb.I().Class("bi bi-people").Style("margin-right:10px;").ToHTML(),
	// 	Title: "Invite a Friend",
	// 	URL:   links.NewUserLinks().InviteFriend(),
	// }

	websiteMenuItem := MenuItem{
		Icon:   hb.I().Class("bi bi-globe").Style("margin-right:10px;").ToHTML(),
		Title:  "To Website",
		URL:    websiteHomeLink,
		Target: "_blank",
	}

	logoutMenuItem := MenuItem{
		Icon:  hb.I().Class("bi bi-arrow-right").Style("margin-right:10px;").ToHTML(),
		Title: "Logout",
		URL:   logoutLink,
	}

	dashboardMenuItem := MenuItem{
		Icon:  hb.I().Class("bi bi-speedometer").Style("margin-right:10px;").ToHTML(),
		Title: "Dashboard",
		URL:   dashboardLink,
	}

	menuItems := []MenuItem{}

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
