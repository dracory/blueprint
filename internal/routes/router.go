package routes

import (
	"project/internal/controllers/admin"
	"project/internal/controllers/auth"
	"project/internal/controllers/liveflux"
	"project/internal/controllers/shared"
	"project/internal/controllers/user"
	"project/internal/controllers/website"
	"project/internal/registry"
	"project/internal/widgets"

	"github.com/dracory/rtr"
)

func routes(registry registry.RegistryInterface) []rtr.RouteInterface {
	routes := []rtr.RouteInterface{}

	routes = append(routes, admin.Routes(registry)...)
	routes = append(routes, auth.Routes(registry)...)
	routes = append(routes, liveflux.Routes(registry)...)
	routes = append(routes, shared.Routes(registry)...)
	routes = append(routes, user.Routes(registry)...)
	routes = append(routes, widgets.Routes(registry)...)
	routes = append(routes, website.Routes(registry)...)

	return routes
}

func RoutesList(registry registry.RegistryInterface) (globalMiddlewareList []rtr.MiddlewareInterface, routeList []rtr.RouteInterface) {
	return globalMiddlewares(registry), routes(registry)
}

// Router creates the router for the registry.
func Router(registry registry.RegistryInterface) rtr.RouterInterface {
	r := rtr.NewRouter()

	// Add global middlewares
	globalMiddlewareList, routeList := RoutesList(registry)

	r.AddBeforeMiddlewares(globalMiddlewareList)

	// Add all routes
	for _, route := range routeList {
		r.AddRoute(route)
	}

	return r
}
