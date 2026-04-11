package main

import (
	"os/user"
	"strings"
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

// TestValidateCommandAllAllowedCommands tests all allowed commands
func TestValidateCommandAllAllowedCommands(t *testing.T) {
	allowedCommands := []string{"ls", "pwd", "cat", "grep", "find", "ps", "df", "du", "whoami", "id", "date", "uptime", "top", "free", "uname"}

	for _, cmd := range allowedCommands {
		t.Run(cmd, func(t *testing.T) {
			err := validateCommand(cmd)
			if err != nil {
				t.Errorf("Command '%s' should be allowed but got error: %v", cmd, err)
			}
		})
	}
}

// TestValidateCommandWithMultipleArguments tests commands with multiple arguments
func TestValidateCommandWithMultipleArguments(t *testing.T) {
	tests := []struct {
		name        string
		cmd         string
		expectError bool
	}{
		{"grep with multiple args", "grep -r pattern /path/to/search", false},
		{"find with multiple args", "find /path -name pattern -type f", false},
		{"ls with multiple flags", "ls -la -h /path", false},
		{"cat with multiple files", "cat file1.txt file2.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommand(tt.cmd)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for command: %s", tt.cmd)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for command '%s': %v", tt.cmd, err)
			}
		})
	}
}

// TestPrivateKeyPath tests the private key path generation
func TestPrivateKeyPath(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("Failed to get current user: %v", err)
	}

	tests := []struct {
		name     string
		sshKey   string
		expected string
	}{
		{"id_rsa", "id_rsa", currentUser.HomeDir + "/.ssh/id_rsa"},
		{"id_ed25519", "id_ed25519", currentUser.HomeDir + "/.ssh/id_ed25519"},
		{"custom_key", "custom_key", currentUser.HomeDir + "/.ssh/custom_key"},
		{"key_with_path", "keys/deploy_key", currentUser.HomeDir + "/.ssh/keys/deploy_key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PrivateKeyPath(tt.sshKey)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestPrivateKeyPathFormat tests that the path has correct format
func TestPrivateKeyPathFormat(t *testing.T) {
	result := PrivateKeyPath("test_key")

	// Should contain .ssh directory
	if !strings.Contains(result, ".ssh") {
		t.Errorf("Path should contain .ssh directory: %s", result)
	}

	// Should end with the key name
	if !strings.HasSuffix(result, "test_key") {
		t.Errorf("Path should end with key name: %s", result)
	}

	// Should be absolute path (start with /)
	if !strings.HasPrefix(result, "/") && !strings.HasPrefix(result, "C:") {
		t.Errorf("Path should be absolute: %s", result)
	}
}

// TestInitConfig tests configuration initialization
func TestInitConfig(t *testing.T) {
	config := InitConfig()

	// Verify all required fields are populated
	if config.Timestamp == "" {
		t.Error("Timestamp should not be empty")
	}

	if config.BuildLocalExecutableTempPath == "" {
		t.Error("BuildLocalExecutableTempPath should not be empty")
	}

	if config.SSHKey == "" {
		t.Error("SSHKey should not be empty")
	}

	if config.SSHUser == "" {
		t.Error("SSHUser should not be empty")
	}

	if config.SSHHost == "" {
		t.Error("SSHHost should not be empty")
	}

	if config.SSHLogin == "" {
		t.Error("SSHLogin should not be empty")
	}

	if config.RemoteAppDir == "" {
		t.Error("RemoteAppDir should not be empty")
	}

	if config.RemoteDeployDir == "" {
		t.Error("RemoteDeployDir should not be empty")
	}

	if config.RemoteTempDeployName == "" {
		t.Error("RemoteTempDeployName should not be empty")
	}

	if config.PM2ProcessName == "" {
		t.Error("PM2ProcessName should not be empty")
	}
}

// TestInitConfigSSHLogin tests that SSHLogin is correctly formatted
func TestInitConfigSSHLogin(t *testing.T) {
	config := InitConfig()

	// SSHLogin should be in format "user@host"
	parts := strings.Split(config.SSHLogin, "@")
	if len(parts) != 2 {
		t.Errorf("SSHLogin should be in format 'user@host', got: %s", config.SSHLogin)
	}

	if parts[0] != config.SSHUser {
		t.Errorf("SSHLogin user part should match SSHUser, got: %s vs %s", parts[0], config.SSHUser)
	}

	if parts[1] != config.SSHHost {
		t.Errorf("SSHLogin host part should match SSHHost, got: %s vs %s", parts[1], config.SSHHost)
	}
}

// TestInitConfigRemoteDeployDir tests that RemoteDeployDir is correctly formatted
func TestInitConfigRemoteDeployDir(t *testing.T) {
	config := InitConfig()

	// RemoteDeployDir should start with /home/
	if !strings.HasPrefix(config.RemoteDeployDir, "/home/") {
		t.Errorf("RemoteDeployDir should start with /home/, got: %s", config.RemoteDeployDir)
	}

	// RemoteDeployDir should contain the SSH user
	if !strings.Contains(config.RemoteDeployDir, config.SSHUser) {
		t.Errorf("RemoteDeployDir should contain SSHUser, got: %s", config.RemoteDeployDir)
	}
}

// TestInitConfigTimestampFormat tests that timestamp has correct format
func TestInitConfigTimestampFormat(t *testing.T) {
	config := InitConfig()

	// Timestamp should be in format Ymd_His (e.g., 20260411_120000)
	// Should be 15 characters long (8 for date + 1 for underscore + 6 for time)
	if len(config.Timestamp) != 15 {
		t.Errorf("Timestamp should be 15 characters long, got: %d (%s)", len(config.Timestamp), config.Timestamp)
	}

	// Should contain underscore separator
	if !strings.Contains(config.Timestamp, "_") {
		t.Errorf("Timestamp should contain underscore separator, got: %s", config.Timestamp)
	}
}

// TestGetDeployCommands tests that deploy commands are generated correctly
func TestGetDeployCommands(t *testing.T) {
	config := InitConfig()
	commands := GetDeployCommands(config)

	// Should have at least some commands
	if len(commands) == 0 {
		t.Error("GetDeployCommands should return at least one command")
	}

	// Verify all commands have required fields
	for i, cmd := range commands {
		if cmd.Reason == "" {
			t.Errorf("Command %d should have a reason", i)
		}

		if cmd.Cmd == "" {
			t.Errorf("Command %d should have a command", i)
		}
	}
}

// TestGetDeployCommandsRequiredFields tests that required commands are marked correctly
func TestGetDeployCommandsRequiredFields(t *testing.T) {
	config := InitConfig()
	commands := GetDeployCommands(config)

	// Count required and optional commands
	requiredCount := 0
	optionalCount := 0

	for _, cmd := range commands {
		if cmd.Required {
			requiredCount++
		} else {
			optionalCount++
		}
	}

	// Should have both required and optional commands
	if requiredCount == 0 {
		t.Error("Should have at least one required command")
	}

	if optionalCount == 0 {
		t.Error("Should have at least one optional command")
	}
}

// TestGetDeployCommandsContainsTimestamp tests that commands contain the timestamp
func TestGetDeployCommandsContainsTimestamp(t *testing.T) {
	config := InitConfig()
	commands := GetDeployCommands(config)

	// At least some commands should contain the timestamp for backup naming
	foundTimestamp := false
	for _, cmd := range commands {
		if strings.Contains(cmd.Cmd, config.Timestamp) {
			foundTimestamp = true
			break
		}
	}

	if !foundTimestamp {
		t.Error("At least one command should contain the timestamp for backup naming")
	}
}

// TestGetDeployCommandsContainsConfigValues tests that commands use config values
func TestGetDeployCommandsContainsConfigValues(t *testing.T) {
	config := InitConfig()
	commands := GetDeployCommands(config)

	commandsStr := ""
	for _, cmd := range commands {
		commandsStr += cmd.Cmd + " "
	}

	// Commands should reference the remote deploy directory
	if !strings.Contains(commandsStr, config.RemoteDeployDir) {
		t.Error("Commands should reference RemoteDeployDir")
	}

	// Commands should reference the PM2 process name
	if !strings.Contains(commandsStr, config.PM2ProcessName) {
		t.Error("Commands should reference PM2ProcessName")
	}
}

// TestValidateCommandWithDangerousCharacters tests all dangerous characters are blocked
func TestValidateCommandWithDangerousCharacters(t *testing.T) {
	dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "<", ">", "\"", "'"}

	for _, char := range dangerousChars {
		t.Run("dangerous_char_"+char, func(t *testing.T) {
			cmd := "ls " + char + " rm -rf /"
			err := validateCommand(cmd)
			if err == nil {
				t.Errorf("Command with '%s' should be rejected", char)
			}
			if !strings.Contains(err.Error(), "dangerous character") {
				t.Errorf("Error should mention dangerous character, got: %v", err)
			}
		})
	}
}

// TestValidateCommandEdgeCases tests edge cases in command validation
func TestValidateCommandEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		cmd         string
		expectError bool
	}{
		{"Command with only spaces", "   ", true},
		{"Command with tabs", "ls\t-la", false},
		{"Very long command", "ls " + strings.Repeat("arg ", 100), false},
		{"Command with numbers", "ls 123", false},
		{"Command with hyphens", "ls -la -h", false},
		{"Command with equals", "ls --color=auto", false},
		{"Command with slashes", "ls /home/user", false},
		{"Command with dots", "ls ../files", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommand(tt.cmd)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for: %s", tt.cmd)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for '%s': %v", tt.cmd, err)
			}
		})
	}
}

