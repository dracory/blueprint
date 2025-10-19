package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	userDir "project/internal/controllers/user"
	"project/internal/links"
	"project/internal/testutils"
	"strings"
	"testing"

	"github.com/dracory/rtr"
)

// TestUserRoutesCount verifies the number of user routes registered.
func TestUserRoutesCount(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithShopStore(true),
		testutils.WithUserStore(true),
	)
	routes := userDir.Routes(app)

	// We expect at least the core user routes
	// (order routes + apikey routes + profile + home + homeCatchAll)
	if len(routes) < 3 {
		t.Fatalf("expected at least 3 user routes, got %d", len(routes))
	}
}

// TestUserRoutesNotNil ensures no route entries are nil.
func TestUserRoutesNotNil(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithShopStore(true),
		testutils.WithUserStore(true),
	)
	routes := userDir.Routes(app)

	for i, rt := range routes {
		if rt == nil {
			t.Fatalf("route at index %d is nil", i)
		}
	}
}

// TestUserRoutesAreAdded verifies that expected user routes are present.
func TestUserRoutesAreAdded(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithShopStore(true),
		testutils.WithUserStore(true),
	)
	routes := userDir.Routes(app)

	expectedPaths := []string{
		links.USER_HOME,
		links.USER_PROFILE,
		links.USER_HOME + links.CATCHALL, // catch-all route
	}

	for _, expectedPath := range expectedPaths {
		found := false
		for _, rt := range routes {
			if rt.GetPath() == expectedPath {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected route with path %s not found", expectedPath)
		}
	}
}

// TestUserCatchAllRouteIsLast verifies that the catch-all route is registered last.
// This is critical to prevent the catch-all from intercepting specific routes.
func TestUserCatchAllRouteIsLast(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithShopStore(true),
		testutils.WithUserStore(true),
	)
	routes := userDir.Routes(app)

	if len(routes) == 0 {
		t.Fatal("no routes found")
	}

	// The catch-all route should be the last route
	catchAllPath := links.USER_HOME + links.CATCHALL
	lastRoute := routes[len(routes)-1]

	if lastRoute.GetPath() != catchAllPath {
		t.Fatalf("expected catch-all route %s to be last, but got %s", catchAllPath, lastRoute.GetPath())
	}
}

// TestUserSpecificRoutesBeforeCatchAll verifies that specific routes come before the catch-all.
// This ensures that specific routes like /user/orders/create/contract-details
// are matched before the catch-all /user/* route.
func TestUserSpecificRoutesBeforeCatchAll(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithShopStore(true),
		testutils.WithUserStore(true),
	)
	routes := userDir.Routes(app)

	catchAllPath := links.USER_HOME + links.CATCHALL
	catchAllIndex := -1

	// Find the index of the catch-all route
	for i, rt := range routes {
		if rt.GetPath() == catchAllPath {
			catchAllIndex = i
			break
		}
	}

	if catchAllIndex == -1 {
		t.Fatal("catch-all route not found")
	}

	// Verify that specific routes come before the catch-all
	specificPaths := []string{
		// links.USER_ORDER_CREATE_CONTRACT_DETAILS,
		// links.USER_ORDER_CREATE_CONTRACT_UPLOAD,
		// links.USER_ORDER_CREATE_SERVICE_SELECT,
		links.USER_PROFILE,
	}

	for _, specificPath := range specificPaths {
		found := false
		for i, rt := range routes {
			if rt.GetPath() == specificPath {
				found = true
				if i >= catchAllIndex {
					t.Fatalf("specific route %s (index %d) comes after catch-all route (index %d)", specificPath, i, catchAllIndex)
				}
				break
			}
		}
		if !found {
			t.Fatalf("specific route %s not found", specificPath)
		}
	}
}

// TestUserHomeRouteBeforeCatchAll verifies that the home route comes before the catch-all.
func TestUserHomeRouteBeforeCatchAll(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithShopStore(true),
		testutils.WithUserStore(true),
	)
	routes := userDir.Routes(app)

	catchAllPath := links.USER_HOME + links.CATCHALL
	homePath := links.USER_HOME

	catchAllIndex := -1
	homeIndex := -1

	for i, rt := range routes {
		if rt.GetPath() == catchAllPath {
			catchAllIndex = i
		}
		if rt.GetPath() == homePath {
			homeIndex = i
		}
	}

	if catchAllIndex == -1 {
		t.Fatal("catch-all route not found")
	}

	if homeIndex == -1 {
		t.Fatal("home route not found")
	}

	if homeIndex >= catchAllIndex {
		t.Fatalf("home route (index %d) should come before catch-all route (index %d)", homeIndex, catchAllIndex)
	}
}

