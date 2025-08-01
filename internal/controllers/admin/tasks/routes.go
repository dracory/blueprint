package admin

import (
	"project/internal/links"

	"github.com/dracory/rtr"
)

func TaskRoutes() []rtr.RouteInterface {
	return []rtr.RouteInterface{
		// &router.Route{
		// 	Name:        "Admin > Users > User Create",
		// 	Path:        links.ADMIN_USERS_USER_CREATE,
		// 	HTMLHandler: NewUserCreateController().Handler,
		// },
		// &router.Route{
		// 	Name:        "Admin > Users > User Delete",
		// 	Path:        links.ADMIN_USERS_USER_DELETE,
		// 	HTMLHandler: NewUserDeleteController().Handler,
		// },
		// &router.Route{
		// 	Name:        "Admin > Users > User Impersonate",
		// 	Path:        links.ADMIN_USERS_USER_IMPERSONATE,
		// 	HTMLHandler: NewUserImpersonateController().Handler,
		// },
		// &router.Route{
		// 	Name:        "Admin > Users > User Manager",
		// 	Path:        links.ADMIN_USERS_USER_MANAGER,
		// 	HTMLHandler: NewUserManagerController().Handler,
		// },
		// &router.Route{
		// 	Name:        "Admin > Users > User Update",
		// 	Path:        links.ADMIN_USERS_USER_UPDATE,
		// 	HTMLHandler: NewUserUpdateController().Handler,
		// },
		rtr.NewRoute().
			SetName("Admin > Tasks > Home").
			SetPath(links.ADMIN_TASKS).
			SetHTMLHandler(TaskController().Handler),
		rtr.NewRoute().
			SetName("Admin > Tasks > Catchall").
			SetPath(links.ADMIN_TASKS + links.CATCHALL).
			SetHTMLHandler(TaskController().Handler),
	}
}
