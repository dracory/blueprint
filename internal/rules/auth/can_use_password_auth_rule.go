package authrules

import (
	"project/internal/app"

	"github.com/dracory/rule"
)

// CanUsePasswordAuthRule checks whether password-based authentication is
// enabled in the application configuration.
type CanUsePasswordAuthRule struct {
	rule.Rule
}

// NewCanUsePasswordAuthRule creates a CanUsePasswordAuthRule using the
// application configuration to determine if password auth is enabled.
func NewCanUsePasswordAuthRule(a app.AppInterface) *CanUsePasswordAuthRule {
	r := &CanUsePasswordAuthRule{}

	r.Rule.SetContext(a.GetConfig().GetPasswordAuthEnabled())

	r.Rule.SetCondition(func(ctx any) bool {
		enabled, ok := ctx.(bool)
		if !ok || !enabled {
			r.AddFailMessage("Password authentication is not enabled")
			return false
		}
		return true
	})

	return r
}