// TestInitConfigConsistency tests that InitConfig produces consistent results
func TestInitConfigConsistency(t *testing.T) {
	config1 := InitConfig()
	config2 := InitConfig()

	// SSHKey, SSHUser, SSHHost should be the same (constants)
	if config1.SSHKey != config2.SSHKey {
		t.Error("SSHKey should be consistent across calls")
	}

	if config1.SSHUser != config2.SSHUser {
		t.Error("SSHUser should be consistent across calls")
	}

	if config1.SSHHost != config2.SSHHost {
		t.Error("SSHHost should be consistent across calls")
	}

	// Timestamps will be different (they're generated at call time)
	// but both should be valid
	if config1.Timestamp == "" || config2.Timestamp == "" {
		t.Error("Both timestamps should be non-empty")
	}
}

// TestInitConfigPathConstruction tests that paths are correctly constructed
func TestInitConfigPathConstruction(t *testing.T) {
	config := InitConfig()

	// BuildLocalExecutableTempPath should contain timestamp
	if !strings.Contains(config.BuildLocalExecutableTempPath, config.Timestamp) {
		t.Errorf("BuildLocalExecutableTempPath should contain timestamp: %s", config.BuildLocalExecutableTempPath)
	}

	// RemoteTempDeployName should contain timestamp
	if !strings.Contains(config.RemoteTempDeployName, config.Timestamp) {
		t.Errorf("RemoteTempDeployName should contain timestamp: %s", config.RemoteTempDeployName)
	}

	// RemoteDeployDir should be constructed from SSHUser and RemoteAppDir
	expectedRemoteDeployDir := "/home/" + config.SSHUser + "/" + config.RemoteAppDir
	if config.RemoteDeployDir != expectedRemoteDeployDir {
		t.Errorf("RemoteDeployDir mismatch: expected %s, got %s", expectedRemoteDeployDir, config.RemoteDeployDir)
	}
}

