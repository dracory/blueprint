package authrules

import (
	"project/internal/testutils"
	"testing"
)

func TestCanUsePasswordAuthRule_Enabled_Passes(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetPasswordAuthEnabled(true)
	app := testutils.Setup(testutils.WithCfg(cfg))

	r := NewCanUsePasswordAuthRule(app)
	if r.Fails() {
		t.Errorf("expected rule to pass when password auth is enabled, got: %s", r.FailMessageFirst())
	}
}

func TestCanUsePasswordAuthRule_Disabled_Fails(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetPasswordAuthEnabled(false)
	app := testutils.Setup(testutils.WithCfg(cfg))

	r := NewCanUsePasswordAuthRule(app)
	if r.Passes() {
		t.Error("expected rule to fail when password auth is disabled")
	}
	if r.FailMessageFirst() == "" {
		t.Error("expected a fail message to be set")
	}
}
