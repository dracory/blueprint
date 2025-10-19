package user

import (
	userAccount "project/internal/controllers/user/account"
	userHome "project/internal/controllers/user/home"
	"project/internal/types"

	"project/internal/links"
	"project/internal/middlewares"

	"github.com/dracory/rtr"
)

func Routes(app types.AppInterface) []rtr.RouteInterface {
	home := rtr.NewRoute().
		SetName("User > Home").
		SetPath(links.USER_HOME).
		SetHTMLHandler(userHome.NewHomeController(app).Handler)

	homeCatchAll := rtr.NewRoute().
		SetName("User > Catch All Controller > Index Page").
		SetPath(links.USER_HOME + links.CATCHALL).
		SetHTMLHandler(userHome.NewHomeController(app).Handler)

	profile := rtr.NewRoute().
		SetName("User > Profile").
		SetPath(links.USER_PROFILE).
		SetHTMLHandler(userAccount.NewProfileController(app).Handler)

	userRoutes := []rtr.RouteInterface{}
	userRoutes = append(userRoutes, profile)
	userRoutes = append(userRoutes, home)
	userRoutes = append(userRoutes, homeCatchAll) // Must be last!

	applyUserMiddleware(app, userRoutes)

	return userRoutes
}

func applyUserMiddleware(app types.AppInterface, routes []rtr.RouteInterface) {
	for _, route := range routes {
		middlewaresToAdd := []rtr.MiddlewareInterface{
			middlewares.NewUserMiddleware(app),
		}

		route.AddBeforeMiddlewares(middlewaresToAdd)
	}
}
