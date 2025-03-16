package website

import (
	"net/http"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/base/test"
)

func TestHomeController_Handler(t *testing.T) {
	// Setup
	testutils.Setup()

	// Execute
	body, response, err := test.CallStringEndpoint(http.MethodPost, newHomeController().Handler, test.NewRequestOptions{})

	// Assert
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`You are at the website home page`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, body)
		}
	}
}
