package middlewares

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"project/internal/testutils"
	"strings"
	"testing"
)

func TestLogRequestMiddleware(t *testing.T) {
	// Arrange
	app := testutils.Setup()

	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	app.SetLogger(logger)

	// Act
	handler := LogRequestMiddleware(app).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := testutils.NewRequest("GET", "/test-path", testutils.NewRequestOptions{})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "GET request") {
		t.Errorf("Expected log to contain 'request', got %s", logOutput)
	}
	if !strings.Contains(logOutput, "test-path") {
		t.Errorf("Expected log to contain 'test-path', got %s", logOutput)
	}
	if !strings.Contains(logOutput, "GET") {
		t.Errorf("Expected log to contain 'GET', got %s", logOutput)
	}
}

func TestLogRequestMiddleware_Filtered(t *testing.T) {
	// Arrange
	app := testutils.Setup()

	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	app.SetLogger(logger)

	tests := []struct {
		name string
		path string
	}{
		{"Static Asset CSS", "/style.css"},
		{"Static Asset JS", "/script.js"},
		{"Static Asset Image", "/image.png"},
		{"Static Asset ICO", "/favicon.ico"},
		{"Health Check", "/health"},
		{"Ping", "/ping"},
		{"Assets Folder", "/assets/image.jpg"},
		{"Static Folder", "/static/style.css"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset() // Clear buffer for each test

			// Act
			handler := LogRequestMiddleware(app).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req, err := testutils.NewRequest("GET", tt.path, testutils.NewRequestOptions{})
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Assert
			if w.Code != http.StatusOK {
				t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}

			logOutput := buf.String()
			if logOutput != "" {
				t.Errorf("Expected empty log for %s, got %s", tt.path, logOutput)
			}
		})
	}
}
