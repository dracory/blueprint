package main

import (
	"testing"
)

// TestConfigStruct tests Config struct fields
func TestConfigStruct(t *testing.T) {
	t.Parallel()
	config := Config{
		Timestamp:                    "20240101_120000",
		BuildLocalExecutableTempPath: "tmp/test",
		SSHKey:                       "test_key",
		SSHUser:                      "test_user",
		SSHHost:                      "test_host",
		SSHLogin:                     "test_user@test_host",
		RemoteAppDir:                 "test_app",
		RemoteDeployDir:              "/home/test_user/test_app",
		RemoteTempDeployName:         "temp_deploy_20240101_120000",
		PM2ProcessName:               "test_app",
		OtherFilesToDeploy:           []DeployFile{},
	}

	if config.Timestamp != "20240101_120000" {
		t.Error("Timestamp field not set correctly")
	}
	if config.BuildLocalExecutableTempPath != "tmp/test" {
		t.Error("BuildLocalExecutableTempPath field not set correctly")
	}
	if config.SSHKey != "test_key" {
		t.Error("SSHKey field not set correctly")
	}
	if config.SSHUser != "test_user" {
		t.Error("SSHUser field not set correctly")
	}
	if config.SSHHost != "test_host" {
		t.Error("SSHHost field not set correctly")
	}
	if config.SSHLogin != "test_user@test_host" {
		t.Error("SSHLogin field not set correctly")
	}
	if config.RemoteAppDir != "test_app" {
		t.Error("RemoteAppDir field not set correctly")
	}
	if config.RemoteDeployDir != "/home/test_user/test_app" {
		t.Error("RemoteDeployDir field not set correctly")
	}
	if config.RemoteTempDeployName != "temp_deploy_20240101_120000" {
		t.Error("RemoteTempDeployName field not set correctly")
	}
	if config.PM2ProcessName != "test_app" {
		t.Error("PM2ProcessName field not set correctly")
	}
	if len(config.OtherFilesToDeploy) != 0 {
		t.Error("OtherFilesToDeploy should be empty")
	}
}

// TestDeployFileStruct tests DeployFile struct
func TestDeployFileStruct(t *testing.T) {
	t.Parallel()
	file := DeployFile{
		LocalPath:  "/local/path/file.txt",
		RemotePath: "file.txt",
	}

	if file.LocalPath != "/local/path/file.txt" {
		t.Error("LocalPath field not set correctly")
	}
	if file.RemotePath != "file.txt" {
		t.Error("RemotePath field not set correctly")
	}
}

// TestDeployCommandStruct tests DeployCommand struct
func TestDeployCommandStruct(t *testing.T) {
	t.Parallel()
	cmd := DeployCommand{
		Reason:   "Test reason",
		Cmd:      "echo test",
		Required: true,
	}

	if cmd.Reason != "Test reason" {
		t.Error("Reason field not set correctly")
	}
	if cmd.Cmd != "echo test" {
		t.Error("Cmd field not set correctly")
	}
	if !cmd.Required {
		t.Error("Required field should be true")
	}
}

// TestDeployFileSlice tests slice of DeployFile
func TestDeployFileSlice(t *testing.T) {
	t.Parallel()
	files := []DeployFile{
		{LocalPath: "file1.txt", RemotePath: "remote1.txt"},
		{LocalPath: "file2.txt", RemotePath: "remote2.txt"},
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}
}

// TestDeployCommandSlice tests slice of DeployCommand
func TestDeployCommandSlice(t *testing.T) {
	t.Parallel()
	commands := []DeployCommand{
		{Reason: "First", Cmd: "cmd1", Required: true},
		{Reason: "Second", Cmd: "cmd2", Required: false},
	}

	if len(commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(commands))
	}

	if commands[0].Required != true {
		t.Error("First command should be required")
	}

	if commands[1].Required != false {
		t.Error("Second command should not be required")
	}
}
