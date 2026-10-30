package authadapters

import (
	"context"

	"project/internal/app"
	"project/internal/emails"
)

// FuncEmailSend adapts the app's generic email sender to the auth
// package's passwordless email hook signature.
func FuncEmailSend(application app.AppInterface) func(ctx context.Context, email string, subject string, body string) error {
	return func(ctx context.Context, email string, subject string, body string) error {
		return emails.SendEmail(emails.SendOptions{
			From:     application.GetConfig().GetMailFromAddress(),
			FromName: application.GetConfig().GetMailFromName(),
			To:       []string{email},
			Subject:  subject,
			HtmlBody: body,
		})
	}
}
