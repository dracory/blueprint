package auth

import (
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
	loginMethod := application.GetConfig().GetLoginMethod()

	authRoutes := loginRoutes(application, loginMethod)

	// The forgot/reset password pages only exist for the password login
	// method — the passwordless methods have no password to recover, so no
	// dead routes are mounted.
	if loginMethod == config.LOGIN_METHOD_PASSWORD {
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
		routes = append(routes, registerRoutes(application, loginMethod)...)
	}

	return routes
}

// loginRoutes mounts the login controller selected by AUTH_LOGIN_METHOD.
// In OTP mode the AuthKnight callback route is omitted entirely so no
// dead external dependency remains.
func loginRoutes(application app.AppInterface, loginMethod string) []rtr.RouteInterface {
	switch loginMethod {
	case config.LOGIN_METHOD_AUTHKNIGHT:
		return []rtr.RouteInterface{
			rtr.GetHTML(links.AUTH_LOGIN, login_authknight.NewLoginController(application).Handler).
				SetName("Auth > Login Controller"),
			rtr.NewRoute().
				SetName("Auth > Auth Controller").
				SetPath(links.AUTH_AUTH).
				SetHTMLHandler(authentication_authknight.NewAuthenticationController(application).Handler),
		}
	case config.LOGIN_METHOD_MAGICLINK:
		magicLinkController := login_magiclink.NewLoginController(application)
		return []rtr.RouteInterface{
			rtr.GetHTML(links.AUTH_LOGIN, magicLinkController.PageHandler).
				SetName("Auth > Login Controller"),
			rtr.PostJSON(links.AUTH_LOGIN, magicLinkController.AjaxHandler).
				SetName("Auth > Login MagicLink Ajax Controller"),
			rtr.GetHTML(links.AUTH_AUTH, magicLinkController.Handler).
				SetName("Auth > MagicLink Controller"),
		}
	case config.LOGIN_METHOD_OTP:
		otpController := login_otp.NewLoginController(application)
		return []rtr.RouteInterface{
			rtr.GetHTML(links.AUTH_LOGIN, otpController.PageHandler).
				SetName("Auth > Login Controller"),
			rtr.PostJSON(links.AUTH_LOGIN, otpController.AjaxHandler).
				SetName("Auth > Login OTP Ajax Controller"),
		}
	case config.LOGIN_METHOD_PASSWORD:
		passwordController := login_password.NewLoginController(application)
		return []rtr.RouteInterface{
			rtr.GetHTML(links.AUTH_LOGIN, passwordController.PageHandler).
				SetName("Auth > Login Controller"),
			rtr.PostJSON(links.AUTH_LOGIN, passwordController.AjaxHandler).
				SetName("Auth > Login Password Ajax Controller"),
		}
	default:
		panic("invalid login method: " + loginMethod)
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

// registerRoutes mounts the register page. In password mode it is the public
// sign-up form that creates the account (email + password) and logs the user
// in; in the passwordless/external modes it is the post-authentication
// profile completion form.
func registerRoutes(application app.AppInterface, loginMethod string) []rtr.RouteInterface {
	var registerRoute, registerAjaxRoute rtr.RouteInterface

	if loginMethod == config.LOGIN_METHOD_PASSWORD {
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
