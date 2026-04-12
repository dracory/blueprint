package stats

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"project/internal/testutils"
)

func TestNewStatsVisitorEnhanceTask(t *testing.T) {
	// Test with nil registry - this will call log.Fatal which terminates the program
	// We cannot test this directly without mocking log.Fatal
	// Test with valid registry instead
	registry := testutils.Setup()
	task := NewStatsVisitorEnhanceTask(registry)
	if task == nil {
		t.Error("NewStatsVisitorEnhanceTask() should not return nil")
	}
	if task.registry != registry {
		t.Error("Task registry should match the provided registry")
	}
}

func TestStatsVisitorEnhanceTask_Alias(t *testing.T) {
	registry := testutils.Setup()
	task := NewStatsVisitorEnhanceTask(registry)

	alias := task.Alias()
	if alias != "StatsVisitorEnhanceTask" {
		t.Errorf("Alias() = %q, want %q", alias, "StatsVisitorEnhanceTask")
	}
}

func TestStatsVisitorEnhanceTask_Title(t *testing.T) {
	registry := testutils.Setup()
	task := NewStatsVisitorEnhanceTask(registry)

	title := task.Title()
	if title != "Stats Visitor Enhance" {
		t.Errorf("Title() = %q, want %q", title, "Stats Visitor Enhance")
	}
}

func TestStatsVisitorEnhanceTask_Description(t *testing.T) {
	registry := testutils.Setup()
	task := NewStatsVisitorEnhanceTask(registry)

	description := task.Description()
	if description != "Enhances the visitor stats by adding the country" {
		t.Errorf("Description() = %q, want %q", description, "Enhances the visitor stats by adding the country")
	}
}

func TestStatsVisitorEnhanceTask_Enqueue(t *testing.T) {
	// Test with task that has nil registry
	task := &statsVisitorEnhanceTask{registry: nil}
	_, err := task.Enqueue()
	if err == nil {
		t.Error("Enqueue() with nil registry should return error")
	}
	if err.Error() != "task store is nil" {
		t.Errorf("Enqueue() error = %q, want 'task store is nil'", err.Error())
	}

	// Test with registry but without task store
	registry := testutils.Setup()
	task = NewStatsVisitorEnhanceTask(registry)
	_, err = task.Enqueue()
	if err == nil {
		t.Error("Enqueue() without task store should return error")
	}
	if err.Error() != "task store is nil" {
		t.Errorf("Enqueue() error = %q, want 'task store is nil'", err.Error())
	}

	// Test with registry with task store
	registry = testutils.Setup(testutils.WithTaskStore(true))
	task = NewStatsVisitorEnhanceTask(registry)
	result, err := task.Enqueue()
	// Task may not be registered, but we should get either success or a specific error
	if err != nil && err.Error() != "task store is nil" {
		t.Logf("Enqueue() with task store returned expected error: %v", err)
	}
	if err == nil && result == nil {
		t.Error("Enqueue() should return a TaskQueueInterface when successful")
	}
}

func TestStatsVisitorEnhanceTask_Handle_NilRegistry(t *testing.T) {
	task := &statsVisitorEnhanceTask{registry: nil}
	result := task.Handle()
	if result != false {
		t.Error("Handle() with nil registry should return false")
	}
}

func TestStatsVisitorEnhanceTask_Handle_NilStatsStore(t *testing.T) {
	registry := testutils.Setup()
	task := NewStatsVisitorEnhanceTask(registry)
	result := task.Handle()
	if result != false {
		t.Error("Handle() with nil stats store should return false")
	}
}

func TestStatsVisitorEnhanceTask_FindCountryByIp_EmptyIP(t *testing.T) {
	registry := testutils.Setup()
	task := NewStatsVisitorEnhanceTask(registry)

	country := task.findCountryByIp(context.TODO(), "")
	if country != "UN" {
		t.Errorf("findCountryByIp() with empty IP = %q, want %q", country, "UN")
	}
}

func TestStatsVisitorEnhanceTask_FindCountryByIp_LocalhostIP(t *testing.T) {
	registry := testutils.Setup()
	task := NewStatsVisitorEnhanceTask(registry)

	country := task.findCountryByIp(context.TODO(), "127.0.0.1")
	if country != "UN" {
		t.Errorf("findCountryByIp() with localhost IP = %q, want %q", country, "UN")
	}
}

