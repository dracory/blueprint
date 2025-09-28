package config

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestDeriveEnvEncKey(t *testing.T) {
	originalPublic := envencKeyPublic
	originalDeobfuscate := deobfuscate
	t.Cleanup(func() {
		envencKeyPublic = originalPublic
		deobfuscate = originalDeobfuscate
	})

	validPublic := "PUBLIC_KEY_WITH_LENGTH_32_CHARS_"
	if len(validPublic) < 32 {
		t.Fatalf("test public key must be at least 32 characters, got %d", len(validPublic))
	}
	validPrivate := "PRIVATE_KEY_WITH_LENGTH_32_CHARS"
	if len(validPrivate) < 32 {
		t.Fatalf("test private key must be at least 32 characters, got %d", len(validPrivate))
	}

	type deobfuscateFunc func(string) (string, error)

	tests := []struct {
		name        string
		pub         string
		private     string
		deobfuscate deobfuscateFunc
		want        string
		wantErr     string
	}{
		{
			name:        "success returns derived key",
			pub:         validPublic,
			private:     validPrivate,
			deobfuscate: func(in string) (string, error) { return in, nil },
			want: func() string {
				sum := sha256.Sum256([]byte(validPublic + validPrivate))
				return fmt.Sprintf("%x", sum)
			}(),
		},
		{
			name:    "error when public key empty",
			pub:     "",
			private: validPrivate,
			wantErr: "envenc public key is empty",
		},
		{
			name:    "error when private key empty",
			pub:     validPublic,
			private: "",
			wantErr: "envenc private key is empty",
		},
		{
			name:    "error when public key short",
			pub:     "short",
			private: validPrivate,
			wantErr: "envenc public key is too short",
		},
		{
			name:    "error when private key short",
			pub:     validPublic,
			private: "short",
			wantErr: "envenc private key is too short",
		},
		{
			name:    "error deobfuscation failure",
			pub:     validPublic,
			private: validPrivate,
			deobfuscate: func(string) (string, error) {
				return "", errors.New("boom")
			},
			wantErr: "failed to deobfuscate public key",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			envencKeyPublic = tc.pub
			if tc.deobfuscate != nil {
				deobfuscate = tc.deobfuscate
			} else {
				deobfuscate = originalDeobfuscate
			}

			got, err := deriveEnvEncKey(tc.private)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error to contain %q, got %v", tc.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("expected key %q, got %q", tc.want, got)
			}
		})
	}
}
