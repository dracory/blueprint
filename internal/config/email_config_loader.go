package config

import (
	"github.com/dracory/env"
	"github.com/spf13/cast"
)

// loadMailConfig loads mail configuration directly into the config.
func loadMailConfig(cfg ConfigInterface) {
	// Mail Driver
	//
	// The mail driver to use for sending emails.
	// Supported values: smtp, sendgrid, mailgun
	// Leave empty to disable email sending.
	cfg.SetMailDriver(env.GetString(KEY_MAIL_DRIVER))

	// Mail From Address
	//
	// The email address that emails will be sent from.
	// Example: noreply@example.com
	cfg.SetMailFromAddress(env.GetString(KEY_MAIL_FROM_ADDRESS))

	// Mail From Name
	//
	// The name that will appear as the sender in outgoing emails.
	// Example: My Application
	cfg.SetMailFromName(env.GetString(KEY_MAIL_FROM_NAME))

	// Mail Host
	//
	// The hostname of the SMTP server.
	// Example: smtp.gmail.com, smtp.sendgrid.net
	cfg.SetMailHost(env.GetString(KEY_MAIL_HOST))

	// Mail Password
	//
	// The password for authenticating with the mail server.
	cfg.SetMailPassword(env.GetString(KEY_MAIL_PASSWORD))

	// Mail Port
	//
	// The port the mail server is listening on.
	// Common values: 25 (SMTP), 465 (SMTPS), 587 (STARTTLS)
	cfg.SetMailPort(cast.ToInt(env.GetString(KEY_MAIL_PORT)))

	// Mail Username
	//
	// The username for authenticating with the mail server.
	cfg.SetMailUsername(env.GetString(KEY_MAIL_USERNAME))
}
