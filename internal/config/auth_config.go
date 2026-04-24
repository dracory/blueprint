package config

import (
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

	return authSettings{
		registrationEnabled: registrationEnabled,
		emailsAllowedAccess: emailsAllowedAccess,
	}
}

type authSettings struct {
	registrationEnabled bool
	emailsAllowedAccess []string
}
