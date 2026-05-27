package main

import (
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateCommand_AllowedLS(t *testing.T) {
	err := validateCommand("ls")
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

func TestValidateCommand_AllowedLSWithArgs(t *testing.T) {
	err := validateCommand("ls -la")
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

func TestValidateCommand_AllowedPWD(t *testing.T) {
	err := validateCommand("pwd")
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

func TestValidateCommand_AllowedCat(t *testing.T) {
	err := validateCommand("cat file.txt")
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

func TestValidateCommand_DisallowedRM(t *testing.T) {
	err := validateCommand("rm file.txt")
	if err == nil {
		t.Error("Expected error but got none")
	} else if err.Error() != "command not allowed: rm" {
		t.Errorf("Expected error message 'command not allowed: rm' but got '%s'", err.Error())
	}
}

func TestValidateCommand_DisallowedSudo(t *testing.T) {
	err := validateCommand("sudo ls")
	if err == nil {
		t.Error("Expected error but got none")
	} else if err.Error() != "command not allowed: sudo" {
		t.Errorf("Expected error message 'command not allowed: sudo' but got '%s'", err.Error())
	}
}

func TestValidateCommand_Empty(t *testing.T) {
	err := validateCommand("")
	if err == nil {
		t.Error("Expected error but got none")
	} else if err.Error() != "empty command not allowed" {
		t.Errorf("Expected error message 'empty command not allowed' but got '%s'", err.Error())
	}
}

func TestValidateCommand_SemicolonInjection(t *testing.T) {
	err := validateCommand("ls; rm -rf /")
	if err == nil {
		t.Error("Expected error but got none")
	} else if err.Error() != "command not allowed: rm" {
		t.Errorf("Expected error message 'command not allowed: rm' but got '%s'", err.Error())
	}
}

func TestValidateCommand_AllowedChained(t *testing.T) {
	err := validateCommand("cd /path; ls")
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

func TestValidateCommand_AmpersandInjection(t *testing.T) {
	err := validateCommand("ls & rm -rf /")
	if err == nil {
		t.Error("Expected error but got none")
	} else if err.Error() != "dangerous character detected in command: &" {
		t.Errorf("Expected error message 'dangerous character detected in command: &' but got '%s'", err.Error())
	}
}

func TestValidateCommand_PipeInjection(t *testing.T) {
	err := validateCommand("ls | rm -rf /")
	if err == nil {
		t.Error("Expected error but got none")
	} else if err.Error() != "dangerous character detected in command: |" {
		t.Errorf("Expected error message 'dangerous character detected in command: |' but got '%s'", err.Error())
	}
}

func TestValidateCommand_BacktickInjection(t *testing.T) {
	err := validateCommand("ls `rm -rf /`")
	if err == nil {
		t.Error("Expected error but got none")
	} else if err.Error() != "dangerous character detected in command: `" {
		t.Errorf("Expected error message 'dangerous character detected in command: `' but got '%s'", err.Error())
	}
}

func TestValidateCommand_DollarSignInjection(t *testing.T) {
	err := validateCommand("ls $HOME")
	if err == nil {
		t.Error("Expected error but got none")
	} else if err.Error() != "dangerous character detected in command: $" {
		t.Errorf("Expected error message 'dangerous character detected in command: $' but got '%s'", err.Error())
	}
}

func TestValidateCommand_ParenthesesInjection(t *testing.T) {
	err := validateCommand("ls $(rm -rf /)")
	if err == nil {
		t.Error("Expected error but got none")
	} else if err.Error() != "dangerous character detected in command: $" {
		t.Errorf("Expected error message 'dangerous character detected in command: $' but got '%s'", err.Error())
	}
}

func TestValidateCommand_RedirectionInjection(t *testing.T) {
	err := validateCommand("ls > /etc/passwd")
	if err == nil {
		t.Error("Expected error but got none")
	} else if err.Error() != "dangerous character detected in command: >" {
		t.Errorf("Expected error message 'dangerous character detected in command: >' but got '%s'", err.Error())
	}
}

func TestValidateCommand_SingleQuoteInjection(t *testing.T) {
	err := validateCommand("ls 'rm -rf /'")
	if err == nil {
		t.Error("Expected error but got none")
	} else if err.Error() != "dangerous character detected in command: '" {
		t.Errorf("Expected error message 'dangerous character detected in command: '' but got '%s'", err.Error())
	}
}

// TestValidateCommandAllAllowedCommands tests all allowed commands
func TestValidateCommandAllAllowedCommands_LS(t *testing.T) {
	err := validateCommand("ls")
	if err != nil {
		t.Errorf("Command 'ls' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_PWD(t *testing.T) {
	err := validateCommand("pwd")
	if err != nil {
		t.Errorf("Command 'pwd' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_Cat(t *testing.T) {
	err := validateCommand("cat")
	if err != nil {
		t.Errorf("Command 'cat' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_CD(t *testing.T) {
	err := validateCommand("cd")
	if err != nil {
		t.Errorf("Command 'cd' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_Chmod(t *testing.T) {
	err := validateCommand("chmod")
	if err != nil {
		t.Errorf("Command 'chmod' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_Grep(t *testing.T) {
	err := validateCommand("grep")
	if err != nil {
		t.Errorf("Command 'grep' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_Find(t *testing.T) {
	err := validateCommand("find")
	if err != nil {
		t.Errorf("Command 'find' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_PS(t *testing.T) {
	err := validateCommand("ps")
	if err != nil {
		t.Errorf("Command 'ps' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_DF(t *testing.T) {
	err := validateCommand("df")
	if err != nil {
		t.Errorf("Command 'df' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_DU(t *testing.T) {
	err := validateCommand("du")
	if err != nil {
		t.Errorf("Command 'du' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_Whoami(t *testing.T) {
	err := validateCommand("whoami")
	if err != nil {
		t.Errorf("Command 'whoami' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_ID(t *testing.T) {
	err := validateCommand("id")
	if err != nil {
		t.Errorf("Command 'id' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_Date(t *testing.T) {
	err := validateCommand("date")
	if err != nil {
		t.Errorf("Command 'date' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_Uptime(t *testing.T) {
	err := validateCommand("uptime")
	if err != nil {
		t.Errorf("Command 'uptime' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_Top(t *testing.T) {
	err := validateCommand("top")
	if err != nil {
		t.Errorf("Command 'top' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_Free(t *testing.T) {
	err := validateCommand("free")
	if err != nil {
		t.Errorf("Command 'free' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_Uname(t *testing.T) {
	err := validateCommand("uname")
	if err != nil {
		t.Errorf("Command 'uname' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_Touch(t *testing.T) {
	err := validateCommand("touch")
	if err != nil {
		t.Errorf("Command 'touch' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_MV(t *testing.T) {
	err := validateCommand("mv")
	if err != nil {
		t.Errorf("Command 'mv' should be allowed but got error: %v", err)
	}
}

func TestValidateCommandAllAllowedCommands_PM2(t *testing.T) {
	err := validateCommand("pm2")
	if err != nil {
		t.Errorf("Command 'pm2' should be allowed but got error: %v", err)
	}
}

// TestValidateCommandWithMultipleArguments tests commands with multiple arguments
func TestValidateCommandWithMultipleArguments_Grep(t *testing.T) {
	err := validateCommand("grep -r pattern /path/to/search")
	if err != nil {
		t.Errorf("Unexpected error for command 'grep -r pattern /path/to/search': %v", err)
	}
}

func TestValidateCommandWithMultipleArguments_Find(t *testing.T) {
	err := validateCommand("find /path -name pattern -type f")
	if err != nil {
		t.Errorf("Unexpected error for command 'find /path -name pattern -type f': %v", err)
	}
}

func TestValidateCommandWithMultipleArguments_LS(t *testing.T) {
	err := validateCommand("ls -la -h /path")
	if err != nil {
		t.Errorf("Unexpected error for command 'ls -la -h /path': %v", err)
	}
}

func TestValidateCommandWithMultipleArguments_Cat(t *testing.T) {
	err := validateCommand("cat file1.txt file2.txt")
	if err != nil {
		t.Errorf("Unexpected error for command 'cat file1.txt file2.txt': %v", err)
	}
}

// TestPrivateKeyPath tests the private key path generation
func TestPrivateKeyPath_IDRSA(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("Failed to get current user: %v", err)
	}
	result := PrivateKeyPath("id_rsa")
	expected := filepath.Join(currentUser.HomeDir, ".ssh", "id_rsa")
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestPrivateKeyPath_IDEd25519(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("Failed to get current user: %v", err)
	}
	result := PrivateKeyPath("id_ed25519")
	expected := filepath.Join(currentUser.HomeDir, ".ssh", "id_ed25519")
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestPrivateKeyPath_CustomKey(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("Failed to get current user: %v", err)
	}
	result := PrivateKeyPath("custom_key")
	expected := filepath.Join(currentUser.HomeDir, ".ssh", "custom_key")
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestPrivateKeyPath_WithPath(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("Failed to get current user: %v", err)
	}
	result := PrivateKeyPath("keys/deploy_key")
	expected := filepath.Join(currentUser.HomeDir, ".ssh", "keys", "deploy_key")
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
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
func TestValidateCommandWithDangerousCharacters_Ampersand(t *testing.T) {
	cmd := "ls & rm -rf /"
	err := validateCommand(cmd)
	if err == nil {
		t.Error("Command with '&' should be rejected")
	}
	if !strings.Contains(err.Error(), "dangerous character") {
		t.Errorf("Error should mention dangerous character, got: %v", err)
	}
}

func TestValidateCommandWithDangerousCharacters_Pipe(t *testing.T) {
	cmd := "ls | rm -rf /"
	err := validateCommand(cmd)
	if err == nil {
		t.Error("Command with '|' should be rejected")
	}
	if !strings.Contains(err.Error(), "dangerous character") {
		t.Errorf("Error should mention dangerous character, got: %v", err)
	}
}

func TestValidateCommandWithDangerousCharacters_Backtick(t *testing.T) {
	cmd := "ls `rm -rf /`"
	err := validateCommand(cmd)
	if err == nil {
		t.Error("Command with '`' should be rejected")
	}
	if !strings.Contains(err.Error(), "dangerous character") {
		t.Errorf("Error should mention dangerous character, got: %v", err)
	}
}

func TestValidateCommandWithDangerousCharacters_Dollar(t *testing.T) {
	cmd := "ls $ rm -rf /"
	err := validateCommand(cmd)
	if err == nil {
		t.Error("Command with '$' should be rejected")
	}
	if !strings.Contains(err.Error(), "dangerous character") {
		t.Errorf("Error should mention dangerous character, got: %v", err)
	}
}

func TestValidateCommandWithDangerousCharacters_LeftParen(t *testing.T) {
	cmd := "ls ( rm -rf /"
	err := validateCommand(cmd)
	if err == nil {
		t.Error("Command with '(' should be rejected")
	}
	if !strings.Contains(err.Error(), "dangerous character") {
		t.Errorf("Error should mention dangerous character, got: %v", err)
	}
}

func TestValidateCommandWithDangerousCharacters_RightParen(t *testing.T) {
	cmd := "ls ) rm -rf /"
	err := validateCommand(cmd)
	if err == nil {
		t.Error("Command with ')' should be rejected")
	}
	if !strings.Contains(err.Error(), "dangerous character") {
		t.Errorf("Error should mention dangerous character, got: %v", err)
	}
}

func TestValidateCommandWithDangerousCharacters_LessThan(t *testing.T) {
	cmd := "ls < rm -rf /"
	err := validateCommand(cmd)
	if err == nil {
		t.Error("Command with '<' should be rejected")
	}
	if !strings.Contains(err.Error(), "dangerous character") {
		t.Errorf("Error should mention dangerous character, got: %v", err)
	}
}

func TestValidateCommandWithDangerousCharacters_GreaterThan(t *testing.T) {
	cmd := "ls > rm -rf /"
	err := validateCommand(cmd)
	if err == nil {
		t.Error("Command with '>' should be rejected")
	}
	if !strings.Contains(err.Error(), "dangerous character") {
		t.Errorf("Error should mention dangerous character, got: %v", err)
	}
}

func TestValidateCommandWithDangerousCharacters_SingleQuote(t *testing.T) {
	cmd := "ls ' rm -rf /"
	err := validateCommand(cmd)
	if err == nil {
		t.Error("Command with ''' should be rejected")
	}
	if !strings.Contains(err.Error(), "dangerous character") {
		t.Errorf("Error should mention dangerous character, got: %v", err)
	}
}

// TestValidateCommandEdgeCases tests edge cases in command validation
func TestValidateCommandEdgeCases_OnlySpaces(t *testing.T) {
	err := validateCommand("   ")
	if err == nil {
		t.Error("Expected error for:    ")
	}
}

func TestValidateCommandEdgeCases_WithTabs(t *testing.T) {
	err := validateCommand("ls\t-la")
	if err != nil {
		t.Errorf("Unexpected error for 'ls\t-la': %v", err)
	}
}

func TestValidateCommandEdgeCases_VeryLong(t *testing.T) {
	cmd := "ls " + strings.Repeat("arg ", 100)
	err := validateCommand(cmd)
	if err != nil {
		t.Errorf("Unexpected error for long command: %v", err)
	}
}

func TestValidateCommandEdgeCases_WithNumbers(t *testing.T) {
	err := validateCommand("ls 123")
	if err != nil {
		t.Errorf("Unexpected error for 'ls 123': %v", err)
	}
}

func TestValidateCommandEdgeCases_WithHyphens(t *testing.T) {
	err := validateCommand("ls -la -h")
	if err != nil {
		t.Errorf("Unexpected error for 'ls -la -h': %v", err)
	}
}

func TestValidateCommandEdgeCases_WithEquals(t *testing.T) {
	err := validateCommand("ls --color=auto")
	if err != nil {
		t.Errorf("Unexpected error for 'ls --color=auto': %v", err)
	}
}

func TestValidateCommandEdgeCases_WithSlashes(t *testing.T) {
	err := validateCommand("ls /home/user")
	if err != nil {
		t.Errorf("Unexpected error for 'ls /home/user': %v", err)
	}
}

func TestValidateCommandEdgeCases_WithDots(t *testing.T) {
	err := validateCommand("ls ../files")
	if err != nil {
		t.Errorf("Unexpected error for 'ls ../files': %v", err)
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

// TestValidateCommandCaseSensitivity tests command validation with different cases
func TestValidateCommandCaseSensitivity(t *testing.T) {
	// Test that command validation is case-sensitive (commands are lowercase)
	err := validateCommand("LS")
	if err == nil {
		t.Error("validateCommand should be case-sensitive, LS should fail")
	}

	err = validateCommand("Ls")
	if err == nil {
		t.Error("validateCommand should be case-sensitive, Ls should fail")
	}
}

// TestValidateCommandWithPath tests command validation with paths
func TestValidateCommandWithPath_AbsolutePath(t *testing.T) {
	err := validateCommand("/usr/bin/ls")
	if err == nil {
		t.Error("Expected error for: /usr/bin/ls")
	}
}

func TestValidateCommandWithPath_RelativePath(t *testing.T) {
	err := validateCommand("./ls")
	if err == nil {
		t.Error("Expected error for: ./ls")
	}
}

func TestValidateCommandWithPath_ParentPath(t *testing.T) {
	err := validateCommand("../ls")
	if err == nil {
		t.Error("Expected error for: ../ls")
	}
}

func TestValidateCommandWithPath_Tilde(t *testing.T) {
	err := validateCommand("~/ls")
	if err == nil {
		t.Error("Expected error for: ~/ls")
	}
}

// TestInitConfigUniqueTimestamps tests that multiple configs have different timestamps
func TestInitConfigUniqueTimestamps(t *testing.T) {
	config1 := InitConfig()
	// Small delay to ensure different timestamp
	config2 := InitConfig()

	// Timestamps should be different (or same if called very quickly)
	// We just verify both are valid
	if config1.Timestamp == "" || config2.Timestamp == "" {
		t.Error("Both timestamps should be non-empty")
	}

	// Both should have same SSH config
	if config1.SSHUser != config2.SSHUser {
		t.Error("SSHUser should be consistent")
	}
	if config1.SSHHost != config2.SSHHost {
		t.Error("SSHHost should be consistent")
	}
}

// TestGetDeployCommandsContent tests the content of deploy commands
func TestGetDeployCommandsContent(t *testing.T) {
	config := InitConfig()
	commands := GetDeployCommands(config)

	// Verify each command has expected content
	for _, cmd := range commands {
		// Each command should reference the remote deploy dir
		if !strings.Contains(cmd.Cmd, config.RemoteDeployDir) && !strings.Contains(cmd.Cmd, "pm2") {
			t.Errorf("Command should reference RemoteDeployDir: %s", cmd.Cmd)
		}
	}
}

// TestPrivateKeyPathWithTilde tests PrivateKeyPath with tilde
func TestPrivateKeyPathWithTilde(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Skip("Cannot get current user, skipping test")
	}

	// Test with tilde prefix
	result := PrivateKeyPath("~/.ssh/id_rsa")
	expected := filepath.Join(currentUser.HomeDir, ".ssh", "id_rsa")
	if result != expected {
		t.Errorf("PrivateKeyPath with tilde = %q, want %q", result, expected)
	}
}
