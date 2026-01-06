package main

import (
	"github.com/dromara/carbon/v2"
)

// InitConfig initializes and returns the configuration for deployment
func InitConfig() Config {
	timestamp := carbon.Now(carbon.UTC).Format("Ymd_His")
	sshLogin := SSH_USER + "@" + SSH_HOST
	remoteDeployDir := "/home/" + SSH_USER + "/" + REMOTE_APP_DIR

	return Config{
		Timestamp:                    timestamp,
		BuildLocalExecutableTempPath: "tmp/application_deploy_" + timestamp,
		SSHKey:                       SSH_KEY,
		SSHUser:                      SSH_USER,
		SSHHost:                      SSH_HOST,
		SSHLogin:                     sshLogin,
		RemoteAppDir:                 REMOTE_APP_DIR,
		RemoteDeployDir:              remoteDeployDir,
		RemoteTempDeployName:         "temp_deploy_" + timestamp,
		PM2ProcessName:               PM2_PROCESS_NAME,
		OtherFilesToDeploy:           OTHER_FILES_TO_DEPLOY,
	}
}

// GetDeployCommands returns the list of commands to be executed during deployment
func GetDeployCommands(config Config) []DeployCommand {
	return []DeployCommand{
		{
			Reason:   "Changing permissions to 750 of the temp file",
			Cmd:      `chmod 750 ` + config.RemoteDeployDir + `/` + config.RemoteTempDeployName,
			Required: true,
		},
		{
			Reason:   "Touch current executable to make sure it exists (only needed for the first time)",
			Cmd:      `touch ` + config.RemoteDeployDir + `/application`,
			Required: true,
		},
		{
			Reason:   "Rename current executable to backup (in case of failure we can restore it manually)",
			Cmd:      `mv ` + config.RemoteDeployDir + `/application  ` + config.RemoteDeployDir + `/` + config.Timestamp + `_backup_application`,
			Required: true,
		},
		{
			Reason:   "Rename temp file to current executable",
			Cmd:      `mv ` + config.RemoteDeployDir + `/` + config.RemoteTempDeployName + `  ` + config.RemoteDeployDir + `/application`,
			Required: true,
		},
		{
			Reason:   "Touch current error log to make sure it exists (only needed for the first time)",
			Cmd:      `touch ` + config.RemoteDeployDir + `/registry.error.log`,
			Required: false,
		},
		{
			Reason:   "Rename current error log to backup (in case of failure we can restore it manually)",
			Cmd:      `mv ` + config.RemoteDeployDir + `/registry.error.log ` + config.RemoteDeployDir + `/` + config.Timestamp + `_backup_registry.error.log`,
			Required: false,
		},
		{
			Reason:   "Touch current log to make sure it exists (only needed for the first time)",
			Cmd:      `touch ` + config.RemoteDeployDir + `/registry.log`,
			Required: false,
		},
		{
			Reason:   "Rename current log to backup (in case of failure we can restore it manually)",
			Cmd:      `mv ` + config.RemoteDeployDir + `/registry.log ` + config.RemoteDeployDir + `/` + config.Timestamp + `_backup_registry.log`,
			Required: false,
		},
		{
			Reason:   "Delete pm2 process for the old executable",
			Cmd:      `pm2 delete ` + config.PM2ProcessName,
			Required: false,
		},
		{
			Reason:   "Start pm2 process for the new executable",
			Cmd:      `cd ` + config.RemoteDeployDir + `; pm2 start "application" --name ` + config.PM2ProcessName + ` --log=` + config.RemoteDeployDir + `/registry.log --error=` + config.RemoteDeployDir + `/registry.error.log --time`,
			Required: true,
		},
	}
}
