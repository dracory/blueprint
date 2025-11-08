package admin

import (
	"net/http"
	"net/http/httptest"
	"project/internal/testutils"
	"testing"
)

func TestImpersonate(t *testing.T) {
	// Setup
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	userID := "test_user"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	// Act
	err := Impersonate(app.GetSessionStore(), w, req, userID)

	// Assert
	if err != nil {
		t.Fatalf("Impersonate failed: %v", err)
	}
}
