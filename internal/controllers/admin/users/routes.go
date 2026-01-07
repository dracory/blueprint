package users

import (
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/rtr"

	userCreate "project/internal/controllers/admin/users/user_create"
	userDelete "project/internal/controllers/admin/users/user_delete"
	userImpersonate "project/internal/controllers/admin/users/user_impersonate"
	userManager "project/internal/controllers/admin/users/user_manager"
	userUpdate "project/internal/controllers/admin/users/user_update"
)

func UserRoutes(registry registry.RegistryInterface) []rtr.RouteInterface {
	userCreate := rtr.NewRoute().
		SetName("Admin > Users > User Create").
		SetPath(links.ADMIN_USERS_USER_CREATE).
		SetHTMLHandler(userCreate.NewUserCreateController(registry).Handler)

	userDelete := rtr.NewRoute().
		SetName("Admin > Users > User Delete").
		SetPath(links.ADMIN_USERS_USER_DELETE).
		SetHTMLHandler(userDelete.NewUserDeleteController(registry).Handler)

	userImpersonate := rtr.NewRoute().
		SetName("Admin > Users > User Impersonate").
		SetPath(links.ADMIN_USERS_USER_IMPERSONATE).
		SetHTMLHandler(userImpersonate.NewUserImpersonateController(registry).Handler)

	userManagerRoute := rtr.NewRoute().
		SetName("Admin > Users > User Manager").
		SetPath(links.ADMIN_USERS_USER_MANAGER).
		SetHTMLHandler(userManager.NewUserManagerController(registry).Handler)

	userUpdate := rtr.NewRoute().
		SetName("Admin > Users > User Update").
		SetPath(links.ADMIN_USERS_USER_UPDATE).
		SetHTMLHandler(userUpdate.NewUserUpdateController(registry).Handler)

	usersHome := rtr.NewRoute().
		SetName("Admin > Users > Home").
		SetPath(links.ADMIN_USERS).
		SetHTMLHandler(userManager.NewUserManagerController(registry).Handler)

	usersCatchAll := rtr.NewRoute().
		SetName("Admin > Users > Catchall").
		SetPath(links.ADMIN_USERS + links.CATCHALL).
		SetHTMLHandler(userManager.NewUserManagerController(registry).Handler)

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
