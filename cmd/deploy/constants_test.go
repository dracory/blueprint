package main

import (
	"testing"
)

// TestConstants tests that constants are defined
func TestConstants_SSH_KEY(t *testing.T) {
	t.Parallel()
	if SSH_KEY == "" {
		t.Error("SSH_KEY should not be empty")
	}
}

func TestConstants_SSH_USER(t *testing.T) {
	t.Parallel()
	if SSH_USER == "" {
		t.Error("SSH_USER should not be empty")
	}
}

func TestConstants_SSH_HOST(t *testing.T) {
	t.Parallel()
	if SSH_HOST == "" {
		t.Error("SSH_HOST should not be empty")
	}
}

func TestConstants_REMOTE_APP_DIR(t *testing.T) {
	t.Parallel()
	if REMOTE_APP_DIR == "" {
		t.Error("REMOTE_APP_DIR should not be empty")
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
