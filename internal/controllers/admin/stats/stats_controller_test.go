package stats

import (
	"context"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/geostore"
)

// countingGeoStore wraps the geo store and counts CountryList calls,
// to prove the stats controller fetches all countries once per request.
type countingGeoStore struct {
	geostore.StoreInterface
	calls *atomic.Int64
}

func (s *countingGeoStore) CountryList(ctx context.Context, options geostore.CountryQueryOptions) ([]geostore.Country, error) {
	s.calls.Add(1)
	us := geostore.NewCountry().SetName("United States")
	us.SetIsoCode2("US")
	de := geostore.NewCountry().SetName("Germany")
	de.SetIsoCode2("DE")
	return []geostore.Country{*us, *de}, nil
}

func TestCountryNameByIso2_LoadsCountriesOnce(t *testing.T) {
	app := testutils.Setup()
	controller := NewStatsController(app)

	var calls atomic.Int64
	app.SetGeoStore(&countingGeoStore{calls: &calls})

	req := httptest.NewRequest("POST", "/admin/stats?controller=sessions&action=load-sessions", nil)
	resolve := controller.countryNameByIso2(req)

	// The sessions page calls this once per unique visitor IP — thousands
	// of calls. A single CountryList must serve all of them.
	const lookups = 100
	for i := 0; i < lookups; i++ {
		code := "US"
		if i%2 == 0 {
			code = "DE"
		}
		if _, err := resolve(code); err != nil {
			t.Fatalf("resolve(%s) failed: %v", code, err)
		}
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("expected 1 geostore call for %d lookups, got %d", lookups, got)
	}
}

func TestCountryNameByIso2_CacheScopedToRequest(t *testing.T) {
	app := testutils.Setup()
	controller := NewStatsController(app)

	var calls atomic.Int64
	app.SetGeoStore(&countingGeoStore{calls: &calls})

	req1 := httptest.NewRequest("POST", "/admin/stats", nil)
	req2 := httptest.NewRequest("POST", "/admin/stats", nil)

	if _, err := controller.countryNameByIso2(req1)("US"); err != nil {
		t.Fatal(err)
	}
	if _, err := controller.countryNameByIso2(req2)("US"); err != nil {
		t.Fatal(err)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("cache must not leak across requests: expected 2 calls, got %d", got)
	}
}

func TestNewStatsController(t *testing.T) {
	app := testutils.Setup()
	controller := NewStatsController(app)

	if controller == nil {
		t.Error("NewStatsController() should not return nil")
	}
	if controller.app != app {
		t.Error("Controller app should match the provided app")
	}
	if controller.logger == nil {
		t.Error("Controller logger should not be nil")
	}
}

func TestStatsController_Handler_DefaultAction(t *testing.T) {
	app := testutils.Setup()
	controller := NewStatsController(app)

	req := httptest.NewRequest("GET", "/admin/stats", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestStatsController_Handler_WithDifferentMethods_GET(t *testing.T) {
	app := testutils.Setup()
	controller := NewStatsController(app)

	req := httptest.NewRequest("GET", "/admin/stats", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestStatsController_Handler_WithDifferentMethods_POST(t *testing.T) {
	app := testutils.Setup()
	controller := NewStatsController(app)

	req := httptest.NewRequest("POST", "/admin/stats", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestStatsController_RegistryField(t *testing.T) {
	app := testutils.Setup()
	controller := NewStatsController(app)

	if controller.app != app {
		t.Error("Controller app should match the provided app")
	}
}

func TestStatsController_NilRegistry(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic with nil app
		}
	}()

	controller := NewStatsController(nil)
	if controller == nil {
		t.Error("NewStatsController(nil) should not return nil")
	}
}

func TestStatsController_MultipleInstances(t *testing.T) {
	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	controller1 := NewStatsController(registry1)
	controller2 := NewStatsController(registry2)

	if controller1 == nil || controller2 == nil {
		t.Fatal("All controllers should be non-nil")
	}

	if controller1 == controller2 {
		t.Error("Controllers should be separate instances")
	}

	if controller1.app != registry1 {
		t.Error("Controller1 should have registry1")
	}

	if controller2.app != registry2 {
		t.Error("Controller2 should have registry2")
	}
}
