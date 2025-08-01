package routes

import (
	"project/app/controllers/admin"
	"project/app/controllers/auth"
	"project/app/controllers/shared"
	"project/app/controllers/user"
	"project/app/controllers/website"
	"project/internal/widgets"

	"github.com/dracory/rtr"
)

func routes() []rtr.RouteInterface {
	routes := []rtr.RouteInterface{}

	routes = append(routes, admin.Routes()...)
	routes = append(routes, auth.Routes()...)
	routes = append(routes, shared.Routes()...)
	routes = append(routes, user.Routes()...)
	routes = append(routes, widgets.Routes()...)
	routes = append(routes, website.Routes()...)

	return routes
}

func RoutesList() (globalMiddlewareList []rtr.MiddlewareInterface, routeList []rtr.RouteInterface) {
	return globalMiddlewares(), routes()
}

// Routes returns the routes of the application
func Routes() rtr.RouterInterface {
	r := rtr.NewRouter()

	// Add global middlewares
	globalMiddlewares, routes := RoutesList()
	r.AddBeforeMiddlewares(globalMiddlewares)

	// Add all routes
	for _, route := range routes {
		r.AddRoute(route)
	}

	return r
}
