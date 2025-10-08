package seo

import (
    "net/http"
    "strings"
    "testing"

    "github.com/dracory/test"
)

func TestRobotsTxtController_Handler(t *testing.T) {
    controller := NewRobotsTxtController()

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

    expectedLines := []string{
        "User-agent: *",
        "Allow: /",
        "Disallow: /admin/",
        "Sitemap: /sitemap.xml",
    }

    for _, expected := range expectedLines {
        if !strings.Contains(body, expected) {
            t.Fatalf("expected body to contain %s, got: %s", expected, body)
        }
    }
}
