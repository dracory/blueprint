package admin

import (
	"project/internal/layouts"
	"project/internal/links"

	"github.com/dracory/hb"
)

func (controller *userManagerController) page(data userManagerControllerData) hb.TagInterface {
	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Home",
			URL:  links.Admin().Home(map[string]string{}),
		},
		{
			Name: "Users",
			URL:  links.Admin().UsersUserManager(map[string]string{}),
		},
		{
			Name: "User Manager",
			URL:  links.Admin().UsersUserManager(map[string]string{}),
		},
	})

	buttonUserNew := hb.Button().
		Class("btn btn-primary float-end").
		Child(hb.I().Class("bi bi-plus-circle").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("New User").
		HxGet(links.Admin().UsersUserCreate(map[string]string{})).
		HxTarget("body").
		HxSwap("beforeend")

	title := hb.Heading1().
		HTML("Users. User Manager").
		Child(buttonUserNew)

	return layouts.AdminPage(
		hb.BR(),
		breadcrumbs,
		hb.HR(),
		title,
		controller.tableUsers(data),
	)
}
