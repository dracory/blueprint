package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"

	"github.com/dracory/env"
	"github.com/samber/lo"
)

// authConfig reads authentication configuration from environment variables.
func authConfig() authSettings {
	// User Registration
	//
	// Controls whether new users can register for an account.
	// Set to false to disable public registration (invite-only or closed systems).
	registrationEnabled := env.GetBool(KEY_AUTH_REGISTRATION_ENABLED)

	// Allowed Emails
	//
	// Comma-separated list of emails allowed to access the application.
	// If empty, all authenticated emails are allowed.
	emailsAllowedAccess := lo.FilterMap(
		strings.FieldsFunc(env.GetString(KEY_AUTH_EMAILS_ALLOWED_ACCESS), func(r rune) bool {
			return r == ',' || r == ';'
		}),
		func(e string, _ int) (string, bool) {
			e = strings.TrimSpace(e)
			return e, e != ""
		},
	)

	if len(emailsAllowedAccess) == 0 {
		emailsAllowedAccess = []string{
			"info@sinevia.com",
			"lesichkovm@gmail.com",
		}
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
			panic(fmt.Sprintf("FATAL: %s must be set in the %s environment", KEY_AUTH_CSRF_SECRET, appEnv))
		}
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			panic(fmt.Sprintf("FATAL: failed to generate CSRF secret: %v", err))
		}
		csrfSecret = hex.EncodeToString(b)
		slog.Warn("AUTH_CSRF_SECRET is not set; a random secret has been generated for this run. " +
			"Set AUTH_CSRF_SECRET in your environment for persistent CSRF protection.")
	}

	// Password Authentication
	//
	// Controls whether email/password authentication is enabled.
	// When false, authentication is handled exclusively by external providers.
	passwordAuthEnabled := env.GetBool(KEY_AUTH_PASSWORD_AUTH_ENABLED)

	return authSettings{
		registrationEnabled: registrationEnabled,
		emailsAllowedAccess: emailsAllowedAccess,
		csrfSecret:          csrfSecret,
		passwordAuthEnabled: passwordAuthEnabled,
	}
}

type authSettings struct {
	registrationEnabled bool
	emailsAllowedAccess []string
	csrfSecret          string
	passwordAuthEnabled bool
}
