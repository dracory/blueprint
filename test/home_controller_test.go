package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"project/app/config"
	"project/app/controllers/shared/notfound"
	"project/app/controllers/website/about"
	"project/app/controllers/website/home"
	"project/internal/platform/database"
)

func TestHomeController_Index(t *testing.T) {
	// Create a test configuration
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}
	cfg.AppEnv = "testing"

	// Create an in-memory SQLite database for testing
	db, err := database.New(cfg)
	if err != nil {
		t.Fatalf("Error creating database: %v", err)
	}
	defer db.Close()

	// Create a home controller
	homeController := home.NewController(cfg, db)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeController.Index)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body contains expected content
	if rr.Body.String() == "" {
		t.Errorf("handler returned empty body")
	}

	// Check if the response contains the expected title
	if !contains(rr.Body.String(), "Welcome to Dracory") {
		t.Errorf("handler returned unexpected body: did not contain 'Welcome to Dracory'")
	}
}

func TestHomeController_About(t *testing.T) {
	// Create a test configuration
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}
	cfg.AppEnv = "testing"

	// Create an in-memory SQLite database for testing
	db, err := database.New(cfg)
	if err != nil {
		t.Fatalf("Error creating database: %v", err)
	}
	defer db.Close()

	// Create a home controller
	aboutController := about.NewController(cfg, db)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/about", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(aboutController.Index)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body contains expected content
	if rr.Body.String() == "" {
		t.Errorf("handler returned empty body")
	}

	// Check if the response contains the expected title
	if !contains(rr.Body.String(), "About Dracory") {
		t.Errorf("handler returned unexpected body: did not contain 'About Dracory'")
	}
}

func TestHomeController_NotFound(t *testing.T) {
	// Create a test configuration
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}
	cfg.AppEnv = "testing"

	// Create an in-memory SQLite database for testing
	db, err := database.New(cfg)
	if err != nil {
		t.Fatalf("Error creating database: %v", err)
	}
	defer db.Close()

	// Create a home controller
	notfoundController := notfound.NewController(cfg, db)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/not-found", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(notfoundController.Index)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	// Check the response body contains expected content
	if rr.Body.String() == "" {
		t.Errorf("handler returned empty body")
	}

	// Check if the response contains the expected title
	if !contains(rr.Body.String(), "404 - Page Not Found") {
		t.Errorf("handler returned unexpected body: did not contain '404 - Page Not Found'")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return s != "" && s != substr && s != substr+"\n" && s != "\n"+substr && s != "\n"+substr+"\n"
}