func TestUserHomePage_RedirectsNonLoggedUser(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithShopStore(true),
		testutils.WithUserStore(true),
	)

	// expected := `<title>Home | Client Dashboard | TinyFunnel`

	// Create a request to pass to the handler.
	req, err := http.NewRequest("GET", links.User().Home(), nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Create a ResponseRecorder so that we can inspect the response.
	rr := httptest.NewRecorder()

	// Call the ServeHTTP method directly and pass in the request and response recorder.
	r := rtr.NewRouter()
	r.AddRoutes(userDir.Routes(app))
	r.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}

	body := rr.Body.String()

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), rr.Result())

	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	if flashMessage.Type != "error" {
		t.Fatal(`Response be of type 'success', but got: `, flashMessage.Type, flashMessage.Message)
	}

	if flashMessage.Message != "Only authenticated users can access this page" {
		t.Fatal(`Response MUST contain 'Only authenticated users can access this page', but got: `, flashMessage.Message)
	}

	expecteds := []string{
		`<a href="/flash?message_id=`,
		`">See Other</a>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, body)
		}
	}
}

// func TestUserHomePage_RedirectsWhenSubscriptionRequired(t *testing.T) {
// 	app := testutils.Setup(
// 		testutils.WithCacheStore(true),
// 		testutils.WithGeoStore(true),
// 		testutils.WithSessionStore(true),
// 		testutils.WithShopStore(true),
// 		testutils.WithUserStore(true),
// 	)

// 	expectedRedirect := "/user/subscriptions/plan-select"

// 	user, session, err := testutils.SeedUserAndSession(app.GetUserStore(), app.GetSessionStore(), testutils.USER_01, httptest.NewRequest("GET", "/", nil), 1)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Populate user
// 	user.SetFirstName("Test")
// 	user.SetLastName("User")
// 	user.SetEmail("test@example.com")
// 	err = app.GetUserStore().UserUpdate(context.Background(), user)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	req, err := testutils.NewRequest(http.MethodGet, links.User().Home(), testutils.NewRequestOptions{
// 		Context: map[any]any{
// 			config.AuthenticatedUserContextKey{}:    user,
// 			config.AuthenticatedSessionContextKey{}: session,
// 		},
// 	})

// 	if err != nil {
// 		t.Fatalf("could not create request: %v", err)
// 	}

// 	// Create a ResponseRecorder so that we can inspect the response.
// 	rr := httptest.NewRecorder()

// 	// Call the ServeHTTP method directly and pass in the request and response recorder.
// 	r := rtr.NewRouter()
// 	r.AddRoutes(userDir.Routes(app))
// 	r.ServeHTTP(rr, req)

// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusTemporaryRedirect {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusTemporaryRedirect)
// 	}

// 	// Check the redirect location is what we expect.
// 	location := rr.Header().Get("Location")
// 	if !strings.Contains(location, expectedRedirect) {
// 		t.Errorf("handler returned unexpected redirect location: got %v want %v", location, expectedRedirect)
// 	}
// }

func TestUserHomePage(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithSessionStore(true),
		testutils.WithShopStore(true),
		testutils.WithUserStore(true),
	)

	expected := `<title>Home | User</title>`

	user, session, err := testutils.SeedUserAndSession(app.GetUserStore(), app.GetSessionStore(), testutils.USER_01, httptest.NewRequest("GET", "/", nil), 1)

	if err != nil {
		t.Fatal(err)
	}

	user.SetFirstName("Test")
	user.SetLastName("User")
	user.SetEmail("test@example.com")
	err = app.GetUserStore().UserUpdate(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}

	// _, err = testutils.SeedSubscription(app.GetSubscriptionStore(), user.ID(), testutils.PLAN_01)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	req, err := testutils.NewRequest(http.MethodGet, links.User().Home(), testutils.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}:    user,
			config.AuthenticatedSessionContextKey{}: session,
		},
	})

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Create a ResponseRecorder so that we can inspect the response.
	rr := httptest.NewRecorder()

	// Call the ServeHTTP method directly and pass in the request and response recorder.
	r := rtr.NewRouter()
	r.AddRoutes(userDir.Routes(app))
	r.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	body := rr.Body.String()
	if !strings.Contains(body, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", body, expected)
	}
}
