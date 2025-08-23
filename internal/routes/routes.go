package routes

import (
	"project/internal/controllers/admin"
	"project/internal/controllers/auth"
	"project/internal/controllers/shared"
	"project/internal/controllers/user"
	"project/internal/controllers/website"
	"project/internal/types"
	"project/internal/widgets"

	"github.com/dracory/rtr"
)

func routes(app types.AppInterface) []rtr.RouteInterface {
	routes := []rtr.RouteInterface{}

	routes = append(routes, admin.Routes(app)...)
	routes = append(routes, auth.Routes(app)...)
	routes = append(routes, shared.Routes(app)...)
	routes = append(routes, user.Routes(app)...)
	routes = append(routes, widgets.Routes()...)
	routes = append(routes, website.Routes(app)...)

	return routes
}

func RoutesList(app types.AppInterface) (globalMiddlewareList []rtr.MiddlewareInterface, routeList []rtr.RouteInterface) {
	return globalMiddlewares(app), routes(app)
}

// Routes returns the routes of the application
func Routes(app types.AppInterface) rtr.RouterInterface {
	r := rtr.NewRouter()

	// Add global middlewares
	globalMiddlewareList, routes := RoutesList(app)

	r.AddBeforeMiddlewares(globalMiddlewareList)

	// Add all routes
	for _, route := range routes {
		r.AddRoute(route)
	}

	return r
}
