package swagger

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSwaggerUIController(t *testing.T) {
	req := httptest.NewRequest("GET", "/swagger", nil)
	w := httptest.NewRecorder()

	SwaggerUIController(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("SwaggerUIController() status = %d, want %d", w.Code, http.StatusOK)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", contentType, "text/html; charset=utf-8")
	}

	body := w.Body.String()
	if len(body) == 0 {
		t.Error("SwaggerUIController() returned empty body")
	}
	if !strings.Contains(body, "swagger") && !strings.Contains(body, "html") {
		t.Error("SwaggerUIController() body does not contain expected swagger UI content")
	}
}

func TestSwaggerYAMLController(t *testing.T) {
	req := httptest.NewRequest("GET", "/swagger.yaml", nil)
	w := httptest.NewRecorder()

	SwaggerYAMLController(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("SwaggerYAMLController() status = %d, want %d", w.Code, http.StatusOK)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/x-yaml; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/x-yaml; charset=utf-8")
	}

	body := w.Body.String()
	if len(body) == 0 {
		t.Error("SwaggerYAMLController() returned empty body")
	}
	if !strings.Contains(body, "openapi") && !strings.Contains(body, "swagger") {
		t.Error("SwaggerYAMLController() body does not contain expected swagger YAML content")
	}
}
