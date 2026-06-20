package authrules

import (
	"project/internal/app"

	"github.com/dracory/rule"
)

// EmailAllowedRule checks whether a given email is permitted to access the
// application. If the configured allowlist is empty, all emails are allowed.
type EmailAllowedRule struct {
	rule.Rule
}

type emailAllowedContext struct {
	allowedEmails []string
	email         string
}

// NewEmailAllowedRule creates an EmailAllowedRule for the given app and email.
func NewEmailAllowedRule(a app.AppInterface, email string) *EmailAllowedRule {
	r := &EmailAllowedRule{}

	r.Rule.SetContext(emailAllowedContext{
		allowedEmails: a.GetConfig().GetEmailsAllowedAccess(),
		email:         email,
	})

	r.Rule.SetCondition(func(ctx any) bool {
		data := ctx.(emailAllowedContext)

		if len(data.allowedEmails) == 0 {
			return true
		}

		for _, allowed := range data.allowedEmails {
			if allowed == data.email {
				return true
			}
		}

		r.AddFailMessage("Your email is not permitted to access this application")
		return false
	})

	return r
}
