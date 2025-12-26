package emails

import (
	"project/internal/types"

	"github.com/samber/lo"
)

func NewEmailNotifyAdmin(app types.AppInterface) *emailNotifyAdmin {
	return &emailNotifyAdmin{app: app}
}

type emailNotifyAdmin struct{ app types.AppInterface }

// Send sends an email notification to the admin
func (e *emailNotifyAdmin) Send(html string) error {
	appName := lo.IfF(e.app != nil, func() string {
		if e.app.GetConfig() == nil {
			return ""
		}
		return e.app.GetConfig().GetAppName()
	}).Else("")

	fromEmail := lo.IfF(e.app != nil, func() string {
		if e.app.GetConfig() == nil {
			return ""
		}
		return e.app.GetConfig().GetMailFromAddress()
	}).Else("")

	fromName := lo.IfF(e.app != nil, func() string {
		if e.app.GetConfig() == nil {
			return ""
		}
		return e.app.GetConfig().GetMailFromName()
	}).Else("")

	emailSubject := appName + ". Admin Notification"
	emailContent := html

	// Use the new CreateEmailTemplate function instead of blankEmailTemplate
	finalHtml := CreateEmailTemplate(e.app, emailSubject, emailContent)

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
