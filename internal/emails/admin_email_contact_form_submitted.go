package emails

import (
	"project/internal/config"
	"project/internal/links"

	"github.com/gouniverse/hb"
)

func NewEmailToAdminOnNewContactFormSubmitted() *emailToAdminOnNewContactFormSubmitted {
	return &emailToAdminOnNewContactFormSubmitted{}
}

type emailToAdminOnNewContactFormSubmitted struct{}

// Send sends an email notification to the admin when a new contact form is submitted
func (e *emailToAdminOnNewContactFormSubmitted) Send() error {
	emailSubject := config.AppName + ". New Contact Form Submitted"
	emailContent := e.template()

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

func (e *emailToAdminOnNewContactFormSubmitted) template() string {
	urlHome := hb.Hyperlink().
		HTML(config.AppName).
		Href(links.NewWebsiteLinks().Home()).
		ToHTML()

	h1 := hb.Heading1().
		HTML(`New Contact Form Submitted`).
		Style(STYLE_HEADING)

	p1 := hb.Paragraph().
		HTML(`There is a new contact form request submitted into ` + config.AppName + `.`).
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
