package authrules

import (
	"strconv"

	"github.com/dracory/rule"
)

const passwordPolicyMinLength = 8

// PasswordPolicyData holds the password fields for policy validation
type PasswordPolicyData struct {
	Password        string
	PasswordConfirm string
}

// PasswordPolicyRule validates password min length and that the
// confirmation matches. It is designed to be composed by other rules
// (password registration, password reset, etc.)
type PasswordPolicyRule struct {
	rule.Rule
}

// NewPasswordPolicyRule creates a new PasswordPolicyRule
func NewPasswordPolicyRule(data PasswordPolicyData) *PasswordPolicyRule {
	r := &PasswordPolicyRule{}

	r.Rule.SetContext(passwordPolicyContext{
		password:        data.Password,
		passwordConfirm: data.PasswordConfirm,
	})

	r.Rule.SetCondition(func(ctx any) bool {
		c := ctx.(passwordPolicyContext)

		// Check if password is empty
		if c.password == "" {
			r.AddFailMessage("Password is required")
			return false
		}

		// Check if password length is valid
		if len(c.password) < passwordPolicyMinLength {
			r.AddFailMessage("Password must be at least " + strconv.Itoa(passwordPolicyMinLength) + " characters")
			return false
		}

		// Check if passwords match
		if c.password != c.passwordConfirm {
			r.AddFailMessage("Passwords do not match")
			return false
		}
		return true
	})

	return r
}

type passwordPolicyContext struct {
	password        string
	passwordConfirm string
}
