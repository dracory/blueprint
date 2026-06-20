package routes

import (
	"project/internal/app"
	"project/internal/controllers/admin"
	"project/internal/controllers/auth"
	"project/internal/controllers/liveflux"
	"project/internal/controllers/shared"
	"project/internal/controllers/user"
	"project/internal/controllers/website"
	"project/internal/middlewares"
	"project/internal/widgets"

	"github.com/dracory/rtr"
)

func routes(app app.AppInterface) []rtr.RouteInterface {
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

func RoutesList(app app.AppInterface) (globalMiddlewareList []rtr.MiddlewareInterface, routeList []rtr.RouteInterface) {
	return globalMiddlewares(app), routes(app)
}

// Router creates the router for the app.
func Router(app app.AppInterface) rtr.RouterInterface {
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

// AiBrowserRouter wraps Router with the AI browser auto-login middleware prepended.
// It must only be used by cmd/ai-browser — never in production.
func AiBrowserRouter(app app.AppInterface) rtr.RouterInterface {
	r := rtr.NewRouter()

	globalMiddlewareList, routeList := RoutesList(app)

	// Prepend auto-login so every unauthenticated request gets a session before
	// the standard AuthMiddleware runs.
	globalMiddlewareList = append([]rtr.MiddlewareInterface{
		middlewares.AiBrowserAutoLoginMiddleware(app),
	}, globalMiddlewareList...)

	r.AddBeforeMiddlewares(globalMiddlewareList)

	for _, route := range routeList {
		r.AddRoute(route)
	}

	return r
}
