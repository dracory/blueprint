package home

import (
	"net/http"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
)

func TestHomeController_Handler(t *testing.T) {
	// Setup
	app := testutils.Setup()
	app.GetConfig().SetAppName("TEST APP NAME")

	// Execute
	body, response, err := test.CallStringEndpoint(http.MethodGet, NewHomeController(app).Handler, test.NewRequestOptions{})

	// Assert
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`Welcome to TEST APP NAME`,
		`<!DOCTYPE html>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, body)
		}
	}
}
