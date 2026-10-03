package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"

	"github.com/samber/lo"
)

// authConfig reads authentication configuration from environment variables.
// Validation failures are collected in env and reported together by NewFromEnv.
func authConfig(env *envValidator) authSettings {
	loginMethods := authLoginMethods(env)

	// User Registration
	//
	// Controls whether new users can register for an account.
	// Set to false to disable public registration (invite-only or closed systems).
	registrationEnabled := env.GetBool(KEY_AUTH_REGISTRATION_ENABLED)

	// Allowed Emails
	//
	// Comma-separated list of emails allowed to access the application.
	// If empty, all authenticated emails are allowed.
	emailsAllowedAccess := env.GetArray(KEY_AUTH_EMAILS_ALLOWED_ACCESS)

	// Remember Me
	//
	// Opt-in persistent login via a long-lived remember session stored in a
	// separate HttpOnly cookie. AUTH_REMEMBER_ME_DAYS controls the lifetime
	// of both the session and the cookie (default 30 days).
	rememberMeEnabled := env.GetBool(KEY_AUTH_REMEMBER_ME_ENABLED)
	rememberMeDays := env.GetIntOrDefault(KEY_AUTH_REMEMBER_ME_DAYS, 30)
	if rememberMeDays < 1 {
		env.Add(fmt.Errorf("%s must be a positive number of days", KEY_AUTH_REMEMBER_ME_DAYS))
		rememberMeDays = 30
	}

	// CSRF Secret
	//
	// Secret key used for CSRF token generation/validation.
	// Required in production/staging. In other environments a random secret is
	// generated automatically and a warning is logged.
	csrfSecret := env.GetString(KEY_AUTH_CSRF_SECRET)
	if csrfSecret == "" {
		appEnv := env.GetString(KEY_APP_ENVIRONMENT)
		if appEnv == APP_ENVIRONMENT_PRODUCTION || appEnv == APP_ENVIRONMENT_STAGING {
			env.Add(fmt.Errorf("%s must be set in the %s environment", KEY_AUTH_CSRF_SECRET, appEnv))
		} else {
			b := make([]byte, 32)
			if _, err := rand.Read(b); err != nil {
				env.Add(fmt.Errorf("failed to generate CSRF secret: %w", err))
			} else {
				csrfSecret = hex.EncodeToString(b)
				slog.Warn("AUTH_CSRF_SECRET is not set; a random secret has been generated for this run. " +
					"Set AUTH_CSRF_SECRET in your environment for persistent CSRF protection.")
			}
		}
	}

	return authSettings{
		registrationEnabled: registrationEnabled,
		emailsAllowedAccess: emailsAllowedAccess,
		csrfSecret:          csrfSecret,
		loginMethods:        loginMethods,
		rememberMeEnabled:   rememberMeEnabled,
		rememberMeDays:      rememberMeDays,
	}
}

// authLoginMethods reads the enabled login mechanisms from
// AUTH_LOGIN_METHODS (comma/semicolon-separated, or a JSON array). The
// first entry is the primary method rendered at links.AUTH_LOGIN; the
// rest are offered as alternatives at /auth/<method>-login.
// Valid values: otp (default), magiclink, password, authknight.
// Invalid entries are collected in env.
func authLoginMethods(env *envValidator) []string {
	loginMethods := env.GetArrayLower(KEY_AUTH_LOGIN_METHODS)
	if len(loginMethods) == 0 {
		loginMethods = env.GetArrayLower(KEY_AUTH_LOGIN_METHOD) // legacy singular
	} else if len(env.GetArrayLower(KEY_AUTH_LOGIN_METHOD)) != 0 {
		slog.Warn("Both AUTH_LOGIN_METHODS and AUTH_LOGIN_METHOD are set; " +
			"AUTH_LOGIN_METHODS takes precedence and AUTH_LOGIN_METHOD is deprecated.")
	}
	if len(loginMethods) == 0 {
		loginMethods = []string{LOGIN_METHOD_OTP}
	}
	for _, m := range loginMethods {
		switch m {
		case LOGIN_METHOD_OTP, LOGIN_METHOD_MAGICLINK, LOGIN_METHOD_PASSWORD, LOGIN_METHOD_AUTHKNIGHT:
		default:
			env.Add(fmt.Errorf("invalid %s value %q (expected a comma-separated list of: %s, %s, %s, %s)",
				KEY_AUTH_LOGIN_METHODS, m,
				LOGIN_METHOD_OTP, LOGIN_METHOD_MAGICLINK, LOGIN_METHOD_PASSWORD, LOGIN_METHOD_AUTHKNIGHT))
		}
	}
	return lo.Uniq(loginMethods)
}

type authSettings struct {
	registrationEnabled bool
	emailsAllowedAccess []string
	csrfSecret          string
	loginMethods        []string
	rememberMeEnabled   bool
	rememberMeDays      int
}
