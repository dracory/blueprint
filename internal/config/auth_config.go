package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
)

// authConfig reads authentication configuration from environment variables.
// Validation failures are collected in env and reported together by NewFromEnv.
func authConfig(env *envValidator) authSettings {
	// Login Method
	//
	// Selects which login mechanism is mounted at links.AUTH_LOGIN.
	// Valid values: otp (default), magiclink, password, authknight.
	loginMethod := strings.ToLower(strings.TrimSpace(env.GetString(KEY_AUTH_LOGIN_METHOD)))
	if loginMethod == "" {
		loginMethod = LOGIN_METHOD_OTP
	}
	switch loginMethod {
	case LOGIN_METHOD_OTP, LOGIN_METHOD_MAGICLINK, LOGIN_METHOD_PASSWORD, LOGIN_METHOD_AUTHKNIGHT:
	default:
		env.Add(fmt.Errorf("invalid %s value %q (expected one of: %s, %s, %s, %s)",
			KEY_AUTH_LOGIN_METHOD, loginMethod,
			LOGIN_METHOD_OTP, LOGIN_METHOD_MAGICLINK, LOGIN_METHOD_PASSWORD, LOGIN_METHOD_AUTHKNIGHT))
	}

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
		loginMethod:         loginMethod,
	}
}

type authSettings struct {
	registrationEnabled bool
	emailsAllowedAccess []string
	csrfSecret          string
	loginMethod         string
}
