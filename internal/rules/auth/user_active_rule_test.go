package authrules

import (
	"testing"

	"github.com/dracory/userstore"
)

func TestUserActiveRule_ActiveUser_Passes(t *testing.T) {
	user := userstore.NewUser().SetStatus(userstore.USER_STATUS_ACTIVE)

	r := NewUserActiveRule(user)
	if r.Fails() {
		t.Errorf("expected rule to pass for active user, got: %s", r.FailMessageFirst())
	}
}

func TestUserActiveRule_InactiveUser_Fails(t *testing.T) {
	user := userstore.NewUser().SetStatus(userstore.USER_STATUS_INACTIVE)

	r := NewUserActiveRule(user)
	if r.Passes() {
		t.Error("expected rule to fail for inactive user")
	}
	if r.FailMessageFirst() == "" {
		t.Error("expected a fail message to be set")
	}
}

func TestUserActiveRule_NilUser_Fails(t *testing.T) {
	r := NewUserActiveRule(nil)
	if r.Passes() {
		t.Error("expected rule to fail for nil user")
	}
	if r.FailMessageFirst() == "" {
		t.Error("expected a fail message to be set")
	}
}
