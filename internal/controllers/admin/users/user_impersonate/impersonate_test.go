package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	"project/internal/testutils"
	"testing"

	"github.com/dracory/userstore"
)

func TestImpersonate(t *testing.T) {
	// Setup
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)

	userID := "test_user"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	// Act
	err := Impersonate(registry.GetSessionStore(), w, req, userID)

	// Assert
	if err != nil {
		t.Fatalf("Impersonate failed: %v", err)
	}
}

func TestImpersonate_NilSessionStore(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	err := Impersonate(nil, w, req, "test_user")
	if err == nil {
		t.Error("Expected error for nil session store")
	}
	if err.Error() != "session store is nil" {
		t.Errorf("Expected 'session store is nil', got: %v", err)
	}
}

func TestNewUserImpersonateController(t *testing.T) {
	registry := testutils.Setup()
	controller := NewUserImpersonateController(registry)

	if controller == nil {
		t.Fatal("Controller should not be nil")
	}
	if controller.registry == nil {
		t.Fatal("Controller registry should not be nil")
	}
}

func TestUserImpersonateController_Handler_NoAuthUser(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	controller := NewUserImpersonateController(registry)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	result := controller.Handler(w, req)
	if result == "" {
		t.Error("Handler should return a result")
	}
}

func TestUserImpersonateController_Handler_WithAuthUser(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCacheStoreUsed(true)
	cfg.SetSessionStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	registry := testutils.Setup(testutils.WithCfg(cfg))

	// Create and set up a test user
	testUser, err := testutils.SeedUser(registry.GetUserStore(), "test_user_01")
	if err != nil {
		t.Fatal(err)
	}

	// Set user as administrator
	testUser.SetRole(userstore.USER_ROLE_ADMINISTRATOR)
	err = registry.GetUserStore().UserUpdate(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}

	// Create request with authenticated user
	req, _ := http.NewRequest("GET", "/?user_id=other_user", nil)
	ctx := context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, testUser)
	req = req.WithContext(ctx)

	controller := NewUserImpersonateController(registry)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)
	if result == "" {
		t.Error("Handler should return a result")
	}
}
