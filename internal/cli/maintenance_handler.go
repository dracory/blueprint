package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"project/internal/app"
)

// Maintenance subcommand constants
const (
	MaintenanceStatus = "status"
)

// Subcommands that enable maintenance mode
var maintenanceEnableSubcommands = map[string]bool{
	"down":   true,
	"enable": true,
	"on":     true,
}

// Subcommands that disable maintenance mode
var maintenanceDisableSubcommands = map[string]bool{
	"up":      true,
	"disable": true,
	"off":     true,
}

// maintenanceState mirrors the JSON structure written to the state file
type maintenanceState struct {
	Message           string   `json:"message"`
	RetryAfterSeconds int      `json:"retry_after_seconds"`
	ExcludeIPs        []string `json:"exclude_ips"`
	ExcludePaths      []string `json:"exclude_paths"`
	CreatedAt         string   `json:"created_at"`
}

// handleMaintenanceCommand handles the 'maintenance' command.
func handleMaintenanceCommand(app app.AppInterface, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing subcommand for 'maintenance'. Use: maintenance enable|disable|status (aliases: on, off, down, up)")
	}

	subcommand := args[0]
	rest := args[1:]

	if maintenanceEnableSubcommands[subcommand] {
		return handleMaintenanceEnable(app, rest)
	}
	if maintenanceDisableSubcommands[subcommand] {
		return handleMaintenanceDisable(app)
	}
	if subcommand == MaintenanceStatus {
		return handleMaintenanceStatus(app)
	}

	return fmt.Errorf("unknown maintenance subcommand '%s'. Use: enable, disable, status (aliases: on, off, down, up)", subcommand)
}

func handleMaintenanceEnable(app app.AppInterface, args []string) error {
	filePath := getMaintenanceFilePath(app)

	opts := parseMaintenanceEnableArgs(args)

	state := maintenanceState{
		Message:           opts.message,
		RetryAfterSeconds: opts.retryAfter,
		ExcludeIPs:        opts.excludeIPs,
		ExcludePaths:      opts.excludePaths,
		CreatedAt:         time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal maintenance state: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write maintenance file '%s': %w", filePath, err)
	}

	fmt.Printf("Maintenance mode: ON (file: %s)\n", filePath)
	if state.Message != "" {
		fmt.Printf("  Message: %s\n", state.Message)
	}
	if state.RetryAfterSeconds > 0 {
		fmt.Printf("  Retry-After: %d seconds\n", state.RetryAfterSeconds)
	}
	if len(state.ExcludeIPs) > 0 {
		fmt.Printf("  Excluded IPs: %s\n", strings.Join(state.ExcludeIPs, ", "))
	}
	if len(state.ExcludePaths) > 0 {
		fmt.Printf("  Excluded paths: %s\n", strings.Join(state.ExcludePaths, ", "))
	}

	return nil
}

func handleMaintenanceDisable(app app.AppInterface) error {
	filePath := getMaintenanceFilePath(app)

	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Maintenance mode: OFF (file does not exist)")
			return nil
		}
		return fmt.Errorf("failed to remove maintenance file '%s': %w", filePath, err)
	}

	fmt.Printf("Maintenance mode: OFF (file removed: %s)\n", filePath)
	return nil
}

func handleMaintenanceStatus(app app.AppInterface) error {
	filePath := getMaintenanceFilePath(app)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Maintenance mode: OFF")
			fmt.Printf("  File: %s (does not exist)\n", filePath)
			return nil
		}
		return fmt.Errorf("failed to read maintenance file: %w", err)
	}

	var state maintenanceState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to parse maintenance file: %w", err)
	}

	fmt.Println("Maintenance mode: ON")
	fmt.Printf("  File: %s\n", filePath)
	if state.Message != "" {
		fmt.Printf("  Message: %s\n", state.Message)
	}
	if state.RetryAfterSeconds > 0 {
		fmt.Printf("  Retry-After: %d seconds\n", state.RetryAfterSeconds)
	}
	if len(state.ExcludeIPs) > 0 {
		fmt.Printf("  Excluded IPs: %s\n", strings.Join(state.ExcludeIPs, ", "))
	}
	if len(state.ExcludePaths) > 0 {
		fmt.Printf("  Excluded paths: %s\n", strings.Join(state.ExcludePaths, ", "))
	}
	if state.CreatedAt != "" {
		fmt.Printf("  Created at: %s\n", state.CreatedAt)
	}

	return nil
}

type maintenanceEnableOptions struct {
	message      string
	retryAfter   int
	excludeIPs   []string
	excludePaths []string
}

func parseMaintenanceEnableArgs(args []string) maintenanceEnableOptions {
	opts := maintenanceEnableOptions{
		message:    "We'll be right back.",
		retryAfter: 60,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if strings.HasPrefix(arg, "--message=") {
			opts.message = strings.TrimPrefix(arg, "--message=")
		} else if arg == "--message" && i+1 < len(args) {
			i++
			opts.message = args[i]
		} else if strings.HasPrefix(arg, "--retry=") {
			val := strings.TrimPrefix(arg, "--retry=")
			if n, err := strconv.Atoi(val); err == nil {
				opts.retryAfter = n
			}
		} else if arg == "--retry" && i+1 < len(args) {
			i++
			if n, err := strconv.Atoi(args[i]); err == nil {
				opts.retryAfter = n
			}
		} else if strings.HasPrefix(arg, "--ips=") {
			val := strings.TrimPrefix(arg, "--ips=")
			opts.excludeIPs = parseCommaList(val)
		} else if arg == "--ips" && i+1 < len(args) {
			i++
			opts.excludeIPs = parseCommaList(args[i])
		} else if strings.HasPrefix(arg, "--exclude=") {
			val := strings.TrimPrefix(arg, "--exclude=")
			opts.excludePaths = parseCommaList(val)
		} else if arg == "--exclude" && i+1 < len(args) {
			i++
			opts.excludePaths = parseCommaList(args[i])
		}
	}

	return opts
}

func parseCommaList(s string) []string {
	var result []string
	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func getMaintenanceFilePath(app app.AppInterface) string {
	if app != nil && app.GetConfig() != nil {
		fp := app.GetConfig().GetAppMaintenanceFilePath()
		if fp != "" {
			return fp
		}
	}
	return "maintenance_mode_state.json"
}
