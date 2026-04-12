package admin

import (
	"testing"

	"project/internal/testutils"
)

// TestHumanFilesize tests the HumanFilesize method
func TestHumanFilesize(t *testing.T) {
	t.Parallel()

	registry := testutils.Setup()
	controller := NewFileManagerController(registry)

	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0 B"},
		{100, "100 B"},
		{999, "999 B"},
		{1000, "1.0 kB"},
		{1500, "1.5 kB"},
		{1000000, "1.0 MB"},
		{1500000, "1.5 MB"},
		{1000000000, "1.0 GB"},
		{1500000000, "1.5 GB"},
		{1000000000000, "1.0 TB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := controller.HumanFilesize(tt.size)
			if result != tt.expected {
				t.Errorf("HumanFilesize(%d) = %q, want %q", tt.size, result, tt.expected)
			}
		})
	}
}

// TestHumanFilesize_Bytes tests byte-level formatting
func TestHumanFilesize_Bytes(t *testing.T) {
	t.Parallel()

	registry := testutils.Setup()
	controller := NewFileManagerController(registry)

	result := controller.HumanFilesize(512)
	if result != "512 B" {
		t.Errorf("HumanFilesize(512) = %q, want '512 B'", result)
	}
}

// TestHumanFilesize_LargeValues tests large file sizes
func TestHumanFilesize_LargeValues(t *testing.T) {
	t.Parallel()

	registry := testutils.Setup()
	controller := NewFileManagerController(registry)

	// Test very large values don't panic
	largeSizes := []int64{
		1 << 30, // ~1 GB
		1 << 40, // ~1 TB
		1 << 50, // ~1 PB
		1 << 60, // ~1 EB
	}

	for _, size := range largeSizes {
		result := controller.HumanFilesize(size)
		if result == "" {
			t.Errorf("HumanFilesize(%d) should not return empty string", size)
		}
	}
}
