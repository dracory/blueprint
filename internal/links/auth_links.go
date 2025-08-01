package links

import "github.com/samber/lo"

type authLinks struct {
}

// Auth is a shortcut for NewAuthLinks
func Auth() *authLinks {
	return NewAuthLinks()
}

// Deprecated: Use Auth() instead. NewAuthLinks will be removed in the next major version.
func NewAuthLinks() *authLinks {
	return &authLinks{}
}

func (l *authLinks) Auth(params ...map[string]string) string {
	p := lo.IfF(len(params) > 0, func() map[string]string { return params[0] }).Else(map[string]string{})
	return URL(AUTH_AUTH, p)
}

func (l *authLinks) AuthKnightLogin(backUrl string) string {
	params := map[string]string{
		"back_url": backUrl,
		"next_url": l.Auth(),
	}
	return "https://authknight.com/app/login" + query(params)
}

func (l *authLinks) Login(backUrl string, params ...map[string]string) string {
	p := lo.IfF(len(params) > 0, func() map[string]string { return params[0] }).Else(map[string]string{})

	if backUrl != "" {
		p["back_url"] = backUrl
	}

	return URL(AUTH_LOGIN, p)
}

func (l *authLinks) Logout(params ...map[string]string) string {
	p := lo.IfF(len(params) > 0, func() map[string]string { return params[0] }).Else(map[string]string{})
	return URL(AUTH_LOGOUT, p)
}

func (l *authLinks) Register(params ...map[string]string) string {
	p := lo.IfF(len(params) > 0, func() map[string]string { return params[0] }).Else(map[string]string{})
	return URL(AUTH_REGISTER, p)
}
