package seo

import (
    "net/http"
    "strings"
    "testing"

    "github.com/dracory/test"
)

func TestSecurityTxtController_Handler(t *testing.T) {
    controller := NewSecurityTxtController()

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

    if !strings.Contains(body, "\r\n") {
        t.Fatalf("expected CRLF line endings, got body: %s", body)
    }

    if strings.Contains(strings.ReplaceAll(body, "\r\n", ""), "\n") {
        t.Fatalf("found bare LF characters, body: %s", body)
    }

    expected := "Contact: https://tiny.vip/BlCe"
    if !strings.Contains(body, expected) {
        t.Fatalf("expected body to contain %s, got: %s", expected, body)
    }
}
