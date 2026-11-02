package auth

import (
	"strings"

	"project/internal/app"
	"project/internal/authadapters"
	"project/internal/config"
	"project/internal/links"

	authpkg "github.com/dracory/auth"
	"github.com/dracory/auth/types"
	"github.com/dracory/rtr"
	rtrMiddleware "github.com/dracory/rtr/middlewares"
)

// Routes builds the /auth surface. The dracory/auth package serves every
// login method plus its built-in pages (login, register, register-code
// verify, logout) through one mux, mounted as a catch-all on /auth/*.
// No app-owned overrides: an authenticated requester on /auth/register
// gets the package's profile-completion form (verified email read-only).
func Routes(application app.AppInterface) []rtr.RouteInterface {
	return packageAuthRoutes(application)
}

// packageAuthRoutes composes the enabled login methods into the package's
// multi-method mux and mounts it on /auth/*.
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
		// Coarse IP cap on the whole /auth surface. Kept loose (30/min)
		// because a single OTP journey legitimately needs ~5 requests
		// (login page → send → verify page → verify → redirect); the
		// package additionally rate-limits the sensitive API endpoints
		// and throttles email sends internally.
		rtrMiddleware.RateLimitByIPMiddleware(30, 60),
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
