package main

import (
	"testing"
)

// TestConstants tests that constants are defined
func TestConstants(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		constant string
	}{
		{"SSH_KEY", SSH_KEY},
		{"SSH_USER", SSH_USER},
		{"SSH_HOST", SSH_HOST},
		{"REMOTE_APP_DIR", REMOTE_APP_DIR},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant == "" {
				t.Errorf("%s should not be empty", tt.name)
			}
		})
	}
}

// TestPM2ProcessName tests PM2 process name constant
func TestPM2ProcessName(t *testing.T) {
	t.Parallel()
	// PM2_PROCESS_NAME is defined as REMOTE_APP_DIR, so they should have the same value
	if PM2_PROCESS_NAME != REMOTE_APP_DIR {
		t.Errorf("PM2_PROCESS_NAME = %q, REMOTE_APP_DIR = %q, should be equal", PM2_PROCESS_NAME, REMOTE_APP_DIR)
	}
}

// TestOtherFilesToDeploy tests OTHER_FILES_TO_DEPLOY variable
func TestOtherFilesToDeploy(t *testing.T) {
	t.Parallel()
	// Should be defined (can be empty slice)
	if OTHER_FILES_TO_DEPLOY == nil {
		t.Error("OTHER_FILES_TO_DEPLOY should not be nil")
	}
}
