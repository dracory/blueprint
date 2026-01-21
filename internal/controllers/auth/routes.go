package auth

import (
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/rtr"
	rtrMiddleware "github.com/dracory/rtr/middlewares"
)

func Routes(application registry.RegistryInterface) []rtr.RouteInterface {
	authRoutes := []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Auth > Auth Controller").
			SetPath(links.AUTH_AUTH).
			SetHTMLHandler(NewAuthenticationController(application).Handler),
		rtr.NewRoute().
			SetName("Auth > Login Controller").
			SetPath(links.AUTH_LOGIN).
			SetHTMLHandler(NewLoginController(application).Handler),
	}

	// Apply stricter rate limiting to sensitive authentication routes only
	for i := range authRoutes {
		authRoutes[i].AddBeforeMiddlewares([]rtr.MiddlewareInterface{
			// Stricter rate limiting for authentication endpoints
			// 5 requests per minute to prevent brute force attacks
			rtrMiddleware.RateLimitByIPMiddleware(5, 60),
		})
	}

	// Logout doesn't need rate limiting - it's not security-sensitive
	logoutRoute := rtr.NewRoute().
		SetName("Auth > Logout Controller").
		SetPath(links.AUTH_LOGOUT).
		SetHTMLHandler(NewLogoutController(application).AnyIndex)

	routes := append(authRoutes, logoutRoute)

	if application.GetConfig().GetRegistrationEnabled() {
		registerRoute := rtr.NewRoute().
			SetName("Auth > Register Controller").
			SetPath(links.AUTH_REGISTER).
			SetHTMLHandler(NewRegisterController(application).Handler)

		// Apply even stricter rate limiting for registration
		registerRoute.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
			// Stricter rate limiting for registration endpoint
			// 3 requests per minute to prevent spam registration
			rtrMiddleware.RateLimitByIPMiddleware(3, 60),
		})

		routes = append(routes, registerRoute)
	}

	return routes
}
