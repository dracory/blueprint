package emails

import (
	"context"
	"errors"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/hb"
	"github.com/dracory/userstore"
	"github.com/samber/lo"
)

func NewInviteFriendEmail(app types.AppInterface, us userstore.StoreInterface) *inviteFriendEmail {
	return &inviteFriendEmail{app: app, userStore: us}
}

type inviteFriendEmail struct {
	app       types.AppInterface
	userStore userstore.StoreInterface
}

// Send sends an invitation email to a friend
func (e *inviteFriendEmail) Send(sendingUserID string, userNote string, recipientEmail string, recipientName string) error {
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

	if e.userStore == nil {
		return errors.New("user store not configured")
	}

	user, err := e.userStore.UserFindByID(context.Background(), sendingUserID)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	userName := user.FirstName()

	if userName == "" {
		userName = user.Email()
	}

	emailSubject := appName + ". Invitation by a Friend"
	emailContent := e.template(appName, userName, userNote, recipientName)

	// Use the new CreateEmailTemplate function instead of blankEmailTemplate
	finalHtml := CreateEmailTemplate(e.app, emailSubject, emailContent)

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

func (e *inviteFriendEmail) template(appName string, userName string, userNote string, recipientName string) string {
	urlHome := hb.Hyperlink().Text(appName).
		Href(links.Website().Home()).ToHTML()

	urlJoin := hb.Hyperlink().Text("Click to Join Me at " + appName).
		Href(links.Website().Home()).ToHTML()

	h1 := hb.Heading1().
		HTML(`You have an awesome friend`).
		Style(STYLE_HEADING)

	p1 := hb.Paragraph().
		HTML(`Hi ` + recipientName + `,`).
		Style(STYLE_PARAGRAPH)

	p2 := hb.Paragraph().
		HTML(`You have been invited by a friend who thinks you will like ` + appName + `.`).
		Style(STYLE_PARAGRAPH)

	p3 := hb.Paragraph().
		HTML(`A note from your friend ` + userName + `:`).
		Style(STYLE_PARAGRAPH)

	p4 := hb.Paragraph().
		HTML(`"` + userNote + `"`).
		Style(STYLE_PARAGRAPH)

	p5 := hb.Paragraph().
		HTML(urlJoin).
		Style(STYLE_PARAGRAPH)

	p6 := hb.Paragraph().
		HTML(``). // Add description
		Style(STYLE_PARAGRAPH)

	p7 := hb.Paragraph().
		Children([]hb.TagInterface{
			hb.Raw(`Thank you for choosing ` + urlHome + `.`),
		}).
		Style(STYLE_PARAGRAPH)

	return hb.Div().Children([]hb.TagInterface{
		h1,
		p1,
		p2,
		p3,
		p4,
		p5,
		p6,
		hb.BR(),
		hb.BR(),
		p7,
	}).ToHTML()
}
