package auth

import (
	"slices"
	"strings"

	"project/internal/app"
	"project/internal/authadapters"
	"project/internal/config"
	"project/internal/controllers/auth/logout"
	"project/internal/controllers/auth/register"
	"project/internal/controllers/auth/register_password"
	"project/internal/links"

	authpkg "github.com/dracory/auth"
	"github.com/dracory/auth/types"
	"github.com/dracory/rtr"
	rtrMiddleware "github.com/dracory/rtr/middlewares"
)

// Routes builds the /auth surface. The standalone dracory/auth package
// serves every login method plus its built-in pages through one mux
// (mounted as a catch-all). App-owned routes mount first so they win on
// their exact paths: register (profile completion + project rules) and
// logout (remember-session cleanup).
func Routes(application app.AppInterface) []rtr.RouteInterface {
	loginMethods := application.GetConfig().GetLoginMethods()

	routes := []rtr.RouteInterface{}

	if application.GetConfig().GetRegistrationEnabled() {
		routes = append(routes, registerRoutes(application, loginMethods)...)
	}

	// Logout stays a local controller — it clears admin's remember
	// session/cookie in addition to the auth session, and is not
	// rate-limited (not security-sensitive).
	routes = append(routes, rtr.NewRoute().
		SetName("Auth > Logout Controller").
		SetPath(links.AUTH_LOGOUT).
		SetHTMLHandler(logout.NewLogoutController(application).AnyIndex))

	return append(routes, packageAuthRoutes(application)...)
}

// packageAuthRoutes composes the enabled login methods into the package's
// multi-method mux and mounts it on /auth/*. Admin's logout and register
// routes still mount ahead of it (they add behaviour the package lacks:
// remember-cookie cleanup, project rules) as thin wrappers.
func packageAuthRoutes(application app.AppInterface) []rtr.RouteInterface {
	loginMethods := application.GetConfig().GetLoginMethods()

	var methods []types.Method
	for i, name := range loginMethods {
		alternatives := loginAlternatives(loginMethods, i)
		switch name {
		case config.LOGIN_METHOD_OTP, config.LOGIN_METHOD_MAGICLINK:
			a, err := authadapters.BuildPasswordlessAuth(application)
			if err != nil {
				panic("auth package: " + err.Error())
			}
			applyBranding(application, a, alternatives)
			if name == config.LOGIN_METHOD_OTP {
				methods = append(methods, authpkg.NewOTPMethodWithPage(a, nil))
			} else {
				methods = append(methods, authpkg.NewMagicLinkMethodFromAuth(a, nil,
					func(token string) string {
						return links.Auth().Auth(map[string]string{"token": token})
					}))
			}
		case config.LOGIN_METHOD_PASSWORD:
			pa, err := authadapters.BuildPasswordAuth(application)
			if err != nil {
				panic("auth package: " + err.Error())
			}
			applyBranding(application, pa, alternatives)
			methods = append(methods, authpkg.NewPasswordMethod(pa, nil))
		case config.LOGIN_METHOD_AUTHKNIGHT:
			ka, err := authadapters.BuildAuthKnightAuth(application)
			if err != nil {
				panic("auth package: " + err.Error())
			}
			applyBranding(application, ka, alternatives)
			methods = append(methods, authpkg.NewAuthKnightMethod(ka))
		}
	}

	mma, err := authpkg.NewMultiMethodAuth("/auth", methods)
	if err != nil {
		panic("auth package mux: " + err.Error())
	}
	mux := mma.Router()

	stricter := []rtr.MiddlewareInterface{
		// 5 requests per minute against brute force on auth endpoints
		rtrMiddleware.RateLimitByIPMiddleware(5, 60),
	}
	base := rtr.NewRoute().SetPath("/auth").SetHandler(mux.ServeHTTP).SetName("Auth > Package")
	base.AddBeforeMiddlewares(stricter)
	catchAll := rtr.NewRoute().SetPath("/auth/*").SetHandler(mux.ServeHTTP).SetName("Auth > Package (catch-all)")
	catchAll.AddBeforeMiddlewares(stricter)

	return []rtr.RouteInterface{base, catchAll}
}

// methodLoginPath mirrors the package's mount convention: the primary
// method lives at /auth/login, secondary methods at /auth/login-<name>.
func methodLoginPath(loginMethods []string, method string) string {
	if len(loginMethods) > 0 && loginMethods[0] == method {
		return "/auth/login"
	}
	return "/auth/login-" + method
}

// loginAlternatives builds cross-links from the method at selfIndex to all
// other configured methods, so each built-in page can offer "use X instead".
func loginAlternatives(loginMethods []string, selfIndex int) []types.LoginAlternative {
	var alternatives []types.LoginAlternative
	for i, method := range loginMethods {
		if i == selfIndex {
			continue
		}
		alternatives = append(alternatives, types.LoginAlternative{
			Label: methodLabel(method),
			URL:   methodLoginPath(loginMethods, method),
			Icon:  methodIcon(method),
		})
	}
	return alternatives
}

func methodLabel(method string) string {
	switch method {
	case config.LOGIN_METHOD_MAGICLINK:
		return "Use magic link instead"
	case config.LOGIN_METHOD_OTP:
		return "Use login code instead"
	case config.LOGIN_METHOD_PASSWORD:
		return "Use password instead"
	case config.LOGIN_METHOD_AUTHKNIGHT:
		return "Sign in with AuthKnight"
	default:
		return "Use " + method + " instead"
	}
}

func methodIcon(method string) string {
	switch method {
	case config.LOGIN_METHOD_MAGICLINK:
		return "bi-envelope"
	case config.LOGIN_METHOD_OTP:
		return "bi-shield-lock"
	case config.LOGIN_METHOD_PASSWORD:
		return "bi-key"
	case config.LOGIN_METHOD_AUTHKNIGHT:
		return "bi-shield-check"
	default:
		return "bi-box-arrow-in-right"
	}
}

// applyBranding pushes the app brand + sibling-method links onto a built
// package auth instance via its optional setter interface.
func applyBranding(application app.AppInterface, a any, alternatives []types.LoginAlternative) {
	type brandingSetter interface {
		SetBrandName(string)
		SetBrandTagline(string)
		SetLoginAlternatives([]types.LoginAlternative)
	}
	s, ok := a.(brandingSetter)
	if !ok {
		return
	}
	s.SetBrandName(strings.ToUpper(application.GetConfig().GetAppName()))
	s.SetBrandTagline("")
	s.SetLoginAlternatives(alternatives)
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

	// 10 requests per minute to allow for legitimate interactions
	// (country/timezone selection)
	middlewares := []rtr.MiddlewareInterface{rtrMiddleware.RateLimitByIPMiddleware(10, 60)}
	registerRoute.AddBeforeMiddlewares(middlewares)
	registerAjaxRoute.AddBeforeMiddlewares(middlewares)

	return []rtr.RouteInterface{registerRoute, registerAjaxRoute}
}
