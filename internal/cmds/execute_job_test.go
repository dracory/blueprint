package cmds

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestExecuteJob_MissingTaskID(t *testing.T) {
	// Create a buffer to capture the log output
	logOutput := &bytes.Buffer{}

	// Save the original logger and replace it with our test logger
	oldLogger := log.Default()
	log.SetOutput(logOutput)
	defer log.SetOutput(oldLogger.Writer())

	// Setup test environment and execute
	registry := testutils.Setup()
	ExecuteJob(registry, []string{})

	// Get the captured output
	output := logOutput.String()

	// Check if the output contains the expected string
	expectedOutput := "Task ID is required and must be the first argument"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected log output to contain %q, got: %q", expectedOutput, output)
	}
}

func TestExecuteJob_TaskNotFound(t *testing.T) {
	// Create a buffer to capture the log output
	logOutput := &bytes.Buffer{}

	// Save the original logger and replace it with our test logger
	oldLogger := log.Default()
	log.SetOutput(logOutput)
	defer log.SetOutput(oldLogger.Writer())

	// Setup test environment and execute
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	// Pass the task ID as a proper command-line argument
	ExecuteJob(registry, []string{"--task_id=nonexistent"})

	// Get the captured output
	output := logOutput.String()

	// Check if the output contains the expected string
	expectedOutput := "Task not found"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected log output to contain %q, got: %q", expectedOutput, output)
	}
}

func TestExecuteJob_NilTaskStore(t *testing.T) {
	// Create a buffer to capture the log output
	logOutput := &bytes.Buffer{}

	// Save the original logger and replace it with our test logger
	oldLogger := log.Default()
	log.SetOutput(logOutput)
	defer log.SetOutput(oldLogger.Writer())

	// Setup test environment WITHOUT task store
	registry := testutils.Setup()

	// Pass the task ID as a proper command-line argument
	ExecuteJob(registry, []string{"--task_id=12345"})

	// Get the captured output
	output := logOutput.String()

	// Check if the output contains the expected string
	expectedOutput := "TaskStore is nil"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected log output to contain %q, got: %q", expectedOutput, output)
	}
}

func TestExecuteJob_EmptyTaskID(t *testing.T) {
	// Create a buffer to capture the log output
	logOutput := &bytes.Buffer{}

	// Save the original logger and replace it with our test logger
	oldLogger := log.Default()
	log.SetOutput(logOutput)
	defer log.SetOutput(oldLogger.Writer())

	// Setup test environment
	cfg := testutils.DefaultConf()
	cfg.SetTaskStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	// Pass empty task ID
	ExecuteJob(registry, []string{"--task_id="})

	// Get the captured output
	output := logOutput.String()

	// Check if the output contains the expected string
	expectedOutput := "Task ID is required"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected log output to contain %q, got: %q", expectedOutput, output)
	}
}
