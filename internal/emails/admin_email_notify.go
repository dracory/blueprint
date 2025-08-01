package emails

import (
	"project/internal/config"
)

func NewEmailNotifyAdmin() *emailNotifyAdmin {
	return &emailNotifyAdmin{}
}

type emailNotifyAdmin struct{}

// Send sends an email notification to the admin
func (e *emailNotifyAdmin) Send(html string) error {
	emailSubject := config.AppName + ". Admin Notification"
	emailContent := html

	// Use the new CreateEmailTemplate function instead of blankEmailTemplate
	finalHtml := CreateEmailTemplate(emailSubject, emailContent)

	recipientEmail := "info@sinevia.com"

	// Use the new SendEmail function instead of Send
	errSend := SendEmail(SendOptions{
		From:     config.MailFromEmailAddress,
		FromName: config.AppName,
		To:       []string{recipientEmail},
		Subject:  emailSubject,
		HtmlBody: finalHtml,
	})
	return errSend
}