func TestStatsVisitorEnhanceTask_ProcessVisitor_NilRegistry(t *testing.T) {
	task := &statsVisitorEnhanceTask{registry: nil}
	result := task.processVisitor(context.TODO(), nil)
	if result != false {
		t.Error("processVisitor() with nil registry should return false")
	}
}

func TestStatsVisitorEnhanceTask_ProcessVisitor_NilStatsStore(t *testing.T) {
	registry := testutils.Setup()
	task := NewStatsVisitorEnhanceTask(registry)
	result := task.processVisitor(context.TODO(), nil)
	if result != false {
		t.Error("processVisitor() with nil stats store should return false")
	}
}

// TestConstants tests the package constants
func TestConstants(t *testing.T) {
	if ipLookupEndpoint != "https://ip2c.org/" {
		t.Errorf("ipLookupEndpoint = %q, want %q", ipLookupEndpoint, "https://ip2c.org/")
	}
	if ipLookupTimeout != 5*time.Second {
		t.Errorf("ipLookupTimeout = %v, want 5 seconds", ipLookupTimeout)
	}
}

// TestIPLookupHTTPClient tests the HTTP client configuration
func TestIPLookupHTTPClient(t *testing.T) {
	// Note: Tests run in parallel, so we only verify the client exists
	// and has correct configuration without mutating global state
	if ipLookupHTTPClient == nil {
		t.Error("ipLookupHTTPClient should not be nil")
		return
	}
	// Verify timeout is configured (read-only check, no race condition)
	if ipLookupHTTPClient.Timeout != 5*time.Second {
		t.Errorf("ipLookupHTTPClient.Timeout = %v, want 5 seconds", ipLookupHTTPClient.Timeout)
	}
}

// mockRoundTripper is a test HTTP transport that returns canned responses
type mockRoundTripper struct {
	response *http.Response
	err      error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

// TestStatsVisitorEnhanceTask_FindCountryByIp_WithMock tests IP lookup with mock HTTP client
func TestStatsVisitorEnhanceTask_FindCountryByIp_WithMock(t *testing.T) {
	registry := testutils.Setup()
	task := NewStatsVisitorEnhanceTask(registry)

	tests := []struct {
		name       string
		response   *http.Response
		err        error
		wantResult string
	}{
		{
			name:       "successful lookup",
			response:   &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString("1;US;USA;United States"))},
			wantResult: "US",
		},
		{
			name:       "empty country code",
			response:   &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString("1;;USA;"))},
			wantResult: "UN",
		},
		{
			name:       "error response",
			response:   &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(bytes.NewBufferString(""))},
			wantResult: "ER",
		},
		{
			name:       "network error",
			err:        http.ErrHandlerTimeout,
			wantResult: "ER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task.httpClient = &http.Client{
				Transport: &mockRoundTripper{response: tt.response, err: tt.err},
			}
			result := task.findCountryByIp(context.TODO(), "1.2.3.4")
			if result != tt.wantResult {
				t.Errorf("findCountryByIp() = %q, want %q", result, tt.wantResult)
			}
		})
	}
}

// TestStatsVisitorEnhanceTask_MultipleInstances tests creating multiple task instances
func TestStatsVisitorEnhanceTask_MultipleInstances(t *testing.T) {
	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	task1 := NewStatsVisitorEnhanceTask(registry1)
	task2 := NewStatsVisitorEnhanceTask(registry2)

	if task1 == task2 {
		t.Error("Multiple instances should be independent")
	}

	if task1.registry != registry1 {
		t.Error("Task1 should have registry1")
	}

	if task2.registry != registry2 {
		t.Error("Task2 should have registry2")
	}
}

// TestStatsVisitorEnhanceTask_StructFields tests the task struct fields
func TestStatsVisitorEnhanceTask_StructFields(t *testing.T) {
	registry := testutils.Setup()
	task := NewStatsVisitorEnhanceTask(registry)

	if task.registry != registry {
		t.Error("Task registry field should match provided registry")
	}
}
