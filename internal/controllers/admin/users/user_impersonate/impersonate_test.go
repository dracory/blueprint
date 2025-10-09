package admin

import (
	"net/http"
	"net/http/httptest"
	"project/internal/testutils"
	"testing"
)

func TestImpersonate(t *testing.T) {
	// Setup
	_, sessionStore, _, cleanup := testutils.SetupTestAuth(t)
	defer cleanup()
	userID := "test_user"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	// Act
	err := Impersonate(sessionStore, w, req, userID)

	// Assert
	if err != nil {
		t.Fatalf("Impersonate failed: %v", err)
	}
}
