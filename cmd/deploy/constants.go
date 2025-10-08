package main

// SSH Configuration
const (
	// SSH private key filename in the .ssh directory
	SSH_KEY = "{{ SSHKEY }}.prv"

	// SSH username for the server
	SSH_USER = "{{ SSHUSER }}"

	// SSH host to connect to
	SSH_HOST = "{{ SSHHOST }}"
)

// Application Configuration
const (
	// Remote application directory name (i.e. example.com)
	REMOTE_APP_DIR = "{{ APP_NAME }}"

	// PM2 process name for the application
	PM2_PROCESS_NAME = REMOTE_APP_DIR
)

// Files to deploy in addition to the main executable
var OTHER_FILES_TO_DEPLOY = []DeployFile{
	// Add files to deploy here, for example:
	// {LocalPath: "config.json", RemotePath: "config.json"},
}
