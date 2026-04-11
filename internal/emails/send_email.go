package emails

import (
	"fmt"
	"project/internal/registry"
	"sync"

	"github.com/dracory/email"
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

var (
	emailSender email.Sender
	senderMu    sync.RWMutex
)

// InitEmailSender initializes the email sender
func InitEmailSender(registry registry.RegistryInterface) {
	if registry == nil {
		return
	}

	if registry.GetConfig() == nil {
		return
	}

	sender := email.NewSMTPSender(email.Config{
		Host:     registry.GetConfig().GetMailHost(),
		Port:     cast.ToString(registry.GetConfig().GetMailPort()),
		Username: registry.GetConfig().GetMailUsername(),
		Password: registry.GetConfig().GetMailPassword(),
		// Skip logger for now as it's causing type compatibility issues
	})

	senderMu.Lock()
	defer senderMu.Unlock()
	emailSender = sender
}

// SendEmail sends an email using the base email package
// This is a new function to avoid conflicts with the original Send function
func SendEmail(options SendOptions) error {
	// Guard against nil sender
	senderMu.RLock()
	sender := emailSender
	senderMu.RUnlock()

	if sender == nil {
		return fmt.Errorf("email sender is not initialized")
	}

	// Convert SendOptions to base email.SendOptions
	baseOptions := email.SendOptions{
		From:     options.From,
		FromName: options.FromName,
		To:       options.To,
		Bcc:      options.Bcc,
		Cc:       options.Cc,
		Subject:  options.Subject,
		HtmlBody: options.HtmlBody,
		TextBody: options.TextBody,
	}

	return sender.Send(baseOptions)
}
