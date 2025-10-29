package emails

import (
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func NewEmailToAdminOnNewUserRegistered(app types.AppInterface) *emailToAdminOnNewUserRegistered {
	return &emailToAdminOnNewUserRegistered{app: app}
}

type emailToAdminOnNewUserRegistered struct{ app types.AppInterface }

// Send sends an email notification to the admin when a new user registers
func (e *emailToAdminOnNewUserRegistered) Send(userID string) error {
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

	emailSubject := appName + ". New User Registered"
	emailContent := e.template(appName, userID)

	// Use the new CreateEmailTemplate function instead of blankEmailTemplate
	finalHtml := CreateEmailTemplate(e.app, emailSubject, emailContent)

	recipientEmail := "info@sinevia.com"

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

func (e *emailToAdminOnNewUserRegistered) template(appName string, userID string) string {
	urlHome := hb.Hyperlink().
		HTML(appName).
		Href(links.Website().Home()).
		ToHTML()

	h1 := hb.Heading1().
		HTML(`New User Registsred`).
		Style(STYLE_HEADING)

	p1 := hb.Paragraph().
		HTML(`There is a new user ID ` + userID + `	that registsred into ` + appName + `.`).
		Style(STYLE_PARAGRAPH)

	p2 := hb.Paragraph().
		HTML(`Please login to admin panel to check the new user.`).
		Style(STYLE_PARAGRAPH)

	p6 := hb.Paragraph().
		Children([]hb.TagInterface{
			hb.Text(`Thank you for choosing ` + urlHome + `.`),
			hb.BR(),
			hb.Text(`The new way to learn`),
		}).
		Style(STYLE_PARAGRAPH)

	return hb.Div().Children([]hb.TagInterface{
		h1,
		p1,
		p2,
		hb.BR(),
		hb.BR(),
		p6,
	}).ToHTML()
}
