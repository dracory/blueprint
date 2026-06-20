package authrules

import (
	"github.com/dracory/rule"
	"github.com/dracory/userstore"
)

// UserActiveRule checks whether a user account is in active status.
type UserActiveRule struct {
	rule.Rule
}

// NewUserActiveRule creates a UserActiveRule for the given user.
// A nil user is treated as inactive.
func NewUserActiveRule(user userstore.UserInterface) *UserActiveRule {
	r := &UserActiveRule{}

	r.Rule.SetContext(user)

	r.Rule.SetCondition(func(ctx any) bool {
		u, ok := ctx.(userstore.UserInterface)
		if !ok || u == nil {
			r.AddFailMessage("User account not found")
			return false
		}

		if u.GetStatus() != userstore.USER_STATUS_ACTIVE {
			r.AddFailMessage("User account is not active. Please contact our support team for assistance.")
			return false
		}

		return true
	})

	return r
}
