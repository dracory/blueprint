package emails

import (
	"project/internal/registry"

	"github.com/samber/lo"
)

func NewEmailNotifyAdmin(registry registry.RegistryInterface) *emailNotifyAdmin {
	return &emailNotifyAdmin{registry: registry}
}

type emailNotifyAdmin struct{ registry registry.RegistryInterface }

// Send sends an email notification to the admin
func (e *emailNotifyAdmin) Send(html string) error {
	appName := lo.IfF(e.registry != nil, func() string {
		if e.registry.GetConfig() == nil {
			return ""
		}
		return e.registry.GetConfig().GetAppName()
	}).Else("")

	fromEmail := lo.IfF(e.registry != nil, func() string {
		if e.registry.GetConfig() == nil {
			return ""
		}
		return e.registry.GetConfig().GetMailFromAddress()
	}).Else("")

	fromName := lo.IfF(e.registry != nil, func() string {
		if e.registry.GetConfig() == nil {
			return ""
		}
		return e.registry.GetConfig().GetMailFromName()
	}).Else("")

	emailSubject := appName + ". Admin Notification"
	emailContent := html

	// Use the new CreateEmailTemplate function instead of blankEmailTemplate
	finalHtml := CreateEmailTemplate(e.registry, emailSubject, emailContent)

	recipientEmail := "info@sinevia.com" // admin email

	// Use the new SendEmail function instead of Send
	errSend := SendEmail(SendOptions{
		From:     fromEmail,
		FromName: fromName,
		To:       []string{recipientEmail},
		Subject:  emailSubject,
		HtmlBody: finalHtml,
	})
	return errSend
}
