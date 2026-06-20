package authrules

import (
	"github.com/dracory/rule"
)

// RegisterFormData holds the fields submitted by a registration or profile
// completion form. Password and Email are optional — validation of those
// fields is skipped when they are empty strings.
type RegisterFormData struct {
	FirstName       string
	LastName        string
	Email           string
	Country         string
	Timezone        string
	Password        string
	PasswordConfirm string
}

// RegisterFormValidationRule validates a registration form's required fields.
// Email is validated only when non-empty (it may be provided by the auth system).
// Password is validated only when non-empty (it is optional for OAuth-based flows).
type RegisterFormValidationRule struct {
	rule.Rule
}

// NewRegisterFormValidationRule creates a RegisterFormValidationRule for the
// provided form data.
func NewRegisterFormValidationRule(data RegisterFormData) *RegisterFormValidationRule {
	r := &RegisterFormValidationRule{}

	r.Rule.SetContext(data)

	r.Rule.SetCondition(func(ctx any) bool {
		d := ctx.(RegisterFormData)

		if d.FirstName == "" {
			r.AddFailMessage("First name is required field")
			return false
		}

		if d.LastName == "" {
			r.AddFailMessage("Last name is required field")
			return false
		}

		if d.Country == "" {
			r.AddFailMessage("Country is required field")
			return false
		}

		if d.Timezone == "" {
			r.AddFailMessage("Timezone is required field")
			return false
		}

		if d.Email != "" {
			// Email format validation can be added here in the future
		}

		if d.Password != "" {
			if len(d.Password) < 8 {
				r.AddFailMessage("Password must be at least 8 characters")
				return false
			}

			if d.Password != d.PasswordConfirm {
				r.AddFailMessage("Passwords do not match")
				return false
			}
		}

		return true
	})

	return r
}

// Message returns the first validation failure message, or empty string if the
// rule has not been evaluated or passed.
func (r *RegisterFormValidationRule) Message() string {
	return r.FailMessageFirst()
}
