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
	registry := testutils.Setup()

	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	registry.SetLogger(logger)

	// Act
	handler := LogRequestMiddleware(registry).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	registry := testutils.Setup()

	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	registry.SetLogger(logger)

	buf.Reset() // Clear buffer for each test

	// Act
	handler := LogRequestMiddleware(registry).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := testutils.NewRequest("GET", "/style.css", testutils.NewRequestOptions{})
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
		t.Errorf("Expected empty log for /style.css, got %s", logOutput)
	}
}

func TestLogRequestMiddleware_Filtered_JS(t *testing.T) {
	// Arrange
	registry := testutils.Setup()

	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	registry.SetLogger(logger)

	buf.Reset() // Clear buffer for each test

	// Act
	handler := LogRequestMiddleware(registry).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := testutils.NewRequest("GET", "/script.js", testutils.NewRequestOptions{})
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
		t.Errorf("Expected empty log for /script.js, got %s", logOutput)
	}
}

func TestLogRequestMiddleware_Filtered_Image(t *testing.T) {
	// Arrange
	registry := testutils.Setup()

	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	registry.SetLogger(logger)

	buf.Reset() // Clear buffer for each test

	// Act
	handler := LogRequestMiddleware(registry).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := testutils.NewRequest("GET", "/image.png", testutils.NewRequestOptions{})
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
		t.Errorf("Expected empty log for /image.png, got %s", logOutput)
	}
}

func TestLogRequestMiddleware_Filtered_Favicon(t *testing.T) {
	// Arrange
	registry := testutils.Setup()

	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	registry.SetLogger(logger)

	buf.Reset() // Clear buffer for each test

	// Act
	handler := LogRequestMiddleware(registry).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := testutils.NewRequest("GET", "/favicon.ico", testutils.NewRequestOptions{})
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
		t.Errorf("Expected empty log for /favicon.ico, got %s", logOutput)
	}
}

func TestLogRequestMiddleware_Filtered_Health(t *testing.T) {
	// Arrange
	registry := testutils.Setup()

	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	registry.SetLogger(logger)

	buf.Reset() // Clear buffer for each test

	// Act
	handler := LogRequestMiddleware(registry).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := testutils.NewRequest("GET", "/health", testutils.NewRequestOptions{})
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
		t.Errorf("Expected empty log for /health, got %s", logOutput)
	}
}

func TestLogRequestMiddleware_Filtered_Ping(t *testing.T) {
	// Arrange
	registry := testutils.Setup()

	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	registry.SetLogger(logger)

	buf.Reset() // Clear buffer for each test

	// Act
	handler := LogRequestMiddleware(registry).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := testutils.NewRequest("GET", "/ping", testutils.NewRequestOptions{})
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
		t.Errorf("Expected empty log for /ping, got %s", logOutput)
	}
}

func TestLogRequestMiddleware_Filtered_AssetsFolder(t *testing.T) {
	// Arrange
	registry := testutils.Setup()

	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	registry.SetLogger(logger)

	buf.Reset() // Clear buffer for each test

	// Act
	handler := LogRequestMiddleware(registry).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := testutils.NewRequest("GET", "/assets/image.jpg", testutils.NewRequestOptions{})
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
		t.Errorf("Expected empty log for /assets/image.jpg, got %s", logOutput)
	}
}

func TestLogRequestMiddleware_Filtered_StaticFolder(t *testing.T) {
	// Arrange
	registry := testutils.Setup()

	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	registry.SetLogger(logger)

	buf.Reset() // Clear buffer for each test

	// Act
	handler := LogRequestMiddleware(registry).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := testutils.NewRequest("GET", "/static/style.css", testutils.NewRequestOptions{})
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
		t.Errorf("Expected empty log for /static/style.css, got %s", logOutput)
	}
}
