package authrules

import "testing"

func TestPasswordPolicyRule_Passes(t *testing.T) {
	r := NewPasswordPolicyRule(PasswordPolicyData{
		Password:        "password123",
		PasswordConfirm: "password123",
	})
	if r.Fails() {
		t.Fatalf("expected rule to pass, got: %s", r.FailMessageFirst())
	}
}

func TestPasswordPolicyRule_Fails(t *testing.T) {
	tests := []struct {
		name            string
		password        string
		passwordConfirm string
		wantMessage     string
	}{
		{"empty password", "", "", "Password is required"},
		{"too short", "short", "short", "Password must be at least 8 characters"},
		{"mismatch", "password123", "password124", "Passwords do not match"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewPasswordPolicyRule(PasswordPolicyData{
				Password:        tc.password,
				PasswordConfirm: tc.passwordConfirm,
			})
			if !r.Fails() {
				t.Fatal("expected rule to fail")
			}
			if r.FailMessageFirst() != tc.wantMessage {
				t.Fatalf("expected message %q, got %q", tc.wantMessage, r.FailMessageFirst())
			}
		})
	}
}
