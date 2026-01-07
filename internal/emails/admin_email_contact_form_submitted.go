package emails

import (
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func NewEmailToAdminOnNewContactFormSubmitted(registry registry.RegistryInterface) *emailToAdminOnNewContactFormSubmitted {
	return &emailToAdminOnNewContactFormSubmitted{registry: registry}
}

type emailToAdminOnNewContactFormSubmitted struct {
	registry registry.RegistryInterface
}

// Send sends an email notification to the admin when a new contact form is submitted
func (e *emailToAdminOnNewContactFormSubmitted) Send() error {
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

	emailSubject := appName + ". New Contact Form Submitted"
	emailContent := e.template(appName)

	// Use the new CreateEmailTemplate function instead of blankEmailTemplate
	finalHtml := CreateEmailTemplate(e.registry, emailSubject, emailContent)

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
