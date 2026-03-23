package config

import (
	"github.com/dracory/env"
	"github.com/spf13/cast"
)

// mailConfig captures email delivery settings.
// It includes all necessary configuration for sending emails through
// various mail drivers (SMTP, SendGrid, Mailgun, etc.).
type mailConfig struct {
	driver      string // Mail driver/service provider (smtp, sendgrid, mailgun, etc.)
	fromAddress string // Default sender email address
	fromName    string // Default sender name
	host        string // SMTP server hostname
	password    string // SMTP or service provider password
	port        int    // SMTP server port number
	username    string // SMTP or service provider username
}

// loadMailConfig loads mail configuration from environment variables.
// It reads all mail-related environment variables and returns a populated
// mailConfig struct. The port is converted from string to int using cast.
//
// Returns:
//   - mailConfig: Populated configuration struct with mail delivery settings
func loadMailConfig() mailConfig {
	return mailConfig{
		driver:      env.GetString(KEY_MAIL_DRIVER),
		fromAddress: env.GetString(KEY_MAIL_FROM_ADDRESS),
		fromName:    env.GetString(KEY_MAIL_FROM_NAME),
		host:        env.GetString(KEY_MAIL_HOST),
		password:    env.GetString(KEY_MAIL_PASSWORD),
		port:        cast.ToInt(env.GetString(KEY_MAIL_PORT)),
		username:    env.GetString(KEY_MAIL_USERNAME),
	}
}
