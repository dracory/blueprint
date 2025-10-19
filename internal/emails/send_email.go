package emails

import (
	"fmt"
	"project/internal/types"

	baseEmail "github.com/dracory/base/email"
	"github.com/spf13/cast"
)

type Mailer struct {
	Host     string
	Port     int
	Username string
	Password string
}

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
func InitEmailSender(app types.AppInterface) {
	if app == nil {
		return
	}

	if app.GetConfig() == nil {
		return
	}

	emailSender = baseEmail.NewSMTPSender(baseEmail.Config{
		Host:     app.GetConfig().GetMailHost(),
		Port:     cast.ToString(app.GetConfig().GetMailPort()),
		Username: app.GetConfig().GetMailUsername(),
		Password: app.GetConfig().GetMailPassword(),
		// Skip logger for now as it's causing type compatibility issues
	})
}

// SendEmail sends an email using the base email package
// This is a new function to avoid conflicts with the original Send function
func SendEmail(options SendOptions) error {
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
