package main

// Config holds all configuration variables for deployment
type Config struct {
	Timestamp                    string
	BuildLocalExecutableTempPath string
	SSHKey                       string
	SSHUser                      string
	SSHHost                      string
	SSHLogin                     string
	RemoteAppDir                 string
	RemoteDeployDir              string
	RemoteTempDeployName         string
	PM2ProcessName               string
	OtherFilesToDeploy           []DeployFile
}

// DeployFile represents a file to be uploaded
// from local to remote server
type DeployFile struct {
	// Local path to the file
	LocalPath string

	// Remote path to the file (relative to the remote deploy directory)
	RemotePath string
}

// DeployCommand represents a command to be executed
// during deployment
type DeployCommand struct {
	Reason   string
	Cmd      string
	Required bool
}
