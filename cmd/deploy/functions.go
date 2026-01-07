package main

import (
	"log"
	"os"
	"os/exec"
	"os/user"

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