// TestGetDeployCommandsCount tests the number of deploy commands
func TestGetDeployCommandsCount(t *testing.T) {
	config := InitConfig()
	commands := GetDeployCommands(config)

	// Should have a reasonable number of commands (at least 5, at most 20)
	if len(commands) < 5 {
		t.Errorf("Should have at least 5 deploy commands, got %d", len(commands))
	}

	if len(commands) > 20 {
		t.Errorf("Should have at most 20 deploy commands, got %d", len(commands))
	}
}

// TestGetDeployCommandsOrder tests that critical commands are in correct order
func TestGetDeployCommandsOrder(t *testing.T) {
	config := InitConfig()
	commands := GetDeployCommands(config)

	// Find indices of key commands
	var chmodIdx, mvIdx, pmStartIdx int
	foundChmod := false
	foundMv := false
	foundPmStart := false

	for i, cmd := range commands {
		if strings.Contains(cmd.Cmd, "chmod") && !foundChmod {
			chmodIdx = i
			foundChmod = true
		}
		if strings.Contains(cmd.Cmd, "mv") && !foundMv {
			mvIdx = i
			foundMv = true
		}
		if strings.Contains(cmd.Cmd, "pm2 start") && !foundPmStart {
			pmStartIdx = i
			foundPmStart = true
		}
	}

	// chmod should come before mv
	if foundChmod && foundMv && chmodIdx >= mvIdx {
		t.Error("chmod command should come before mv command")
	}

	// pm2 start should be last
	if foundPmStart && pmStartIdx != len(commands)-1 {
		t.Error("pm2 start command should be the last command")
	}
}
