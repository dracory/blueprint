package emails

import (
	"project/app/links"
	"project/config"

	"github.com/gouniverse/hb"
)

func NewEmailToAdminOnNewUserRegistered() *emailToAdminOnNewUserRegistered {
	return &emailToAdminOnNewUserRegistered{}
}

type emailToAdminOnNewUserRegistered struct{}

// Send sends an email notification to the admin when a new user registers
func (e *emailToAdminOnNewUserRegistered) Send(userID string) error {
	emailSubject := config.AppName + ". New User Registered"
	emailContent := e.template(userID)

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

func (e *emailToAdminOnNewUserRegistered) template(userID string) string {
	urlHome := hb.Hyperlink().
		HTML(config.AppName).
		Href(links.NewWebsiteLinks().Home()).
		ToHTML()

	h1 := hb.Heading1().
		HTML(`New User Registsred`).
		Style(STYLE_HEADING)

	p1 := hb.Paragraph().
		HTML(`There is a new user ID ` + userID + `	that registsred into ` + config.AppName + `.`).
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
