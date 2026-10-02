package links

import "github.com/samber/lo"

type authLinks struct {
}

// Auth is a shortcut for NewAuthLinks
func Auth() *authLinks {
	return &authLinks{}
}

// Auth returns the magic-link verification callback URL.
func (l *authLinks) Auth(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(AUTH_CALLBACK_MAGICLINK, p)
}

func (l *authLinks) AuthKnightLogin(backUrl string) string {
	params := map[string]string{
		"back_url": backUrl,
		"next_url": URL(AUTH_CALLBACK_AUTHKNIGHT, nil),
	}
	return "https://authknight.com/app/login" + query(params)
}

func (l *authLinks) ForgotPassword(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(AUTH_FORGOT_PASSWORD, p)
}

func (l *authLinks) Login(backUrl string, params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})

	if backUrl != "" {
		p["back_url"] = backUrl
	}

	return URL(AUTH_LOGIN, p)
}

// LoginAuthknight returns the AuthKnight login page URL, used when
// authknight is enabled as a secondary login method (the primary method
// is always served at Login()).
func (l *authLinks) LoginAuthknight(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(AUTH_LOGIN_AUTHKNIGHT, p)
}

// LoginMagiclink returns the magic-link login page URL, used when
// magiclink is enabled as a secondary login method.
func (l *authLinks) LoginMagiclink(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(AUTH_LOGIN_MAGICLINK, p)
}

// LoginOtp returns the OTP login page URL, used when otp is enabled as a
// secondary login method.
func (l *authLinks) LoginOtp(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(AUTH_LOGIN_OTP, p)
}

// LoginPassword returns the password login page URL, used when password
// is enabled as a secondary login method.
func (l *authLinks) LoginPassword(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(AUTH_LOGIN_PASSWORD, p)
}

func (l *authLinks) Logout(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(AUTH_LOGOUT, p)
}

func (l *authLinks) PasswordReset(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(AUTH_PASSWORD_RESET, p)
}

func (l *authLinks) Register(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(AUTH_REGISTER, p)
}
