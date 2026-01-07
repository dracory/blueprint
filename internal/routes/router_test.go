package routes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/registry"
	"project/internal/routes"
	"project/internal/testutils"
	"project/internal/types"
)

func TestRoutes_HTTPWorkflows(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		method         string
		path           string
		setup          func() types.RegistryInterface
		expectedStatus int
		assert         func(*testing.T, *httptest.ResponseRecorder, types.RegistryInterface)
	}{
		{
			name:   "website home returns ok",
			method: http.MethodGet,
			path:   "/",
			setup: func() types.RegistryInterface {
				return testutils.Setup(
					testutils.WithCacheStore(true),
					testutils.WithSessionStore(true),
					testutils.WithUserStore(true),
				)
			},
			expectedStatus: http.StatusOK,
			assert: func(t *testing.T, rr *httptest.ResponseRecorder, _ types.RegistryInterface) {
				t.Helper()

				body := rr.Body.String()
				if !strings.Contains(body, "You are at the website home page") {
					t.Fatalf("expected response to contain website home copy, got %q", body)
				}
			},
		},
		{
			name:   "user home redirects unauthenticated users",
			method: http.MethodGet,
			path:   "/user",
			setup: func() types.RegistryInterface {
				return testutils.Setup(
					testutils.WithCacheStore(true),
					testutils.WithGeoStore(true),
					testutils.WithSessionStore(true),
					testutils.WithShopStore(true),
					testutils.WithUserStore(true),
				)
			},
			expectedStatus: http.StatusSeeOther,
			assert: func(t *testing.T, rr *httptest.ResponseRecorder, registry registry.RegistryInterface) {
				t.Helper()

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
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := tc.setup()
			router := routes.Router(app)

			req := httptest.NewRequest(tc.method, tc.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			tc.assert(t, rr, app)
		})
	}
}
