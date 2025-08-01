package user

import (
	userAccount "project/app/controllers/user/account"
	userHome "project/app/controllers/user/home"

	"project/app/middlewares"
	"project/internal/links"

	"github.com/dracory/rtr"
)

func Routes() []rtr.RouteInterface {
	home := rtr.NewRoute().
		SetName("User > Home").
		SetPath(links.USER_HOME).
		SetHTMLHandler(userHome.NewHomeController().Handler)

	homeCatchAll := rtr.NewRoute().
		SetName("User > Catch All Controller > Index Page").
		SetPath(links.USER_HOME + links.CATCHALL).
		SetHTMLHandler(userHome.NewHomeController().Handler)

	profile := rtr.NewRoute().
		SetName("User > Profile").
		SetPath(links.USER_PROFILE).
		SetHTMLHandler(userAccount.NewProfileController().Handler)

	userRoutes := []rtr.RouteInterface{
		profile,
		home,
		homeCatchAll,
	}

	for _, route := range userRoutes {
		route.AddBeforeMiddlewares([]rtr.MiddlewareInterface{middlewares.NewUserMiddleware()})
	}

	return userRoutes
}
