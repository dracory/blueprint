package main

import (
	"errors"

	"github.com/dracory/base/cfmt"
	cli "github.com/dracory/base/cmd"
	"github.com/samber/lo"
)

// BuildApp builds the executable for deployment
func BuildApp(config Config) error {
	cfmt.Infoln("🛠️  1. Building executable...")

	err := BuildExecutable(config.BuildLocalExecutableTempPath)
	if err != nil {
		return err
	}

	cfmt.Successln("✅  - Executable built successfully")
	return nil
}

// UploadFiles uploads any additional files needed for deployment
func UploadFiles(config Config) error {
	cfmt.Infoln("📤  2. Uploading files...")

	for _, file := range config.OtherFilesToDeploy {
		if err := UploadFileToRemoteDeployDir(config, file.LocalPath, file.RemotePath); err != nil {
			return err
		}
	}

	cfmt.Successln("✅  - Additional files uploaded")
	return nil
}

// UploadFileToRemoteDeployDir uploads a single file to the remote deploy directory
func UploadFileToRemoteDeployDir(config Config, fileLocalPath string, fileRemoteName string) error {
	cmd := `scp -o stricthostkeychecking=no -i ` + PrivateKeyPath(config.SSHKey) + ` ` + fileLocalPath + ` ` + config.SSHLogin + `:` + config.RemoteDeployDir + `/` + fileRemoteName
	cfmt.Infoln("🖥️  - Executing: " + cmd)

	output, err := cli.ExecLine(cmd)
	if err != nil {
		cfmt.Errorln("❌  - Error:", err)
		cfmt.Errorln("📝  - Output:", output)
		return errors.New("failed to upload file: " + fileLocalPath + ", error: " + err.Error())
	}
	cfmt.Successln("✅  - Output:", output)

	return nil
}

// UploadExecutable uploads the built executable to the server
func UploadExecutable(config Config) error {
	cfmt.Infoln("🚀  3. Uploading executable...")

	return UploadFileToRemoteDeployDir(config, config.BuildLocalExecutableTempPath, config.RemoteTempDeployName)
}

// ReplaceExecutable replaces the current executable with the new one on the server
func ReplaceExecutable(config Config) error {
	cfmt.Infoln("♻️  4. Replacing current executable...")

	cmds := GetDeployCommands(config)

	for _, entry := range cmds {
		cfmt.Infoln("🖥️  - Executing: " + entry.Cmd)

		output, err := SSH(config.SSHHost, config.SSHUser, config.SSHKey, entry.Cmd)

		if err != nil {
			cfmt.Errorln("❌  - Error:", err)
			cfmt.Errorln("📝  - Output:", output)
			if entry.Required {
				return err // stop on first error, if required
			}
		}

		cfmt.Successln("✅  - Output: ", lo.Ternary(output == "", "no output", output))
	}

	return nil
}
