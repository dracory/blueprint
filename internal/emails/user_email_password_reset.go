package emails

import (
	"html"
	"project/internal/app"
	"project/internal/links"

	"github.com/dracory/email"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

// NewEmailPasswordReset creates a new password reset email sender.
func NewEmailPasswordReset(app app.AppInterface) *emailPasswordReset {
	return &emailPasswordReset{app: app}
}

type emailPasswordReset struct {
	app app.AppInterface
}

// Send sends a single-use password reset link to the given email address.
func (e *emailPasswordReset) Send(recipientEmail string, resetURL string) error {
	appName := lo.IfF(e.app != nil && e.app.GetConfig() != nil, func() string {
		return e.app.GetConfig().GetAppName()
	}).Else("")

	fromEmail := lo.IfF(e.app != nil && e.app.GetConfig() != nil, func() string {
		return e.app.GetConfig().GetMailFromAddress()
	}).Else("")

	fromName := lo.IfF(e.app != nil && e.app.GetConfig() != nil, func() string {
		return e.app.GetConfig().GetMailFromName()
	}).Else("")

	emailSubject := appName + ". Your Password Reset Link"
	emailContent := e.template(appName, resetURL)
	finalHtml := CreateEmailTemplate(e.app, emailSubject, emailContent)

	return SendEmail(SendOptions{
		From:     fromEmail,
		FromName: fromName,
		To:       []string{recipientEmail},
		Subject:  emailSubject,
		HtmlBody: finalHtml,
	})
}

func (e *emailPasswordReset) template(appName string, resetURL string) string {
	urlHome := hb.Hyperlink().Text(appName).
		Href(links.Website().Home()).ToHTML()

	h1 := hb.Heading1().
		HTML(`Reset your password`).
		Style(email.StyleHeading1)

	p1 := hb.Paragraph().
		HTML(`We received a request to reset the password for your ` + appName + ` account. Click the button below to choose a new password:`).
		Style(email.StyleParagraph)

	button := hb.Div().Style(`text-align:center;padding:16px;`).
		Child(hb.Hyperlink().
			Href(resetURL).
			Text("Reset Password").
			Style(`display:inline-block;padding:14px 32px;background:#667eea;color:#ffffff;text-decoration:none;font-weight:bold;border-radius:8px;`))

	p2 := hb.Paragraph().
		HTML(`This link expires in 15 minutes and can only be used once. If you did not request a password reset, you can safely ignore this email — your password will remain unchanged.`).
		Style(email.StyleParagraph)

	p3 := hb.Paragraph().
		HTML(`If the button does not work, copy and paste this URL into your browser:`).
		Style(email.StyleParagraph)

	p4 := hb.Paragraph().
		HTML(html.EscapeString(resetURL)).
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
