package shared

import (
	"log/slog"
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

// TestHeaderNilStore verifies Header handles nil store
func TestHeaderNilStore(t *testing.T) {
	t.Parallel()
	logger := slog.Default()
	req := httptest.NewRequest("GET", "/", nil)

	result := Header(nil, logger, req)

	// Should return nil when store is nil
	if result != nil {
		t.Error("Header(nil) should return nil")
	}
}

// TestHeaderNotNil verifies Header returns non-nil result with valid inputs
func TestHeaderNotNil(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	logger := slog.Default()
	req := httptest.NewRequest("GET", "/", nil)

	result := Header(app.GetShopStore(), logger, req)

	if result == nil {
		t.Error("Header() should return non-nil result with valid store")
	}
}
