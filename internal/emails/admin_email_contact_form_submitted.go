package emails

import (
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func NewEmailToAdminOnNewContactFormSubmitted(app types.AppInterface) *emailToAdminOnNewContactFormSubmitted {
	return &emailToAdminOnNewContactFormSubmitted{app: app}
}

type emailToAdminOnNewContactFormSubmitted struct {
	app types.AppInterface
}

// Send sends an email notification to the admin when a new contact form is submitted
func (e *emailToAdminOnNewContactFormSubmitted) Send() error {
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

	emailSubject := appName + ". New Contact Form Submitted"
	emailContent := e.template(appName)

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

func (e *emailToAdminOnNewContactFormSubmitted) template(appName string) string {
	urlHome := hb.Hyperlink().
		HTML(appName).
		Href(links.Website().Home()).
		ToHTML()

	h1 := hb.Heading1().
		HTML(`New Contact Form Submitted`).
		Style(STYLE_HEADING)

	p1 := hb.Paragraph().
		HTML(`There is a new contact form request submitted into ` + appName + `.`).
		Style(STYLE_PARAGRAPH)

	p2 := hb.Paragraph().
		HTML(`Please login to admin panel to check the new contact request.`).
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
