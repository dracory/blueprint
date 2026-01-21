package main

import (
	"testing"
)

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name        string
		cmd         string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Allowed command - ls",
			cmd:         "ls",
			expectError: false,
		},
		{
			name:        "Allowed command with arguments - ls -la",
			cmd:         "ls -la",
			expectError: false,
		},
		{
			name:        "Allowed command - pwd",
			cmd:         "pwd",
			expectError: false,
		},
		{
			name:        "Allowed command - cat",
			cmd:         "cat file.txt",
			expectError: false,
		},
		{
			name:        "Disallowed command - rm",
			cmd:         "rm file.txt",
			expectError: true,
			errorMsg:    "command not allowed: rm",
		},
		{
			name:        "Disallowed command - sudo",
			cmd:         "sudo ls",
			expectError: true,
			errorMsg:    "command not allowed: sudo",
		},
		{
			name:        "Empty command",
			cmd:         "",
			expectError: true,
			errorMsg:    "empty command not allowed",
		},
		{
			name:        "Command with semicolon injection",
			cmd:         "ls; rm -rf /",
			expectError: true,
			errorMsg:    "dangerous character detected in command: ;",
		},
		{
			name:        "Command with ampersand injection",
			cmd:         "ls & rm -rf /",
			expectError: true,
			errorMsg:    "dangerous character detected in command: &",
		},
		{
			name:        "Command with pipe injection",
			cmd:         "ls | rm -rf /",
			expectError: true,
			errorMsg:    "dangerous character detected in command: |",
		},
		{
			name:        "Command with backtick injection",
			cmd:         "ls `rm -rf /`",
			expectError: true,
			errorMsg:    "dangerous character detected in command: `",
		},
		{
			name:        "Command with dollar sign injection",
			cmd:         "ls $HOME",
			expectError: true,
			errorMsg:    "dangerous character detected in command: $",
		},
		{
			name:        "Command with parentheses injection",
			cmd:         "ls $(rm -rf /)",
			expectError: true,
			errorMsg:    "dangerous character detected in command: $",
		},
		{
			name:        "Command with redirection injection",
			cmd:         "ls > /etc/passwd",
			expectError: true,
			errorMsg:    "dangerous character detected in command: >",
		},
		{
			name:        "Command with quote injection",
			cmd:         "ls \"rm -rf /\"",
			expectError: true,
			errorMsg:    "dangerous character detected in command: \"",
		},
		{
			name:        "Command with single quote injection",
			cmd:         "ls 'rm -rf /'",
			expectError: true,
			errorMsg:    "dangerous character detected in command: '",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommand(tt.cmd)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s' but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}
