package admin

import (
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"

	userCreate "project/internal/controllers/admin/users/user_create"
	userDelete "project/internal/controllers/admin/users/user_delete"
	userImpersonate "project/internal/controllers/admin/users/user_impersonate"
	userManager "project/internal/controllers/admin/users/user_manager"
	userUpdate "project/internal/controllers/admin/users/user_update"
)

func UserRoutes(app types.AppInterface) []rtr.RouteInterface {
	userCreate := rtr.NewRoute().
		SetName("Admin > Users > User Create").
		SetPath(links.ADMIN_USERS_USER_CREATE).
		SetHTMLHandler(userCreate.NewUserCreateController(app).Handler)

	userDelete := rtr.NewRoute().
		SetName("Admin > Users > User Delete").
		SetPath(links.ADMIN_USERS_USER_DELETE).
		SetHTMLHandler(userDelete.NewUserDeleteController(app).Handler)

	userImpersonate := rtr.NewRoute().
		SetName("Admin > Users > User Impersonate").
		SetPath(links.ADMIN_USERS_USER_IMPERSONATE).
		SetHTMLHandler(userImpersonate.NewUserImpersonateController(app).Handler)

	userManagerRoute := rtr.NewRoute().
		SetName("Admin > Users > User Manager").
		SetPath(links.ADMIN_USERS_USER_MANAGER).
		SetHTMLHandler(userManager.NewUserManagerController(app).Handler)

	userUpdate := rtr.NewRoute().
		SetName("Admin > Users > User Update").
		SetPath(links.ADMIN_USERS_USER_UPDATE).
		SetHTMLHandler(userUpdate.NewUserUpdateController(app).Handler)

	usersHome := rtr.NewRoute().
		SetName("Admin > Users > Home").
		SetPath(links.ADMIN_USERS).
		SetHTMLHandler(userManager.NewUserManagerController(app).Handler)

	usersCatchAll := rtr.NewRoute().
		SetName("Admin > Users > Catchall").
		SetPath(links.ADMIN_USERS + links.CATCHALL).
		SetHTMLHandler(userManager.NewUserManagerController(app).Handler)

	return []rtr.RouteInterface{
		userCreate,
		userDelete,
		userImpersonate,
		userManagerRoute,
		userUpdate,
		usersHome,
		usersCatchAll,
	}
}
