package routes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/routes"
	"project/internal/testutils"
)

func TestRoutes_HTTPWorkflows_WebsiteHomeReturnsOk(t *testing.T) {
	t.Parallel()

	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSessionStore(true),
		testutils.WithUserStore(true),
	)
	router := routes.Router(registry)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Fatalf("expected response to contain html doctype, got %q", body)
	}
}

func TestRoutes_HTTPWorkflows_UserHomeRedirectsUnauthenticatedUsers(t *testing.T) {
	t.Parallel()

	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithShopStore(true),
		testutils.WithUserStore(true),
	)
	router := routes.Router(registry)

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	if rr.Header().Get("Location") == "" {
		t.Fatalf("expected redirect Location header to be set")
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(registry.GetCacheStore(), rr.Result())
	if err != nil {
		t.Fatalf("failed to fetch flash message: %v", err)
	}

	if flashMessage == nil {
		t.Fatalf("expected flash message to be stored for redirect")
	}

	if flashMessage.Type != "error" {
		t.Fatalf("expected flash message type error, got %s", flashMessage.Type)
	}

	if flashMessage.Message != "Only authenticated users can access this page" {
		t.Fatalf("unexpected flash message: %s", flashMessage.Message)
	}
}
