package emails

import (
	"html"
	"project/internal/app"
	"project/internal/links"

	"github.com/dracory/email"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

// NewEmailOtp creates a new OTP email sender.
func NewEmailOtp(app app.AppInterface) *emailOtp {
	return &emailOtp{app: app}
}

type emailOtp struct {
	app app.AppInterface
}

// Send sends a one-time-password login code to the given email address.
func (e *emailOtp) Send(recipientEmail string, otp string) error {
	appName := lo.IfF(e.app != nil && e.app.GetConfig() != nil, func() string {
		return e.app.GetConfig().GetAppName()
	}).Else("")

	fromEmail := lo.IfF(e.app != nil && e.app.GetConfig() != nil, func() string {
		return e.app.GetConfig().GetMailFromAddress()
	}).Else("")

	fromName := lo.IfF(e.app != nil && e.app.GetConfig() != nil, func() string {
		return e.app.GetConfig().GetMailFromName()
	}).Else("")

	emailSubject := appName + ". Your Login Code"
	emailContent := e.template(appName, otp)
	finalHtml := CreateEmailTemplate(e.app, emailSubject, emailContent)

	return SendEmail(SendOptions{
		From:     fromEmail,
		FromName: fromName,
		To:       []string{recipientEmail},
		Subject:  emailSubject,
		HtmlBody: finalHtml,
	})
}

func (e *emailOtp) template(appName string, otp string) string {
	urlHome := hb.Hyperlink().Text(appName).
		Href(links.Website().Home()).ToHTML()

	h1 := hb.Heading1().
		HTML(`Your login code`).
		Style(email.StyleHeading1)

	p1 := hb.Paragraph().
		HTML(`Use the following code to sign in to ` + appName + `:`).
		Style(email.StyleParagraph)

	code := hb.Div().
		HTML(html.EscapeString(otp)).
		Style(`font-size:32px;font-weight:bold;letter-spacing:8px;text-align:center;padding:16px;background:#f4f4f5;border-radius:8px;`)

	p2 := hb.Paragraph().
		HTML(`This code expires in 15 minutes. If you did not request it, you can safely ignore this email.`).
		Style(email.StyleParagraph)

	p3 := hb.Paragraph().
		Children([]hb.TagInterface{
			hb.Raw(`Thank you for choosing ` + urlHome + `.`),
		}).
		Style(email.StyleParagraph)

	return hb.Div().Children([]hb.TagInterface{
		h1,
		p1,
		code,
		hb.BR(),
		p2,
		hb.BR(),
		p3,
	}).ToHTML()
}
