package emails

import (
	"context"
	"errors"
	"project/internal/types"
	"project/internal/links"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/userstore"
)

func NewInviteFriendEmail(cfg types.ConfigInterface, us userstore.StoreInterface) *inviteFriendEmail {
	return &inviteFriendEmail{cfg: cfg, userStore: us}
}

type inviteFriendEmail struct {
	cfg       types.ConfigInterface
	userStore userstore.StoreInterface
}

// Send sends an invitation email to a friend
func (e *inviteFriendEmail) Send(sendingUserID string, userNote string, recipientEmail string, recipientName string) error {
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

	appName := ""
	if e.cfg != nil {
		appName = e.cfg.GetAppName()
	}
	emailSubject := appName + ". Invitation by a Friend"
	emailContent := e.template(userName, userNote, recipientName)

	// Use the new CreateEmailTemplate function instead of blankEmailTemplate
	finalHtml := CreateEmailTemplate(emailSubject, emailContent)

	// Use the new SendEmail function instead of Send
	errSend := SendEmail(SendOptions{
		From:     func() string { if e.cfg != nil { return e.cfg.GetMailFromEmail() }; return "" }(),
		FromName: func() string { if e.cfg != nil { return e.cfg.GetMailFromName() }; return appName }(),
		To:       []string{recipientEmail},
		Subject:  emailSubject,
		HtmlBody: finalHtml,
	})
	return errSend
}

func (e *inviteFriendEmail) template(userName string, userNote string, recipientName string) string {

	urlHome := hb.Hyperlink().Text("ProvedExpert").
		Href(links.NewWebsiteLinks().Home()).ToHTML()

	urlJoin := hb.Hyperlink().Text("Click to Join Me at ProvedExpert").
		Href(links.NewWebsiteLinks().Home()).ToHTML()

	h1 := hb.Heading1().
		HTML(`You have an awesome friend`).
		Style(STYLE_HEADING)

	p1 := hb.Paragraph().
		HTML(`Hi ` + recipientName + `,`).
		Style(STYLE_PARAGRAPH)

	p2 := hb.Paragraph().
		HTML(`You have been invited by a friend who thinks you will like ` + func() string { if e.cfg != nil { return e.cfg.GetAppName() }; return "" }() + `.`).
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
