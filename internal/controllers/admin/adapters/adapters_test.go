package adapters

import (
	"context"
	"net/http/httptest"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
)

func TestUnionIDSets(t *testing.T) {
	got := unionIDSets([]string{"a", "b"}, []string{"b", "c"})
	if len(got) != 3 {
		t.Errorf("unionIDSets = %v, want 3 unique ids", got)
	}
}

func TestIntersectIDSets(t *testing.T) {
	got := intersectIDSets([][]string{{"a", "b", "c"}, {"b", "c"}, {"c", "b", "d"}})
	if len(got) != 2 {
		t.Errorf("intersectIDSets = %v, want [b c]", got)
	}

	if got := intersectIDSets(nil); got != nil {
		t.Errorf("intersectIDSets(nil) = %v, want nil", got)
	}

	if got := intersectIDSets([][]string{{"a"}, {}}); len(got) != 0 {
		t.Errorf("intersect with empty set = %v, want empty", got)
	}
}

func TestUserStoreCustomerResolverNilStore(t *testing.T) {
	r := NewUserStoreCustomerResolver(nil)
	name, email := r.FindByID(context.Background(), "any")
	if name != "" || email != "" {
		t.Errorf("expected empty result, got %q %q", name, email)
	}
	ids, err := r.SearchIDs(context.Background(), "x", "y")
	if ids != nil || err != nil {
		t.Errorf("expected nil ids, got %v, err %v", ids, err)
	}
}

func TestUserStoreCustomerResolver(t *testing.T) {
	app := testutils.Setup(testutils.WithUserStore(true))
	defer app.Close()

	ctx := context.Background()
	user := userstore.NewUser().
		SetEmail("jane@example.com").
		SetFirstName("Jane").
		SetLastName("Doe").
		SetStatus(userstore.USER_STATUS_ACTIVE)
	if err := app.GetUserStore().UserCreate(ctx, user); err != nil {
		t.Fatalf("UserCreate failed: %v", err)
	}

	r := NewUserStoreCustomerResolver(app.GetUserStore())

	name, email := r.FindByID(ctx, user.GetID())
	if name != "Jane Doe" || email != "jane@example.com" {
		t.Errorf("FindByID = %q %q, want Jane Doe jane@example.com", name, email)
	}

	name, email = r.FindByID(ctx, "missing")
	if name != "" || email != "" {
		t.Errorf("FindByID(missing) = %q %q, want empty", name, email)
	}

	ids, err := r.SearchIDs(ctx, "", "jane@")
	if err != nil || len(ids) != 1 || ids[0] != user.GetID() {
		t.Errorf("SearchIDs email = %v, err %v", ids, err)
	}

	ids, err = r.SearchIDs(ctx, "Jane", "")
	if err != nil || len(ids) != 1 {
		t.Errorf("SearchIDs name = %v, err %v", ids, err)
	}
}

func TestPiiSealUnsealPassthroughWithoutVault(t *testing.T) {
	app := testutils.Setup(testutils.WithUserStore(true))
	defer app.Close()

	ctx := context.Background()
	user := userstore.NewUser().
		SetEmail("a@b.com").
		SetFirstName("A").
		SetLastName("B")

	seal := NewUserPiiSealFunc(app)
	sealed, err := seal(ctx, user)
	if err != nil || sealed.GetEmail() != "a@b.com" {
		t.Errorf("seal should passthrough, got %v err %v", sealed, err)
	}

	unseal := NewUserPiiUnsealFunc(app)
	unsealed, err := unseal(ctx, user)
	if err != nil || unsealed.GetEmail() != "a@b.com" {
		t.Errorf("unseal should passthrough, got %v err %v", unsealed, err)
	}

	nilUser, err := unseal(ctx, nil)
	if nilUser != nil || err != nil {
		t.Errorf("unseal(nil) = %v err %v, want nil nil", nilUser, err)
	}

	batch := NewUsersPiiUnsealFunc(app)
	users, err := batch(ctx, []userstore.UserInterface{user})
	if err != nil || len(users) != 1 {
		t.Errorf("batch unseal = %v err %v", users, err)
	}
}

func TestOnUserSearchFuncEmpty(t *testing.T) {
	fn := NewOnUserSearchFunc(nil)
	ids, err := fn(context.Background(), nil)
	if ids != nil || err != nil {
		t.Errorf("expected nil ids for nil app, got %v err %v", ids, err)
	}
}

func TestOnUserImpersonateFunc(t *testing.T) {
	app := testutils.Setup(testutils.WithSessionStore(true))
	defer app.Close()

	fn := NewOnUserImpersonateFunc(app)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	if err := fn(rec, req, "user-123"); err != nil {
		t.Fatalf("impersonate failed: %v", err)
	}

	sessions, err := app.GetSessionStore().SessionList(context.Background(),
		sessionstore.NewSessionQuery().SetUserID("user-123"))
	if err != nil || len(sessions) == 0 {
		t.Fatalf("expected impersonation session, got %v err %v", sessions, err)
	}

	if len(rec.Result().Cookies()) == 0 {
		t.Error("expected auth cookie to be set")
	}
}

func TestOnUserImpersonateFuncNoSessionStore(t *testing.T) {
	app := testutils.Setup() // session store disabled
	defer app.Close()

	fn := NewOnUserImpersonateFunc(app)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	if err := fn(rec, req, "user-123"); err == nil {
		t.Fatal("expected error when session store is disabled")
	}
}

func TestOnUserUpdatedFuncNoTaskStore(t *testing.T) {
	app := testutils.Setup() // task store disabled
	defer app.Close()

	fn := NewOnUserUpdatedFunc(app, "some-alias")
	fn(context.Background(), "user-123") // must not panic
}

func TestGeoResolverNilStore(t *testing.T) {
	r := NewGeoResolver(nil)
	countries, err := r.Countries(context.Background())
	if countries != nil || err != nil {
		t.Errorf("expected nil, got %v err %v", countries, err)
	}
	tzs, err := r.Timezones(context.Background(), "US")
	if tzs != nil || err != nil {
		t.Errorf("expected nil, got %v err %v", tzs, err)
	}
}

func TestGeoResolver(t *testing.T) {
	app := testutils.Setup(testutils.WithGeoStore(true))
	defer app.Close()

	r := NewGeoResolver(app.GetGeoStore())

	if _, err := r.Countries(context.Background()); err != nil {
		t.Fatalf("Countries failed: %v", err)
	}

	tzs, err := r.Timezones(context.Background())
	if tzs != nil || err != nil {
		t.Errorf("expected nil without country code, got %v err %v", tzs, err)
	}
}

func TestLlmFactoryTestingProvider(t *testing.T) {
	app := testutils.Setup() // DefaultConf sets env=testing
	defer app.Close()

	factory := NewLlmFactory(app)
	llmInstance, err := factory()
	if err != nil || llmInstance == nil {
		t.Errorf("expected mock LLM in testing env, got %v err %v", llmInstance, err)
	}
}

func TestFlashRedirectFunc(t *testing.T) {
	app := testutils.Setup(testutils.WithCacheStore(true))
	defer app.Close()

	fn := NewFlashRedirectFunc(app)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	out := fn(rec, req, "success", "Done!", "/next", 0)
	if out == "" {
		t.Error("expected non-empty redirect response")
	}
}
