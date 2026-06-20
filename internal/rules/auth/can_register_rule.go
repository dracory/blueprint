package authrules

import (
	"project/internal/app"

	"github.com/dracory/rule"
)

// CanRegisterRule checks whether registration is currently enabled and,
// when an email is provided, whether that email is on the allowed access list.
type CanRegisterRule struct {
	rule.Rule
}

type canRegisterContext struct {
	registrationEnabled bool
	allowedEmails       []string
	email               string
}

// NewCanRegisterRule creates a CanRegisterRule. Pass a non-empty email to also
// validate the email allowlist; pass "" to check only registration-enabled.
func NewCanRegisterRule(a app.AppInterface, email string) *CanRegisterRule {
	r := &CanRegisterRule{}

	r.Rule.SetContext(canRegisterContext{
		registrationEnabled: a.GetConfig().GetRegistrationEnabled(),
		allowedEmails:       a.GetConfig().GetEmailsAllowedAccess(),
		email:               email,
	})

	r.Rule.SetCondition(func(ctx any) bool {
		data := ctx.(canRegisterContext)

		if !data.registrationEnabled {
			r.AddFailMessage("Registrations are currently disabled")
			return false
		}

		if data.email != "" && len(data.allowedEmails) > 0 {
			found := false
			for _, allowed := range data.allowedEmails {
				if allowed == data.email {
					found = true
					break
				}
			}
			if !found {
				r.AddFailMessage("Your email is not permitted to register")
				return false
			}
		}

		return true
	})

	return r
}
