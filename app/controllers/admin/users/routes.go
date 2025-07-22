package admin

import (
	"project/app/links"

	"github.com/dracory/rtr"
)

func UserRoutes() []rtr.RouteInterface {
	return []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Admin > Users > User Create").
			SetPath(links.ADMIN_USERS_USER_CREATE).
			SetHTMLHandler(NewUserCreateController().Handler),
		rtr.NewRoute().
			SetName("Admin > Users > User Delete").
			SetPath(links.ADMIN_USERS_USER_DELETE).
			SetHTMLHandler(NewUserDeleteController().Handler),
		rtr.NewRoute().
			SetName("Admin > Users > User Impersonate").
			SetPath(links.ADMIN_USERS_USER_IMPERSONATE).
			SetHTMLHandler(NewUserImpersonateController().Handler),
		rtr.NewRoute().
			SetName("Admin > Users > User Manager").
			SetPath(links.ADMIN_USERS_USER_MANAGER).
			SetHTMLHandler(NewUserManagerController().Handler),
		rtr.NewRoute().
			SetName("Admin > Users > User Update").
			SetPath(links.ADMIN_USERS_USER_UPDATE).
			SetHTMLHandler(NewUserUpdateController().Handler),
		rtr.NewRoute().
			SetName("Admin > Users > Home").
			SetPath(links.ADMIN_USERS).
			SetHTMLHandler(NewUserManagerController().Handler),
		rtr.NewRoute().
			SetName("Admin > Users > Catchall").
			SetPath(links.ADMIN_USERS + links.CATCHALL).
			SetHTMLHandler(NewUserManagerController().Handler),
		rtr.NewRoute().
			SetName("Admin > Users > Home").
			SetPath(links.ADMIN_USERS).
			SetHTMLHandler(NewUserManagerController().Handler),
		rtr.NewRoute().
			SetName("Admin > Users > Catchall").
			SetPath(links.ADMIN_USERS + links.CATCHALL).
			SetHTMLHandler(NewUserManagerController().Handler),
	}
}
