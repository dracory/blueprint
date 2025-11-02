package layouts

import (
	"project/internal/links"

	dashboardTypes "github.com/dracory/dashboard/types"
	"github.com/dracory/hb"
	"github.com/dracory/userstore"
)

func adminLayoutMainMenu(user userstore.UserInterface) []dashboardTypes.MenuItem {
	websiteHomeLink := links.Website().Home()
	dashboardLink := links.Admin().Home()
	loginLink := links.Auth().Login(dashboardLink)
	logoutLink := links.Auth().Logout()

	homeMenuItem := dashboardTypes.MenuItem{
		Icon:  hb.I().Class("bi bi-house").Style("margin-right:10px;").ToHTML(),
		Title: "Home",
		URL:   websiteHomeLink,
	}

	loginMenuItem := dashboardTypes.MenuItem{
		Icon:  hb.I().Class("bi bi-arrow-right").Style("margin-right:10px;").ToHTML(),
		Title: "Login",
		URL:   loginLink,
	}

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
		menuItems = append(menuItems, websiteMenuItem)
		menuItems = append(menuItems, logoutMenuItem)
	} else {
		menuItems = append(menuItems, homeMenuItem)
		menuItems = append(menuItems, loginMenuItem)
	}

	return menuItems
}
