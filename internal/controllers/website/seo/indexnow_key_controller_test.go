package seo

import (
	"net/http"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
)

func TestIndexNowKeyController_Handler(t *testing.T) {
	// Setup test app with config
	app := testutils.Setup()
	app.GetConfig().SetIndexNowKey("cd325dd195454606a8316fb303224f37")

	controller := NewIndexNowKeyController(app)

	body, response, err := test.CallStringEndpoint(http.MethodGet, controller.Handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if got := response.Header.Get("Content-Type"); got != "text/plain" {
		t.Fatalf("expected Content-Type text/plain, got %s", got)
	}

	expected := "cd325dd195454606a8316fb303224f37"
	if body != expected {
		t.Fatalf("expected body %s, got %s", expected, body)
	}
}
