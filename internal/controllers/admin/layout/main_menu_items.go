package layout

import (
	"project/internal/layouts"
	"project/internal/links"

	"github.com/dracory/hb"
	"github.com/dracory/userstore"
)

// mainMenuItems generates the main navigation items for the admin navbar.
// The set differs for guests and authenticated users.
func mainMenuItems(user userstore.UserInterface) []layouts.MenuItem {
	websiteHomeLink := links.Website().Home()
	dashboardLink := links.Admin().Home()
	loginLink := links.Auth().Login(dashboardLink)
	logoutLink := links.Auth().Logout()

	homeMenuItem := layouts.MenuItem{
		Icon:  hb.I().Class("bi bi-house").Style("margin-right:10px;").ToHTML(),
		Title: "Home",
		URL:   websiteHomeLink,
	}

	loginMenuItem := layouts.MenuItem{
		Icon:  hb.I().Class("bi bi-arrow-right").Style("margin-right:10px;").ToHTML(),
		Title: "Login",
		URL:   loginLink,
	}

	websiteMenuItem := layouts.MenuItem{
		Icon:   hb.I().Class("bi bi-globe").Style("margin-right:10px;").ToHTML(),
		Title:  "To Website",
		URL:    websiteHomeLink,
		Target: "_blank",
	}

	logoutMenuItem := layouts.MenuItem{
		Icon:  hb.I().Class("bi bi-arrow-right").Style("margin-right:10px;").ToHTML(),
		Title: "Logout",
		URL:   logoutLink,
	}

	dashboardMenuItem := layouts.MenuItem{
		Icon:  hb.I().Class("bi bi-speedometer").Style("margin-right:10px;").ToHTML(),
		Title: "Dashboard",
		URL:   dashboardLink,
	}

	menuItems := []layouts.MenuItem{}

	if user != nil {
		menuItems = append(menuItems, dashboardMenuItem)
		menuItems = append(menuItems, websiteMenuItem)
		menuItems = append(menuItems, logoutMenuItem)
	} else {
		menuItems = append(menuItems, homeMenuItem)
		menuItems = append(menuItems, loginMenuItem)
	}

	return menuItems
}
