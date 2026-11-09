package helpers

import (
	"os"
	"testing"
)

func TestIsCliMode(t *testing.T) {
	os.Args = []string{"main", "task", "testTask"}
	if !IsCliMode() {
		t.Errorf("IsCliMode() should return true")
	}

	os.Args = []string{"main"}
	if IsCliMode() {
		t.Errorf("IsCliMode() should return false")
	}

	// Test with empty args
	os.Args = []string{}
	if IsCliMode() {
		t.Errorf("IsCliMode() with empty args should return false")
	}

	// Test with single argument (program name only)
	os.Args = []string{"./server"}
	if IsCliMode() {
		t.Errorf("IsCliMode() with single arg should return false")
	}
}

func TestIsCliMode_VariousArguments_NoArguments(t *testing.T) {
	os.Args = []string{}
	result := IsCliMode()
	if result != false {
		t.Errorf("IsCliMode() = %v, want false", result)
	}
}

func TestIsCliMode_VariousArguments_ProgramNameOnly(t *testing.T) {
	os.Args = []string{"server"}
	result := IsCliMode()
	if result != false {
		t.Errorf("IsCliMode() = %v, want false", result)
	}
}

func TestIsCliMode_VariousArguments_OneArgument(t *testing.T) {
	os.Args = []string{"server", "task"}
	result := IsCliMode()
	if result != true {
		t.Errorf("IsCliMode() = %v, want true", result)
	}
}

func TestIsCliMode_VariousArguments_MultipleArguments(t *testing.T) {
	os.Args = []string{"server", "task", "myTask"}
	result := IsCliMode()
	if result != true {
		t.Errorf("IsCliMode() = %v, want true", result)
	}
}

func TestIsCliMode_VariousArguments_JobCommand(t *testing.T) {
	os.Args = []string{"server", "job"}
	result := IsCliMode()
	if result != true {
		t.Errorf("IsCliMode() = %v, want true", result)
	}
}
