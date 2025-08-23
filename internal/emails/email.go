package emails

import (
	"fmt"
	"project/internal/types"

	baseEmail "github.com/dracory/base/email"
	"github.com/spf13/cast"
)

// SendOptions defines the options for sending an email
// This maintains compatibility with the original implementation
type SendOptions struct {
	From     string
	FromName string
	To       []string
	Bcc      []string
	Cc       []string
	Subject  string
	HtmlBody string
	TextBody string
}

var emailSender baseEmail.Sender
var cfg types.ConfigInterface

// InitEmailSender initializes the email sender
func InitEmailSender() {
	if cfg == nil {
		// cannot initialize without config; leave sender nil so caller gets an error on send
		return
	}
	emailSender = baseEmail.NewSMTPSender(baseEmail.Config{
		Host:     cfg.GetMailHost(),
		Port:     cast.ToString(cfg.GetMailPort()),
		Username: cfg.GetMailUsername(),
		Password: cfg.GetMailPassword(),
		// Skip logger for now as it's causing type compatibility issues
	})
}

// Init sets the package-level configuration for emails
func Init(c types.ConfigInterface) {
	cfg = c
}

// SendEmail sends an email using the base email package
// This is a new function to avoid conflicts with the original Send function
func SendEmail(options SendOptions) error {
	// Initialize the email sender if it hasn't been initialized yet
	if emailSender == nil {
		InitEmailSender()
	}

	// Guard against nil sender
	if emailSender == nil {
		return fmt.Errorf("email sender is not initialized")
	}

	// Convert SendOptions to base email.SendOptions
	baseOptions := baseEmail.SendOptions{
		From:     options.From,
		FromName: options.FromName,
		To:       options.To,
		Bcc:      options.Bcc,
		Cc:       options.Cc,
		Subject:  options.Subject,
		HtmlBody: options.HtmlBody,
		TextBody: options.TextBody,
	}

	return emailSender.Send(baseOptions)
}

// CreateEmailTemplate creates an email template using the base email package
// This is a new function to avoid conflicts with the original blankEmailTemplate function
func CreateEmailTemplate(title string, htmlContent string) string {
	// Create header links
	headerLinks := map[string]string{}

	// Use the base email template
	return baseEmail.DefaultTemplate(baseEmail.TemplateOptions{
		Title:   title,
		Content: htmlContent,
		AppName: func() string {
			if cfg != nil {
				return cfg.GetAppName()
			}
			return ""
		}(),
		HeaderLinks: headerLinks,
	})
}
