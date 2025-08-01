package emails

import (
	"project/internal/config"

	baseEmail "github.com/dracory/base/email"
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

// InitEmailSender initializes the email sender
func InitEmailSender() {
	emailSender = baseEmail.NewSMTPSender(baseEmail.Config{
		Host:     config.MailHost,
		Port:     config.MailPort,
		Username: config.MailUsername,
		Password: config.MailPassword,
		// Skip logger for now as it's causing type compatibility issues
	})
}

// SendEmail sends an email using the base email package
// This is a new function to avoid conflicts with the original Send function
func SendEmail(options SendOptions) error {
	// Initialize the email sender if it hasn't been initialized yet
	if emailSender == nil {
		InitEmailSender()
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
		Title:       title,
		Content:     htmlContent,
		AppName:     config.AppName,
		HeaderLinks: headerLinks,
	})
}
