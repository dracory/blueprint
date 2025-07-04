package main

import (
	"log"
	"os"
	"os/exec"
	"os/user"

	"github.com/dracory/base/cmd"
	"github.com/dromara/carbon/v2"
	"github.com/mingrammer/cfmt"
	"github.com/samber/lo"
	"github.com/sfreiberg/simplessh"
)

var timestamp = carbon.Now(carbon.UTC).Format("Ymd_His")
var buildExecutablePath = "tmp/application_deploy_" + timestamp
var sshKey = "{{ SSHKEY }}.prv"
var sshUser = "{{ SSHUSER }}"
var sshHost = "{{ SSHHOST }}"
var sshLogin = sshUser + "@" + sshHost
var remAppDir = "{{ APP_NAME }}"
var remDeployDir = "/home/" + sshUser + "/" + remAppDir
var remTempDeployName = "temp_deploy_" + timestamp
var pm2ProcessName = "{{ PROCESSNAME }}"
var otherFilesToDeploy = []struct {
	LocalPath  string
	RemotePath string
}{}

func main() {
	cfmt.Infoln("1. Building executable...")

	err := buildExecutable(buildExecutablePath)

	if err != nil {
		log.Fatal(err)
		return
	}

	cfmt.Infoln("2. Uploading files...")

	for _, file := range otherFilesToDeploy {
		command := `scp -o stricthostkeychecking=no -i ` + privateKeyPath(sshKey) + ` ` + file.LocalPath + ` ` + sshLogin + `:` + remDeployDir + `/` + file.RemotePath
		cfmt.Infoln(" - Executing:" + command)
		cmd.ExecLine(command)
	}

	cfmt.Infoln("3. Uploading executable...")

	command := `scp -o stricthostkeychecking=no -i ` + privateKeyPath(sshKey) + ` ` + buildExecutablePath + ` ` + sshLogin + `:` + remDeployDir + `/` + remTempDeployName

	cfmt.Infoln(" - Executing:" + command)
	cmd.ExecLine(command)

	cfmt.Infoln("3. Replace current executable...")

	commands := []struct {
		command  string
		required bool
	}{
		{
			command:  `chmod 750 ` + remDeployDir + `/` + remTempDeployName,
			required: true,
		},
		{
			command:  `mv ` + remDeployDir + `/application  ` + remDeployDir + `/` + timestamp + `_backup_application`,
			required: true,
		},
		{
			command:  `mv ` + remDeployDir + `/` + remTempDeployName + `  ` + remDeployDir + `/application`,
			required: true,
		},
		{
			command:  `mv ` + remDeployDir + `/application.error.log ` + remDeployDir + `/` + timestamp + `_backup_application.error.log`,
			required: false,
		},
		{
			command:  `mv ` + remDeployDir + `/application.log ` + remDeployDir + `/` + timestamp + `_backup_application.log`,
			required: false,
		},
		// {
		// 	command:  `cd ` + remDeployDir + ` && pm2 start "application" --name ` + pm2ProcessName + ` --log="` + remDeployDir + `/application.log" --error="` + remDeployDir + `/application.error.log" --time`,
		// 	required: true,
		// },
		{
			command:  `pm2 restart ` + pm2ProcessName,
			required: true,
		},
	}

	for _, entry := range commands {
		cfmt.Infoln(" - Executing:" + entry.command)

		output, error := ssh(sshHost, sshUser, sshKey, entry.command)

		if error != nil {
			cfmt.Errorln("  - Error:", error)
			cfmt.Errorln("  - Output: ", output)
			if entry.required {
				return // stop on first error, if required
			}
		}

		cfmt.Successln("  - Output: ", lo.Ternary(output == "", "no output", output))
	}

	cfmt.Infoln("Deployed!")

}

// buildExecutable builds an executable at the specified path.
//
// Parameters:
// - pathExec: string - the path where the executable will be built.
//
// Returns:
// - error - error if the build process encounters any issues.
func buildExecutable(pathExec string) error {
	newEnv := os.Environ()
	newEnv = append(newEnv, "GOOS=linux")
	newEnv = append(newEnv, "GOARCH=amd64")
	newEnv = append(newEnv, "CGO_ENABLED=0")

	cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-v", "-o", pathExec, "main.go")
	cmd.Env = newEnv
	out, err := cmd.CombinedOutput()

	if err != nil {
		cfmt.Errorln(string(out))
	} else {
		cfmt.Successln(string(out))
	}

	return err
}

// privateKeyPath returns the full path of the private key for the given SSH key.
//
// Parameters:
// - sshKey: a string representing the name of the SSH key
//
// Returns:
// - string: the full path to the private key
func privateKeyPath(sshKey string) string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err.Error())
	}
	homeDirectory := user.HomeDir
	privateKeyPath := homeDirectory + "/.ssh/" + sshKey
	return privateKeyPath
}

// ssh connects to an SSH server, executes a command, and returns the output.
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
func ssh(sshHost, sshUser, sshKey, cmd string) (output string, err error) {
	client, err := simplessh.ConnectWithKeyFile(sshHost+":22", sshUser, privateKeyPath(sshKey))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	outputBytes, err := client.Exec(cmd)

	if err != nil {
		return string(outputBytes), err
	}

	return string(outputBytes), nil
}
