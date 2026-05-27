package main

import (
	"os"
	"testing"
)

func TestIsCliMode(t *testing.T) {
	os.Args = []string{"main", "task", "testTask"}
	if !isCliMode() {
		t.Errorf("isCliMode() should return true")
	}

	os.Args = []string{"main"}
	if isCliMode() {
		t.Errorf("isCliMode() should return false")
	}

	// Test with empty args
	os.Args = []string{}
	if isCliMode() {
		t.Errorf("isCliMode() with empty args should return false")
	}

	// Test with single argument (program name only)
	os.Args = []string{"./server"}
	if isCliMode() {
		t.Errorf("isCliMode() with single arg should return false")
	}
}

func TestIsCliMode_VariousArguments_NoArguments(t *testing.T) {
	os.Args = []string{}
	result := isCliMode()
	if result != false {
		t.Errorf("isCliMode() = %v, want false", result)
	}
}

func TestIsCliMode_VariousArguments_ProgramNameOnly(t *testing.T) {
	os.Args = []string{"server"}
	result := isCliMode()
	if result != false {
		t.Errorf("isCliMode() = %v, want false", result)
	}
}

func TestIsCliMode_VariousArguments_OneArgument(t *testing.T) {
	os.Args = []string{"server", "task"}
	result := isCliMode()
	if result != true {
		t.Errorf("isCliMode() = %v, want true", result)
	}
}

func TestIsCliMode_VariousArguments_MultipleArguments(t *testing.T) {
	os.Args = []string{"server", "task", "myTask"}
	result := isCliMode()
	if result != true {
		t.Errorf("isCliMode() = %v, want true", result)
	}
}

func TestIsCliMode_VariousArguments_JobCommand(t *testing.T) {
	os.Args = []string{"server", "job"}
	result := isCliMode()
	if result != true {
		t.Errorf("isCliMode() = %v, want true", result)
	}
}
