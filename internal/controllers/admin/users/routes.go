package admin

import (
	"project/internal/links"

	"github.com/dracory/rtr"
)

func UserRoutes() []rtr.RouteInterface {
	userCreate := rtr.NewRoute().
		SetName("Admin > Users > User Create").
		SetPath(links.ADMIN_USERS_USER_CREATE).
		SetHTMLHandler(NewUserCreateController().Handler)

	userDelete := rtr.NewRoute().
		SetName("Admin > Users > User Delete").
		SetPath(links.ADMIN_USERS_USER_DELETE).
		SetHTMLHandler(NewUserDeleteController().Handler)

	userImpersonate := rtr.NewRoute().
		SetName("Admin > Users > User Impersonate").
		SetPath(links.ADMIN_USERS_USER_IMPERSONATE).
		SetHTMLHandler(NewUserImpersonateController().Handler)

	userManager := rtr.NewRoute().
		SetName("Admin > Users > User Manager").
		SetPath(links.ADMIN_USERS_USER_MANAGER).
		SetHTMLHandler(NewUserManagerController().Handler)

	userUpdate := rtr.NewRoute().
		SetName("Admin > Users > User Update").
		SetPath(links.ADMIN_USERS_USER_UPDATE).
		SetHTMLHandler(NewUserUpdateController().Handler)

	usersHome := rtr.NewRoute().
		SetName("Admin > Users > Home").
		SetPath(links.ADMIN_USERS).
		SetHTMLHandler(NewUserManagerController().Handler)

	usersCatchAll := rtr.NewRoute().
		SetName("Admin > Users > Catchall").
		SetPath(links.ADMIN_USERS + links.CATCHALL).
		SetHTMLHandler(NewUserManagerController().Handler)

	return []rtr.RouteInterface{
		userCreate,
		userDelete,
		userImpersonate,
		userManager,
		userUpdate,
		usersHome,
		usersCatchAll,
	}
}
