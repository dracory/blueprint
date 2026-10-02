package auth

import (
	"slices"

	"project/internal/app"
	"project/internal/config"
	"project/internal/controllers/auth/authentication_authknight"
	"project/internal/controllers/auth/forgot_password"
	"project/internal/controllers/auth/login_authknight"
	"project/internal/controllers/auth/login_magiclink"
	"project/internal/controllers/auth/login_otp"
	"project/internal/controllers/auth/login_password"
	"project/internal/controllers/auth/logout"
	"project/internal/controllers/auth/password_reset"
	"project/internal/controllers/auth/register"
	"project/internal/controllers/auth/register_password"
	"project/internal/links"

	"github.com/dracory/rtr"
	rtrMiddleware "github.com/dracory/rtr/middlewares"
)

func Routes(application app.AppInterface) []rtr.RouteInterface {
	loginMethods := application.GetConfig().GetLoginMethods()

	authRoutes := loginRoutes(application, loginMethods)

	// Verification callbacks mount once per callback-bearing method,
	// independent of whether the method is primary or secondary.
	if slices.Contains(loginMethods, config.LOGIN_METHOD_MAGICLINK) {
		authRoutes = append(authRoutes,
			rtr.GetHTML(links.AUTH_CALLBACK_MAGICLINK, login_magiclink.NewLoginController(application).Handler).
				SetName("Auth > MagicLink Controller"))
	}
	if slices.Contains(loginMethods, config.LOGIN_METHOD_AUTHKNIGHT) {
		authRoutes = append(authRoutes,
			rtr.GetHTML(links.AUTH_CALLBACK_AUTHKNIGHT, authentication_authknight.NewAuthenticationController(application).Handler).
				SetName("Auth > AuthKnight Controller"))
	}

	// The forgot/reset password pages exist whenever password auth is
	// enabled — the passwordless methods have no password to recover,
	// so no dead routes are mounted.
	if slices.Contains(loginMethods, config.LOGIN_METHOD_PASSWORD) {
		authRoutes = append(authRoutes, passwordRecoveryRoutes(application)...)
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
	routes := append(authRoutes, rtr.NewRoute().
		SetName("Auth > Logout Controller").
		SetPath(links.AUTH_LOGOUT).
		SetHTMLHandler(logout.NewLogoutController(application).AnyIndex))

	if application.GetConfig().GetRegistrationEnabled() {
		routes = append(routes, registerRoutes(application, loginMethods)...)
	}

	return routes
}

// loginRoutes mounts every enabled login method. The primary method
// (methods[0]) is mounted at links.AUTH_LOGIN; secondary methods get
// their own paths (/auth/<method>-login).
func loginRoutes(application app.AppInterface, loginMethods []string) []rtr.RouteInterface {
	routes := mountMethod(application, loginMethods[0], true)
	for _, method := range loginMethods[1:] {
		routes = append(routes, mountMethod(application, method, false)...)
	}
	return routes
}

// mountMethod mounts the GET page + POST ajax routes for one login
// method. primary=true mounts at links.AUTH_LOGIN with the canonical
// "Auth > Login Controller" route name; secondary methods mount at their
// dedicated paths with a " (alt)" name suffix so names stay unique.
func mountMethod(application app.AppInterface, method string, primary bool) []rtr.RouteInterface {
	suffix := ""
	if !primary {
		suffix = " (alt)"
	}

	switch method {
	case config.LOGIN_METHOD_MAGICLINK:
		basePath := links.AUTH_LOGIN
		if !primary {
			basePath = links.AUTH_LOGIN_MAGICLINK
		}
		controller := login_magiclink.NewLoginController(application)
		controller.SetBasePath(basePath)
		return []rtr.RouteInterface{
			rtr.GetHTML(basePath, controller.PageHandler).
				SetName("Auth > Login Controller" + suffix),
			rtr.PostJSON(basePath, controller.AjaxHandler).
				SetName("Auth > Login MagicLink Ajax Controller" + suffix),
		}
	case config.LOGIN_METHOD_OTP:
		basePath := links.AUTH_LOGIN
		if !primary {
			basePath = links.AUTH_LOGIN_OTP
		}
		controller := login_otp.NewLoginController(application)
		controller.SetBasePath(basePath)
		return []rtr.RouteInterface{
			rtr.GetHTML(basePath, controller.PageHandler).
				SetName("Auth > Login Controller" + suffix),
			rtr.PostJSON(basePath, controller.AjaxHandler).
				SetName("Auth > Login OTP Ajax Controller" + suffix),
		}
	case config.LOGIN_METHOD_PASSWORD:
		basePath := links.AUTH_LOGIN
		if !primary {
			basePath = links.AUTH_LOGIN_PASSWORD
		}
		controller := login_password.NewLoginController(application)
		controller.SetBasePath(basePath)
		return []rtr.RouteInterface{
			rtr.GetHTML(basePath, controller.PageHandler).
				SetName("Auth > Login Controller" + suffix),
			rtr.PostJSON(basePath, controller.AjaxHandler).
				SetName("Auth > Login Password Ajax Controller" + suffix),
		}
	case config.LOGIN_METHOD_AUTHKNIGHT:
		basePath := links.AUTH_LOGIN
		if !primary {
			basePath = links.AUTH_LOGIN_AUTHKNIGHT
		}
		return []rtr.RouteInterface{
			rtr.GetHTML(basePath, login_authknight.NewLoginController(application).Handler).
				SetName("Auth > Login Controller" + suffix),
		}
	default:
		panic("invalid login method: " + method)
	}
}

// passwordRecoveryRoutes mounts the forgot-password and password-reset pages.
func passwordRecoveryRoutes(application app.AppInterface) []rtr.RouteInterface {
	forgotPasswordController := forgot_password.NewForgotPasswordController(application)
	passwordResetController := password_reset.NewPasswordResetController(application)

	return []rtr.RouteInterface{
		rtr.GetHTML(links.AUTH_FORGOT_PASSWORD, forgotPasswordController.PageHandler).
			SetName("Auth > Forgot Password Controller"),
		rtr.PostJSON(links.AUTH_FORGOT_PASSWORD, forgotPasswordController.AjaxHandler).
			SetName("Auth > Forgot Password Ajax Controller"),
		rtr.GetHTML(links.AUTH_PASSWORD_RESET, passwordResetController.PageHandler).
			SetName("Auth > Password Reset Controller"),
		rtr.PostJSON(links.AUTH_PASSWORD_RESET, passwordResetController.AjaxHandler).
			SetName("Auth > Password Reset Ajax Controller"),
	}
}

// registerRoutes mounts the register page. When password auth is among
// the enabled methods it is the public sign-up form that creates the
// account (email + password) and logs the user in; otherwise it is the
// post-authentication profile completion form.
func registerRoutes(application app.AppInterface, loginMethods []string) []rtr.RouteInterface {
	var registerRoute, registerAjaxRoute rtr.RouteInterface

	if slices.Contains(loginMethods, config.LOGIN_METHOD_PASSWORD) {
		registerPasswordController := register_password.NewRegisterController(application)
		registerRoute = rtr.GetHTML(links.AUTH_REGISTER, registerPasswordController.PageHandler).
			SetName("Auth > Register Controller")
		registerAjaxRoute = rtr.PostJSON(links.AUTH_REGISTER, registerPasswordController.AjaxHandler).
			SetName("Auth > Register Ajax Controller")
	} else {
		registerController := register.NewRegisterController(application)
		registerRoute = rtr.GetHTML(links.AUTH_REGISTER, registerController.PageHandler).
			SetName("Auth > Register Controller")
		registerAjaxRoute = rtr.PostJSON(links.AUTH_REGISTER, registerController.AjaxHandler).
			SetName("Auth > Register Ajax Controller")
	}

	// Apply moderate rate limiting for registration
	// 10 requests per minute to allow for legitimate interactions (country/timezone selection)
	middlewares := []rtr.MiddlewareInterface{rtrMiddleware.RateLimitByIPMiddleware(10, 60)}
	registerRoute.AddBeforeMiddlewares(middlewares)
	registerAjaxRoute.AddBeforeMiddlewares(middlewares)

	return []rtr.RouteInterface{registerRoute, registerAjaxRoute}
}
