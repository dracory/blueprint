package config

import "testing"

func TestMissingEnvErrorError(t *testing.T) {
	tests := []struct {
		name    string
		err     MissingEnvError
		want   string
	}{
		{
			name: "with context",
			err:  MissingEnvError{Key: "SOME_KEY", Context: "provide value"},
			want: "config: required env \"SOME_KEY\" is missing: provide value",
		},
		{
			name: "without context",
			err:  MissingEnvError{Key: "OTHER_KEY"},
			want: "config: required env \"OTHER_KEY\" is missing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.err.Error(); got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestEnsureRequired(t *testing.T) {
	t.Run("returns nil when value present", func(t *testing.T) {
		if err := ensureRequired("value", "KEY", "ctx"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns error when empty", func(t *testing.T) {
		err := ensureRequired(" \t\n ", "KEY", "ctx")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if _, ok := err.(MissingEnvError); !ok {
			t.Fatalf("expected MissingEnvError, got %T", err)
		}
	})
}

func TestRequireWhen(t *testing.T) {
	t.Run("skips when condition false", func(t *testing.T) {
		if err := requireWhen(false, "KEY", "ctx", ""); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("validates when condition true", func(t *testing.T) {
		err := requireWhen(true, "KEY", "ctx", "")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
