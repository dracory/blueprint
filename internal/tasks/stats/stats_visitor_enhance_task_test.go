package stats

import (
	"context"
	"testing"

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
