package config

import (
	"github.com/dracory/env"
	"github.com/spf13/cast"
)

// readMailConfig reads mail configuration from environment variables.
func emailConfig(cfg *configImplementation) {
	// Mail Driver
	//
	// The mail driver to use for sending emails.
	// Supported values: smtp, sendgrid, mailgun
	// Leave empty to disable email sending.
	driver := env.GetString(KEY_MAIL_DRIVER)

	// Mail From Address
	//
	// The email address that emails will be sent from.
	// Example: noreply@example.com
	fromAddress := env.GetStringOrDefault(KEY_MAIL_FROM_ADDRESS, "noreply@example.com")

	// Mail From Name
	//
	// The name that will appear as the sender in outgoing emails.
	fromName := env.GetStringOrDefault(KEY_MAIL_FROM_NAME, "Blueprint")

	// Mail Host
	//
	// The hostname of the SMTP server.
	// Example: smtp.gmail.com, smtp.sendgrid.net
	host := env.GetString(KEY_MAIL_HOST)

	// Mail Password
	//
	// The password for authenticating with the mail server.
	password := env.GetString(KEY_MAIL_PASSWORD)

	// Mail Port
	//
	// The port the mail server is listening on.
	// Common values: 25 (SMTP), 465 (SMTPS), 587 (STARTTLS)
	port := cast.ToInt(env.GetStringOrDefault(KEY_MAIL_PORT, "587"))

	// Mail Username
	//
	// The username for authenticating with the mail server.
	username := env.GetString(KEY_MAIL_USERNAME)

	// -------------------------------------------------------------------------
	// Do not edit below this line
	// -------------------------------------------------------------------------
	cfg.setMailConfig(driver, fromAddress, fromName, host, password, port, username)
}
