package constants

import "testing"

func TestTaskAliasesAreUniqueAndNonEmpty(t *testing.T) {
	aliases := map[string]string{
		"BlindIndexRebuildTaskAlias":                     BlindIndexRebuildTaskAlias,
		"CleanUpTaskAlias":                               CleanUpTaskAlias,
		"EmailTestTaskAlias":                             EmailTestTaskAlias,
		"EmailToAdminTaskAlias":                          EmailToAdminTaskAlias,
		"EmailAuthDeliveryTaskAlias":                     EmailAuthDeliveryTaskAlias,
		"EmailToAdminOnNewContactFormSubmittedTaskAlias": EmailToAdminOnNewContactFormSubmittedTaskAlias,
		"EmailToAdminOnNewUserRegisteredTaskAlias":       EmailToAdminOnNewUserRegisteredTaskAlias,
		"HelloWorldTaskAlias":                            HelloWorldTaskAlias,
		"StatsVisitorEnhanceTaskAlias":                   StatsVisitorEnhanceTaskAlias,
	}

	seen := map[string]string{}
	for name, alias := range aliases {
		if alias == "" {
			t.Errorf("%s is empty", name)
		}
		if prev, exists := seen[alias]; exists {
			t.Errorf("alias %q is shared by %s and %s", alias, prev, name)
		}
		seen[alias] = name
	}
}
