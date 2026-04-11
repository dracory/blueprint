package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test IsFlashRoute
func TestIsFlashRoute(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "flash path returns true",
			path:     "/flash/message",
			expected: true,
		},
		{
			name:     "regular path returns false",
			path:     "/home",
			expected: false,
		},
		{
			name:     "path with flash in middle returns true",
			path:     "/api/flash/notify",
			expected: true,
		},
		{
			name:     "path ending with flash returns true",
			path:     "/admin/flash",
			expected: true,
		},
		{
			name:     "path starting with flash returns true",
			path:     "/flash",
			expected: true,
		},
		{
			name:     "complex path without flash returns false",
			path:     "/api/v1/users/profile",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			result := IsFlashRoute(req)
			if result != tt.expected {
				t.Errorf("IsFlashRoute() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test Flash Constants
func TestFlashConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{name: "FLASH_ERROR", constant: FLASH_ERROR, expected: "error"},
		{name: "FLASH_SUCCESS", constant: FLASH_SUCCESS, expected: "success"},
		{name: "FLASH_INFO", constant: FLASH_INFO, expected: "info"},
		{name: "FLASH_WARNING", constant: FLASH_WARNING, expected: "warning"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

// Test ToFlashURL with nil cache store
func TestToFlashURL_NilCacheStore(t *testing.T) {
	result := ToFlashURL(nil, FLASH_ERROR, "test message", "/redirect", 5)
	if result != "to_flash_url: cache store is nil" {
		t.Errorf("ToFlashURL(nil cache) = %v, want 'to_flash_url: cache store is nil'", result)
	}
}

// Test ToFlashErrorURL convenience function
func TestToFlashErrorURL_NilCacheStore(t *testing.T) {
	result := ToFlashErrorURL(nil, "error message", "/redirect", 5)
	if result != "to_flash_url: cache store is nil" {
		t.Errorf("ToFlashErrorURL(nil cache) = %v, want error message", result)
	}
}

// Test ToFlashInfoURL convenience function
func TestToFlashInfoURL_NilCacheStore(t *testing.T) {
	result := ToFlashInfoURL(nil, "info message", "/redirect", 5)
	if result != "to_flash_url: cache store is nil" {
		t.Errorf("ToFlashInfoURL(nil cache) = %v, want error message", result)
	}
}

// Test ToFlashSuccessURL convenience function
func TestToFlashSuccessURL_NilCacheStore(t *testing.T) {
	result := ToFlashSuccessURL(nil, "success message", "/redirect", 5)
	if result != "to_flash_url: cache store is nil" {
		t.Errorf("ToFlashSuccessURL(nil cache) = %v, want error message", result)
	}
}

// Test ToFlashWarningURL convenience function
func TestToFlashWarningURL_NilCacheStore(t *testing.T) {
	result := ToFlashWarningURL(nil, "warning message", "/redirect", 5)
	if result != "to_flash_url: cache store is nil" {
		t.Errorf("ToFlashWarningURL(nil cache) = %v, want error message", result)
	}
}
