package authrules

import (
	"project/internal/testutils"
	"testing"
)

func TestEmailAllowedRule_EmptyAllowlist_Passes(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetEmailsAllowedAccess([]string{})
	app := testutils.Setup(testutils.WithCfg(cfg))

	r := NewEmailAllowedRule(app, "anyone@example.com")
	if r.Fails() {
		t.Errorf("expected rule to pass with empty allowlist, got: %s", r.FailMessageFirst())
	}
}

func TestEmailAllowedRule_EmailInAllowlist_Passes(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetEmailsAllowedAccess([]string{"allowed@example.com", "also@example.com"})
	app := testutils.Setup(testutils.WithCfg(cfg))

	r := NewEmailAllowedRule(app, "allowed@example.com")
	if r.Fails() {
		t.Errorf("expected rule to pass for allowed email, got: %s", r.FailMessageFirst())
	}
}

func TestEmailAllowedRule_EmailNotInAllowlist_Fails(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetEmailsAllowedAccess([]string{"allowed@example.com"})
	app := testutils.Setup(testutils.WithCfg(cfg))

	r := NewEmailAllowedRule(app, "stranger@example.com")
	if r.Passes() {
		t.Error("expected rule to fail for email not in allowlist")
	}
	if r.FailMessageFirst() == "" {
		t.Error("expected a fail message to be set")
	}
}
