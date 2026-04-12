package main

import (
	"strings"
	"testing"
)

// TestConstants tests that constants are defined
func TestConstants(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		constant string
		contains string
	}{
		{"SSH_KEY", SSH_KEY, "{{ SSHKEY }}"},
		{"SSH_USER", SSH_USER, "{{ SSHUSER }}"},
		{"SSH_HOST", SSH_HOST, "{{ SSHHOST }}"},
		{"REMOTE_APP_DIR", REMOTE_APP_DIR, "{{ APP_NAME }}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant == "" {
				t.Errorf("%s should not be empty", tt.name)
			}
			if !strings.Contains(tt.constant, tt.contains) {
				t.Errorf("%s = %q, should contain %q", tt.name, tt.constant, tt.contains)
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
