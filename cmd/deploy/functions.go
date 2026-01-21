package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/dracory/base/cfmt"
	"github.com/sfreiberg/simplessh"
)

// BuildExecutable builds an executable at the specified path.
//
// Parameters:
// - pathExec: string - the path where the executable will be built.
//
// Returns:
// - error - error if the build process encounters any issues.
func BuildExecutable(pathExec string) error {
	newEnv := os.Environ()
	newEnv = append(newEnv, "GOOS=linux")
	newEnv = append(newEnv, "GOARCH=amd64")
	newEnv = append(newEnv, "CGO_ENABLED=0")

	cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-v", "-o", pathExec, "./cmd/server")
	cmd.Env = newEnv
	out, err := cmd.CombinedOutput()

	if err != nil {
		cfmt.Errorln(string(out))
	} else {
		cfmt.Successln(string(out))
	}

	return err
}

// validateCommand checks if the provided command is in the allowed commands list.
//
// Parameters:
// - cmd: the command to validate
//
// Returns:
// - error: error if command is not allowed, nil otherwise
func validateCommand(cmd string) error {
	// Define allowed commands for security
	allowedCommands := map[string]bool{
		"ls":     true,
		"pwd":    true,
		"cat":    true,
		"grep":   true,
		"find":   true,
		"ps":     true,
		"df":     true,
		"du":     true,
		"whoami": true,
		"id":     true,
		"date":   true,
		"uptime": true,
		"top":    true,
		"free":   true,
		"uname":  true,
	}

	// Check for dangerous characters first
	dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "<", ">", "\"", "'"}
	for _, char := range dangerousChars {
		if strings.Contains(cmd, char) {
			return errors.New("dangerous character detected in command: " + char)
		}
	}

	// Extract the base command (first word before any arguments)
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return errors.New("empty command not allowed")
	}

	baseCommand := parts[0]

	// Check if the base command is allowed
	if !allowedCommands[baseCommand] {
		return errors.New("command not allowed: " + baseCommand)
	}

	return nil
}

// PrivateKeyPath returns the full path of the private key for the given SSH key.
//
// Parameters:
// - sshKey: a string representing the name of the SSH key
//
// Returns:
// - string: the full path to the private key
func PrivateKeyPath(sshKey string) string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err.Error())
	}
	homeDirectory := user.HomeDir
	privateKeyPath := homeDirectory + "/.ssh/" + sshKey
	return privateKeyPath
}

// SSH connects to an SSH server, executes a command, and returns the output.
//
// Parameters:
// - sshHost: the hostname of the SSH server.
// - sshUser: the username to authenticate with.
// - sshKey: the path to the SSH private key file.
// - cmd: the command to execute on the SSH server.
//
// Return:
// - output: the output of the executed command.
// - err: an error, if any, nil otherwise.
func SSH(sshHost, sshUser, sshKey, cmd string) (output string, err error) {
	// Validate command before execution to prevent command injection
	if err := validateCommand(cmd); err != nil {
		return "", err
	}

	client, err := simplessh.ConnectWithKeyFile(sshHost+":22", sshUser, PrivateKeyPath(sshKey))
	if err != nil {
		panic(err)
	}
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			log.Printf("Warning: failed to close SSH client: %v", closeErr)
		}
	}()

	outputBytes, err := client.Exec(cmd)

	if err != nil {
		return string(outputBytes), err
	}

	return string(outputBytes), nil
}
