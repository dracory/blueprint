package emails

import (
	"html"
	"project/internal/app"
	"project/internal/links"

	"github.com/dracory/email"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

// NewEmailMagicLink creates a new magic link email sender.
func NewEmailMagicLink(app app.AppInterface) *emailMagicLink {
	return &emailMagicLink{app: app}
}

type emailMagicLink struct {
	app app.AppInterface
}

// Send sends a single-use login link to the given email address.
func (e *emailMagicLink) Send(recipientEmail string, magicLinkURL string) error {
	appName := lo.IfF(e.app != nil && e.app.GetConfig() != nil, func() string {
		return e.app.GetConfig().GetAppName()
	}).Else("")

	fromEmail := lo.IfF(e.app != nil && e.app.GetConfig() != nil, func() string {
		return e.app.GetConfig().GetMailFromAddress()
	}).Else("")

	fromName := lo.IfF(e.app != nil && e.app.GetConfig() != nil, func() string {
		return e.app.GetConfig().GetMailFromName()
	}).Else("")

	emailSubject := appName + ". Your Login Link"
	emailContent := e.template(appName, magicLinkURL)
	finalHtml := CreateEmailTemplate(e.app, emailSubject, emailContent)

	return SendEmail(SendOptions{
		From:     fromEmail,
		FromName: fromName,
		To:       []string{recipientEmail},
		Subject:  emailSubject,
		HtmlBody: finalHtml,
	})
}

func (e *emailMagicLink) template(appName string, magicLinkURL string) string {
	urlHome := hb.Hyperlink().Text(appName).
		Href(links.Website().Home()).ToHTML()

	h1 := hb.Heading1().
		HTML(`Your login link`).
		Style(email.StyleHeading1)

	p1 := hb.Paragraph().
		HTML(`Click the button below to sign in to ` + appName + `:`).
		Style(email.StyleParagraph)

	button := hb.Div().Style(`text-align:center;padding:16px;`).
		Child(hb.Hyperlink().
			Href(magicLinkURL).
			Text("Sign In").
			Style(`display:inline-block;padding:14px 32px;background:#667eea;color:#ffffff;text-decoration:none;font-weight:bold;border-radius:8px;`))

	p2 := hb.Paragraph().
		HTML(`This link expires in 15 minutes and can only be used once. If you did not request it, you can safely ignore this email.`).
		Style(email.StyleParagraph)

	p3 := hb.Paragraph().
		HTML(`If the button does not work, copy and paste this URL into your browser:`).
		Style(email.StyleParagraph)

	p4 := hb.Paragraph().
		HTML(html.EscapeString(magicLinkURL)).
		Style(`word-break:break-all;font-size:12px;color:#6b7280;`)

	p5 := hb.Paragraph().
		Children([]hb.TagInterface{
			hb.Raw(`Thank you for choosing ` + urlHome + `.`),
		}).
		Style(email.StyleParagraph)

	return hb.Div().Children([]hb.TagInterface{
		h1,
		p1,
		button,
		hb.BR(),
		p2,
		p3,
		p4,
		hb.BR(),
		p5,
	}).ToHTML()
}
