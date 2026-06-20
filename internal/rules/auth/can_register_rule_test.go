package authrules

import (
	"project/internal/testutils"
	"testing"
)

func TestCanRegisterRule_RegistrationDisabled_Fails(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetRegistrationEnabled(false)
	cfg.SetEmailsAllowedAccess([]string{})
	app := testutils.Setup(testutils.WithCfg(cfg))

	r := NewCanRegisterRule(app, "user@example.com")
	if r.Passes() {
		t.Error("expected rule to fail when registration is disabled")
	}
	if r.FailMessageFirst() != "Registrations are currently disabled" {
		t.Errorf("expected message 'Registrations are currently disabled', got: %s", r.FailMessageFirst())
	}
}

func TestCanRegisterRule_RegistrationEnabled_NoAllowlist_Passes(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetRegistrationEnabled(true)
	cfg.SetEmailsAllowedAccess([]string{})
	app := testutils.Setup(testutils.WithCfg(cfg))

	r := NewCanRegisterRule(app, "anyone@example.com")
	if r.Fails() {
		t.Errorf("expected rule to pass with open access, got: %s", r.FailMessageFirst())
	}
}

func TestCanRegisterRule_RegistrationEnabled_EmailInAllowlist_Passes(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetRegistrationEnabled(true)
	cfg.SetEmailsAllowedAccess([]string{"permitted@example.com"})
	app := testutils.Setup(testutils.WithCfg(cfg))

	r := NewCanRegisterRule(app, "permitted@example.com")
	if r.Fails() {
		t.Errorf("expected rule to pass for permitted email, got: %s", r.FailMessageFirst())
	}
}

func TestCanRegisterRule_RegistrationEnabled_EmailNotInAllowlist_Fails(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetRegistrationEnabled(true)
	cfg.SetEmailsAllowedAccess([]string{"permitted@example.com"})
	app := testutils.Setup(testutils.WithCfg(cfg))

	r := NewCanRegisterRule(app, "stranger@example.com")
	if r.Passes() {
		t.Error("expected rule to fail for email not in allowlist")
	}
}

func TestCanRegisterRule_NoEmail_ChecksOnlyRegistrationEnabled(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetRegistrationEnabled(true)
	cfg.SetEmailsAllowedAccess([]string{"permitted@example.com"})
	app := testutils.Setup(testutils.WithCfg(cfg))

	r := NewCanRegisterRule(app, "")
	if r.Fails() {
		t.Errorf("expected rule to pass when email is empty (skip allowlist check), got: %s", r.FailMessageFirst())
	}
}
