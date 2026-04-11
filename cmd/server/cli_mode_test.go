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

func TestIsCliMode_VariousArguments(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "No arguments",
			args:     []string{},
			expected: false,
		},
		{
			name:     "Program name only",
			args:     []string{"server"},
			expected: false,
		},
		{
			name:     "Program name with one argument",
			args:     []string{"server", "task"},
			expected: true,
		},
		{
			name:     "Program name with multiple arguments",
			args:     []string{"server", "task", "myTask"},
			expected: true,
		},
		{
			name:     "Program name with job command",
			args:     []string{"server", "job"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			result := isCliMode()
			if result != tt.expected {
				t.Errorf("isCliMode() = %v, want %v", result, tt.expected)
			}
		})
	}
}
