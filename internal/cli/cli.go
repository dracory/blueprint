package cli

import (
	"errors"
	"fmt"
	"project/internal/cmds"
	"project/internal/emails"
	"project/internal/registry"
	"project/internal/routes"

	"github.com/dracory/base/cfmt"
)

// Constants for command names
const (
	CommandTask    = "task"
	CommandJob     = "job"
	CommandRoutes  = "routes"
	SubcommandList = "list"
)

// commandHandler defines the function signature for command handlers.
type commandHandler func(registry registry.RegistryInterface, args []string) error

// commandHandlers maps command strings to their handler functions.
var commandHandlers = map[string]commandHandler{
	CommandTask:   handleTaskCommand,
	CommandJob:    handleJobCommand,
	CommandRoutes: handleRoutesCommand,
}

// ExecuteCliCommand executes a CLI command based on the provided arguments using a handler map.
//
// The command can be one of the following:
// - task <alias> <arguments>
// - job <arguments>
// - routes list
//
// Business logic:
// 1. Logs the command being executed.
// 2. Validates that at least one argument (the command) is provided.
// 3. Looks up the command in the `commandHandlers` map.
// 4. If a handler is found, executes it with the remaining arguments.
// 5. If no handler is found, returns an "unrecognized command" error.
// 6. Returns specific errors for invalid commands, missing arguments, or nil TaskStore via the handlers.
//
// Parameters:
// - registry registry.RegistryInterface : The registry instance to be passed to command handlers.
// - args []string : The command line arguments (excluding the program name).
//
// Returns:
// - error: An error if the command execution fails or is invalid, otherwise nil.
func ExecuteCliCommand(registry registry.RegistryInterface, args []string) error {
	cfmt.Infoln("Executing command: ", args)

	if len(args) == 0 {
		cfmt.Warningln("No command provided.")
		// Optionally, print usage instructions here
		return errors.New("no command provided")
	}

	command := args[0]
	remainingArgs := args[1:] // Arguments after the main command

	// Look up the handler for the command
	handler, found := commandHandlers[command]
	if !found {
		err := fmt.Errorf("unrecognized command: %s", command)
		cfmt.Warningln(err.Error())
		// Optionally, print usage instructions here
		return err
	}

	// Execute the found handler with registry
	return handler(registry, remainingArgs)
}

// handleTaskCommand handles the 'task' command.
func handleTaskCommand(registry registry.RegistryInterface, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing task alias for command '%s'", CommandTask)
	}
	if registry.GetTaskStore() == nil {
		err := errors.New("task store is nil")
		cfmt.Errorln(err.Error())
		return err
	}

	// Initialize email sender for tasks that require it
	emails.InitEmailSender(registry)

	taskAlias := args[0]
	taskArgs := args[1:]
	// Assuming TaskExecuteCli handles its own errors/logging internally
	registry.GetTaskStore().TaskDefinitionExecuteCli(taskAlias, taskArgs)
	// Assuming success unless TaskExecuteCli panics or indicates failure differently
	return nil
}

// handleJobCommand handles the 'job' command.
func handleJobCommand(registry registry.RegistryInterface, args []string) error {
	// Assuming ExecuteJob handles its own errors/logging internally
	cmds.ExecuteJob(registry, args)
	// Assuming success unless ExecuteJob panics or indicates failure differently
	return nil
}

// handleRoutesCommand handles the 'routes' command.
func handleRoutesCommand(registry registry.RegistryInterface, args []string) error {
	if len(args) == 0 || args[0] != SubcommandList {
		return fmt.Errorf("invalid or missing subcommand for '%s'. Use '%s %s'", CommandRoutes, CommandRoutes, SubcommandList)
	}
	// m, r := routes.RoutesList()
	// router.List(m, r)
	r := routes.Router(registry)
	r.List()

	return nil
}
