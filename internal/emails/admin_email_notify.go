package emails

import (
	"project/internal/types"
)

func NewEmailNotifyAdmin(cfg types.ConfigInterface) *emailNotifyAdmin {
	return &emailNotifyAdmin{cfg: cfg}
}

type emailNotifyAdmin struct{ cfg types.ConfigInterface }

// Send sends an email notification to the admin
func (e *emailNotifyAdmin) Send(html string) error {
	appName := e.cfg.GetAppName()
	emailSubject := appName + ". Admin Notification"
	emailContent := html

	// Use the new CreateEmailTemplate function instead of blankEmailTemplate
	finalHtml := CreateEmailTemplate(emailSubject, emailContent)

	recipientEmail := "info@sinevia.com"

	// Use the new SendEmail function instead of Send
	errSend := SendEmail(SendOptions{
		From:     e.cfg.GetMailFromAddress(),
		FromName: appName,
		To:       []string{recipientEmail},
		Subject:  emailSubject,
		HtmlBody: finalHtml,
	})
	return errSend
}
