package auth

import (
	"project/internal/app"
	"project/internal/config"
	"project/internal/controllers/auth/authentication_authknight"
	"project/internal/controllers/auth/login_authknight"
	"project/internal/controllers/auth/login_otp"
	"project/internal/controllers/auth/logout"
	"project/internal/controllers/auth/register"
	"project/internal/links"

	"github.com/dracory/rtr"
	rtrMiddleware "github.com/dracory/rtr/middlewares"
)

func Routes(application app.AppInterface) []rtr.RouteInterface {
	loginRoute := rtr.NewRoute().
		SetName("Auth > Login Controller").
		SetMethod("GET").
		SetPath(links.AUTH_LOGIN)

	authRoutes := []rtr.RouteInterface{}

	// Mount the login controller selected by config.LOGIN_METHOD.
	// In OTP mode the AuthKnight callback route is omitted entirely so no
	// dead external dependency remains.
	switch config.LOGIN_METHOD {
	case config.LOGIN_METHOD_AUTHKNIGHT:
		loginRoute.SetHTMLHandler(login_authknight.NewLoginController(application).Handler)

		authRoutes = append(authRoutes, rtr.NewRoute().
			SetName("Auth > Auth Controller").
			SetPath(links.AUTH_AUTH).
			SetHTMLHandler(authentication_authknight.NewAuthenticationController(application).Handler))
	case config.LOGIN_METHOD_OTP:
		otpController := login_otp.NewLoginController(application)
		loginRoute.SetHTMLHandler(otpController.PageHandler)
		authRoutes = append(authRoutes, rtr.PostJSON(links.AUTH_LOGIN, otpController.AjaxHandler).
			SetName("Auth > Login OTP Ajax Controller"))
	default:
		panic("invalid config.LOGIN_METHOD: " + config.LOGIN_METHOD)
	}

	authRoutes = append(authRoutes, loginRoute)

	logoutRoute := rtr.NewRoute().
		SetName("Auth > Logout Controller").
		SetPath(links.AUTH_LOGOUT).
		SetHTMLHandler(logout.NewLogoutController(application).AnyIndex)

	registerRoute := rtr.NewRoute().
		SetName("Auth > Register Controller").
		SetPath(links.AUTH_REGISTER).
		SetHTMLHandler(register.NewRegisterController(application).Handler)

	// Apply stricter rate limiting to sensitive authentication routes only
	for i := range authRoutes {
		authRoutes[i].AddBeforeMiddlewares([]rtr.MiddlewareInterface{
			// Stricter rate limiting for authentication endpoints
			// 5 requests per minute to prevent brute force attacks
			rtrMiddleware.RateLimitByIPMiddleware(5, 60),
		})
	}

	// Logout doesn't need rate limiting - it's not security-sensitive
	routes := append(authRoutes, logoutRoute)

	if application.GetConfig().GetRegistrationEnabled() {
		// Apply moderate rate limiting for registration
		registerRoute.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
			// Moderate rate limiting for registration endpoint
			// 10 requests per minute to allow for legitimate interactions (country/timezone selection)
			rtrMiddleware.RateLimitByIPMiddleware(10, 60),
		})

		routes = append(routes, registerRoute)
	}

	return routes
}
