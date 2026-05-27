package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test IsFlashRoute
func TestIsFlashRoute_FlashPathReturnsTrue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/flash/message", nil)
	result := IsFlashRoute(req)
	if result != true {
		t.Errorf("IsFlashRoute() = %v, want true", result)
	}
}

func TestIsFlashRoute_RegularPathReturnsFalse(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/home", nil)
	result := IsFlashRoute(req)
	if result != false {
		t.Errorf("IsFlashRoute() = %v, want false", result)
	}
}

func TestIsFlashRoute_PathWithFlashInMiddleReturnsTrue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/flash/notify", nil)
	result := IsFlashRoute(req)
	if result != true {
		t.Errorf("IsFlashRoute() = %v, want true", result)
	}
}

func TestIsFlashRoute_PathEndingWithFlashReturnsTrue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/admin/flash", nil)
	result := IsFlashRoute(req)
	if result != true {
		t.Errorf("IsFlashRoute() = %v, want true", result)
	}
}

func TestIsFlashRoute_PathStartingWithFlashReturnsTrue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/flash", nil)
	result := IsFlashRoute(req)
	if result != true {
		t.Errorf("IsFlashRoute() = %v, want true", result)
	}
}

func TestIsFlashRoute_ComplexPathWithoutFlashReturnsFalse(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
	result := IsFlashRoute(req)
	if result != false {
		t.Errorf("IsFlashRoute() = %v, want false", result)
	}
}

// Test Flash Constants
func TestFlashConstants_FLASH_ERROR(t *testing.T) {
	if FLASH_ERROR != "error" {
		t.Errorf("FLASH_ERROR = %v, want \"error\"", FLASH_ERROR)
	}
}

func TestFlashConstants_FLASH_SUCCESS(t *testing.T) {
	if FLASH_SUCCESS != "success" {
		t.Errorf("FLASH_SUCCESS = %v, want \"success\"", FLASH_SUCCESS)
	}
}

func TestFlashConstants_FLASH_INFO(t *testing.T) {
	if FLASH_INFO != "info" {
		t.Errorf("FLASH_INFO = %v, want \"info\"", FLASH_INFO)
	}
}

func TestFlashConstants_FLASH_WARNING(t *testing.T) {
	if FLASH_WARNING != "warning" {
		t.Errorf("FLASH_WARNING = %v, want \"warning\"", FLASH_WARNING)
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
