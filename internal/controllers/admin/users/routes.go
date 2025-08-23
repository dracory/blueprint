package admin

import (
    "project/internal/links"
    "project/internal/types"

    "github.com/dracory/rtr"
)

func UserRoutes(app types.AppInterface) []rtr.RouteInterface {
    userCreate := rtr.NewRoute().
        SetName("Admin > Users > User Create").
        SetPath(links.ADMIN_USERS_USER_CREATE).
        SetHTMLHandler(NewUserCreateController(app).Handler)

    userDelete := rtr.NewRoute().
        SetName("Admin > Users > User Delete").
        SetPath(links.ADMIN_USERS_USER_DELETE).
        SetHTMLHandler(NewUserDeleteController(app).Handler)

    userImpersonate := rtr.NewRoute().
        SetName("Admin > Users > User Impersonate").
        SetPath(links.ADMIN_USERS_USER_IMPERSONATE).
        SetHTMLHandler(NewUserImpersonateController(app).Handler)

    userManager := rtr.NewRoute().
        SetName("Admin > Users > User Manager").
        SetPath(links.ADMIN_USERS_USER_MANAGER).
        SetHTMLHandler(NewUserManagerController(app).Handler)

    userUpdate := rtr.NewRoute().
        SetName("Admin > Users > User Update").
        SetPath(links.ADMIN_USERS_USER_UPDATE).
        SetHTMLHandler(NewUserUpdateController(app).Handler)

    usersHome := rtr.NewRoute().
        SetName("Admin > Users > Home").
        SetPath(links.ADMIN_USERS).
        SetHTMLHandler(NewUserManagerController(app).Handler)

    usersCatchAll := rtr.NewRoute().
        SetName("Admin > Users > Catchall").
        SetPath(links.ADMIN_USERS + links.CATCHALL).
        SetHTMLHandler(NewUserManagerController(app).Handler)

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
