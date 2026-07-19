package cli

import (
	"encoding/json"
	"os"
	"testing"
)

func TestHandleMaintenanceEnable_CreatesFile(t *testing.T) {
	path := "test_maintenance_enable.json"
	defer os.Remove(path)

	opts := maintenanceEnableOptions{
		message:    "Test maintenance",
		retryAfter: 30,
	}

	state := maintenanceState{
		Message:           opts.message,
		RetryAfterSeconds: opts.retryAfter,
		ExcludeIPs:        opts.excludeIPs,
		ExcludePaths:      opts.excludePaths,
		CreatedAt:         "2026-01-01T00:00:00Z",
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	readData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	var readState maintenanceState
	if err := json.Unmarshal(readData, &readState); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if readState.Message != "Test maintenance" {
		t.Fatalf("expected message 'Test maintenance', got '%s'", readState.Message)
	}
	if readState.RetryAfterSeconds != 30 {
		t.Fatalf("expected retry 30, got %d", readState.RetryAfterSeconds)
	}
}

func TestHandleMaintenanceEnable_WritesValidJSON(t *testing.T) {
	path := getMaintenanceFilePath(nil)
	defer os.Remove(path)

	err := handleMaintenanceEnable(nil, []string{
		"--message=Handler test",
		"--retry=45",
		"--ips=10.0.0.1",
		"--exclude=/api/*",
	})
	if err != nil {
		t.Fatalf("handleMaintenanceEnable failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	var state maintenanceState
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if state.Message != "Handler test" {
		t.Fatalf("expected message 'Handler test', got '%s'", state.Message)
	}
	if state.RetryAfterSeconds != 45 {
		t.Fatalf("expected retry 45, got %d", state.RetryAfterSeconds)
	}
	if len(state.ExcludeIPs) != 1 || state.ExcludeIPs[0] != "10.0.0.1" {
		t.Fatalf("expected excludeIPs [10.0.0.1], got %v", state.ExcludeIPs)
	}
	if len(state.ExcludePaths) != 1 || state.ExcludePaths[0] != "/api/*" {
		t.Fatalf("expected excludePaths [/api/*], got %v", state.ExcludePaths)
	}
	if state.CreatedAt == "" {
		t.Fatal("expected CreatedAt to be set")
	}
}

func TestHandleMaintenanceDisable_RemovesFile(t *testing.T) {
	path := getMaintenanceFilePath(nil)
	defer os.Remove(path)

	if err := os.WriteFile(path, []byte(`{"message":"test"}`), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := handleMaintenanceDisable(nil); err != nil {
		t.Fatalf("handleMaintenanceDisable failed: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("file should not exist after handleMaintenanceDisable")
	}
}

func TestHandleMaintenanceDisable_NoFile_IsNoOp(t *testing.T) {
	path := getMaintenanceFilePath(nil)
	os.Remove(path)
	defer os.Remove(path)

	if err := handleMaintenanceDisable(nil); err != nil {
		t.Fatalf("handleMaintenanceDisable should not error when file doesn't exist: %v", err)
	}
}

func TestHandleMaintenanceStatus_NoFile_ReportsOff(t *testing.T) {
	path := getMaintenanceFilePath(nil)
	os.Remove(path)
	defer os.Remove(path)

	if err := handleMaintenanceStatus(nil); err != nil {
		t.Fatalf("handleMaintenanceStatus failed: %v", err)
	}
}

func TestHandleMaintenanceStatus_FileExists_ReportsOn(t *testing.T) {
	path := getMaintenanceFilePath(nil)
	defer os.Remove(path)

	state := maintenanceState{
		Message:           "Status test",
		RetryAfterSeconds: 90,
		ExcludeIPs:        []string{"192.168.1.1"},
		CreatedAt:         "2026-01-01T00:00:00Z",
	}
	data, _ := json.MarshalIndent(state, "", "  ")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := handleMaintenanceStatus(nil); err != nil {
		t.Fatalf("handleMaintenanceStatus failed: %v", err)
	}
}

func TestHandleMaintenanceCommand_MissingSubcommand(t *testing.T) {
	err := handleMaintenanceCommand(nil, []string{})
	if err == nil {
		t.Fatal("expected error for missing subcommand")
	}
}

func TestHandleMaintenanceCommand_UnknownSubcommand(t *testing.T) {
	err := handleMaintenanceCommand(nil, []string{"invalid"})
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}

func TestParseMaintenanceEnableArgs_Defaults(t *testing.T) {
	opts := parseMaintenanceEnableArgs([]string{})

	if opts.message != "We'll be right back." {
		t.Fatalf("expected default message, got '%s'", opts.message)
	}
	if opts.retryAfter != 60 {
		t.Fatalf("expected default retry 60, got %d", opts.retryAfter)
	}
}

func TestParseMaintenanceEnableArgs_AllOptions(t *testing.T) {
	args := []string{
		"--message=Database migration",
		"--retry=120",
		"--ips=203.0.113.5,198.51.100.10",
		"--exclude=/admin/*,/api/health",
	}

	opts := parseMaintenanceEnableArgs(args)

	if opts.message != "Database migration" {
		t.Fatalf("expected message 'Database migration', got '%s'", opts.message)
	}
	if opts.retryAfter != 120 {
		t.Fatalf("expected retry 120, got %d", opts.retryAfter)
	}
	if len(opts.excludeIPs) != 2 {
		t.Fatalf("expected 2 exclude IPs, got %d", len(opts.excludeIPs))
	}
	if opts.excludeIPs[0] != "203.0.113.5" {
		t.Fatalf("expected first IP '203.0.113.5', got '%s'", opts.excludeIPs[0])
	}
	if len(opts.excludePaths) != 2 {
		t.Fatalf("expected 2 exclude paths, got %d", len(opts.excludePaths))
	}
	if opts.excludePaths[0] != "/admin/*" {
		t.Fatalf("expected first path '/admin/*', got '%s'", opts.excludePaths[0])
	}
}

func TestParseMaintenanceEnableArgs_SpaceSeparatedValues(t *testing.T) {
	args := []string{
		"--message", "Custom message",
		"--retry", "300",
		"--ips", "10.0.0.1,10.0.0.2",
		"--exclude", "/admin/*",
	}

	opts := parseMaintenanceEnableArgs(args)

	if opts.message != "Custom message" {
		t.Fatalf("expected 'Custom message', got '%s'", opts.message)
	}
	if opts.retryAfter != 300 {
		t.Fatalf("expected retry 300, got %d", opts.retryAfter)
	}
	if len(opts.excludeIPs) != 2 {
		t.Fatalf("expected 2 IPs, got %d", len(opts.excludeIPs))
	}
	if len(opts.excludePaths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(opts.excludePaths))
	}
}

func TestParseCommaList(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"a,b,c", []string{"a", "b", "c"}},
		{"a, b , c", []string{"a", "b", "c"}},
		{"single", []string{"single"}},
		{"", []string{}},
		{" , , ", []string{}},
	}

	for _, tt := range tests {
		result := parseCommaList(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("parseCommaList(%q) = %v, expected %v", tt.input, result, tt.expected)
			continue
		}
		for i, v := range result {
			if v != tt.expected[i] {
				t.Errorf("parseCommaList(%q)[%d] = %q, expected %q", tt.input, i, v, tt.expected[i])
			}
		}
	}
}

func TestGetMaintenanceFilePath_Default(t *testing.T) {
	fp := getMaintenanceFilePath(nil)
	if fp != "maintenance_mode_state.json" {
		t.Fatalf("expected default file path, got '%s'", fp)
	}
}
