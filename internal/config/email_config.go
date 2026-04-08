package config

import (
	"github.com/dracory/env"
	"github.com/spf13/cast"
)

// ============================================================================
// Interface
// ============================================================================

// EmailConfigInterface defines email/mail configuration methods.
type EmailConfigInterface interface {
	SetMailDriver(string)
	GetMailDriver() string

	SetMailHost(string)
	GetMailHost() string

	SetMailPort(int)
	GetMailPort() int

	SetMailUsername(string)
	GetMailUsername() string

	SetMailPassword(string)
	GetMailPassword() string

	SetMailFromAddress(string)
	GetMailFromAddress() string

	SetMailFromName(string)
	GetMailFromName() string
}

// ============================================================================
// Types
// ============================================================================

// mailConfig captures email delivery settings.
type mailConfig struct {
	driver      string // Mail driver/service provider (smtp, sendgrid, mailgun, etc.)
	fromAddress string // Default sender email address
	fromName    string // Default sender name
	host        string // SMTP server hostname
	password    string // SMTP or service provider password
	port        int    // SMTP server port number
	username    string // SMTP or service provider username
}

// ============================================================================
// Loader
// ============================================================================

// loadMailConfig loads mail configuration from environment variables.
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

// ============================================================================
// Implementation (Getters/Setters)
// ============================================================================

func (c *configImplementation) SetMailDriver(v string) {
	c.emailDriver = v
}

func (c *configImplementation) GetMailDriver() string {
	return c.emailDriver
}

func (c *configImplementation) SetMailHost(v string) {
	c.emailHost = v
}

func (c *configImplementation) GetMailHost() string {
	return c.emailHost
}

func (c *configImplementation) SetMailPort(v int) {
	c.emailPort = v
}

func (c *configImplementation) GetMailPort() int {
	return c.emailPort
}

func (c *configImplementation) SetMailUsername(v string) {
	c.emailUsername = v
}

func (c *configImplementation) GetMailUsername() string {
	return c.emailUsername
}

func (c *configImplementation) SetMailPassword(v string) {
	c.emailPassword = v
}

func (c *configImplementation) GetMailPassword() string {
	return c.emailPassword
}

func (c *configImplementation) SetMailFromName(v string) {
	c.emailFromName = v
}

func (c *configImplementation) GetMailFromName() string {
	return c.emailFromName
}

func (c *configImplementation) SetMailFromAddress(v string) {
	c.emailFromAddress = v
}

func (c *configImplementation) GetMailFromAddress() string {
	return c.emailFromAddress
}
