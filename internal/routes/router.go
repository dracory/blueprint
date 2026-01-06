package routes

import (
	"project/internal/controllers/admin"
	"project/internal/controllers/auth"
	"project/internal/controllers/liveflux"
	"project/internal/controllers/shared"
	"project/internal/controllers/user"
	"project/internal/controllers/website"
	"project/internal/types"
	"project/internal/widgets"

	"github.com/dracory/rtr"
)

func routes(app types.RegistryInterface) []rtr.RouteInterface {
	routes := []rtr.RouteInterface{}

	routes = append(routes, admin.Routes(app)...)
	routes = append(routes, auth.Routes(app)...)
	routes = append(routes, liveflux.Routes(app)...)
	routes = append(routes, shared.Routes(app)...)
	routes = append(routes, user.Routes(app)...)
	routes = append(routes, widgets.Routes(app)...)
	routes = append(routes, website.Routes(app)...)

	return routes
}

func RoutesList(app types.RegistryInterface) (globalMiddlewareList []rtr.MiddlewareInterface, routeList []rtr.RouteInterface) {
	return globalMiddlewares(app), routes(app)
}

// Router creates the router for the registry.
func Router(app types.RegistryInterface) rtr.RouterInterface {
	r := rtr.NewRouter()

	// Add global middlewares
	globalMiddlewareList, routeList := RoutesList(app)

	r.AddBeforeMiddlewares(globalMiddlewareList)

	// Add all routes
	for _, route := range routeList {
		r.AddRoute(route)
	}

	return r
}
