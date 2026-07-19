package cli

import (
	"fmt"
	"project/internal/app"
	"project/internal/cmds"
	"project/internal/emails"
	"project/internal/routes"

	baseCfmt "github.com/dracory/base/cfmt"
	"github.com/dracory/base/cli"
)

// Constants for command names
const (
	CommandTask        = "task"
	CommandJob         = "job"
	CommandRoutes      = "routes"
	CommandMaintenance = "maintenance"
	SubcommandList     = "list"
)

// NewDispatcher creates a new CLI dispatcher with blueprint-specific commands registered.
func NewDispatcher() *cli.Dispatcher[app.AppInterface] {
	dispatcher := cli.NewDispatcher[app.AppInterface]()

	// Register blueprint-specific commands
	dispatcher.RegisterCommand(CommandTask, "Execute a task by alias", handleTaskCommand)
	dispatcher.RegisterCommand(CommandJob, "Execute a job with arguments", handleJobCommand)
	dispatcher.RegisterCommand(CommandRoutes, "List all registered routes", handleRoutesCommand)
	dispatcher.RegisterCommand(CommandMaintenance, "Manage maintenance mode", handleMaintenanceCommand)

	return dispatcher
}

// ExecuteCliCommand executes a CLI command using the generic dispatcher.
//
// This function creates a new dispatcher with blueprint-specific commands
// and delegates execution to the generic dispatcher.
//
// Parameters:
// - app app.AppInterface : The app instance to be passed to command handlers.
// - args []string : The command line arguments (excluding the program name).
//
// Returns:
// - error: An error if the command execution fails or is invalid, otherwise nil.
func ExecuteCliCommand(app app.AppInterface, args []string) error {
	dispatcher := NewDispatcher()
	return dispatcher.ExecuteCommand(app, args)
}

// handleTaskCommand handles the 'task' command.
func handleTaskCommand(app app.AppInterface, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing task alias for command '%s'", CommandTask)
	}
	if app.GetTaskStore() == nil {
		err := fmt.Errorf("task store is nil")
		baseCfmt.Errorln(err.Error())
		return err
	}

	// Initialize email sender for tasks that require it
	emails.InitEmailSender(app)

	taskAlias := args[0]
	taskArgs := args[1:]
	// Assuming TaskExecuteCli handles its own errors/logging internally
	app.GetTaskStore().TaskDefinitionExecuteCli(taskAlias, taskArgs)
	// Assuming success unless TaskExecuteCli panics or indicates failure differently
	return nil
}

// handleJobCommand handles the 'job' command.
func handleJobCommand(app app.AppInterface, args []string) error {
	// Assuming ExecuteJob handles its own errors/logging internally
	cmds.ExecuteJob(app, args)
	// Assuming success unless ExecuteJob panics or indicates failure differently
	return nil
}

// handleRoutesCommand handles the 'routes' command.
func handleRoutesCommand(app app.AppInterface, args []string) error {
	if len(args) == 0 || args[0] != SubcommandList {
		return fmt.Errorf("invalid or missing subcommand for '%s'. Use '%s %s'", CommandRoutes, CommandRoutes, SubcommandList)
	}

	r := routes.Router(app)
	r.List()

	return nil
}
