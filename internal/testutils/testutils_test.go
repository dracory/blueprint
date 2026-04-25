package testutils

import (
	"bytes"
	"net/url"
	"project/internal/config"
	"testing"

	"github.com/dracory/test"
)

func TestTestKeyIntegration(t *testing.T) {
	// Test that the TestKey function from the blueprint project
	// produces the same result as the TestKey function from the base project
	cfg := config.New()
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseHost("")
	cfg.SetDatabasePort("")
	cfg.SetDatabaseName("file::memory:?cache=shared")
	cfg.SetDatabaseUsername("")
	cfg.SetDatabasePassword("")
	blueprintKey := TestKey(cfg)
	baseKey := test.TestKey(cfg.GetDatabaseDriver(), cfg.GetDatabaseHost(), cfg.GetDatabasePort(), cfg.GetDatabaseName(), cfg.GetDatabaseUsername(), cfg.GetDatabasePassword())

	if blueprintKey != baseKey {
		t.Errorf("Blueprint TestKey and base TestKey should produce the same result")
	}
}

func TestTestConfigIntegration(t *testing.T) {
	// Create a test configuration
	testConfig := config.New()
	// App
	testConfig.SetAppHost("localhost")
	testConfig.SetAppPort("8080")
	testConfig.SetAppName("Test App")
	testConfig.SetAppUrl("http://localhost:8080")
	testConfig.SetAppEnv(config.APP_ENVIRONMENT_TESTING)

	// Database
	testConfig.SetDatabaseDriver("sqlite")
	testConfig.SetDatabaseHost("")
	testConfig.SetDatabasePort("")
	testConfig.SetDatabaseName("file::memory:?cache=shared")
	testConfig.SetDatabaseUsername("")
	testConfig.SetDatabasePassword("")
	testConfig.SetAppDebug(true)

	// Encryption
	testConfig.SetEnvEncryptionKey("123456")

	// Mail
	testConfig.SetMailDriver("smtp")
	testConfig.SetMailHost("127.0.0.1")
	testConfig.SetMailPort(32435)
	testConfig.SetMailUsername("")
	testConfig.SetMailPassword("")
	testConfig.SetMailFromAddress("admintest@test.com")
	testConfig.SetMailFromName("Admin Test")

	// Stores
	testConfig.SetCmsStoreUsed(false)
	testConfig.SetCmsStoreTemplateID("default")
	testConfig.SetVaultStoreUsed(false)
	testConfig.SetVaultStoreKey("abcdefghijklmnopqrstuvwxyz1234567890")

	// Artificial Intelligence LLMs
	testConfig.SetAnthropicApiUsed(false)
	testConfig.SetAnthropicApiKey("anthropic_api_key")
	testConfig.SetAnthropicApiDefaultModel("anthropic_api_default_model")
	testConfig.SetGoogleGeminiApiUsed(false)
	testConfig.SetGoogleGeminiApiKey("google_gemini_api_key")
	testConfig.SetGoogleGeminiApiDefaultModel("google_gemini_api_default_model")
	testConfig.SetOpenAiApiUsed(false)
	testConfig.SetOpenAiApiKey("openai_api_key")
	testConfig.SetOpenRouterApiUsed(false)
	testConfig.SetOpenRouterApiKey("openrouter_api_key")
	testConfig.SetOpenRouterApiDefaultModel("openrouter_api_default_model")
	testConfig.SetVertexAiApiUsed(false)
	testConfig.SetVertexAiApiDefaultModel("vertex_ai_api_default_model")
	testConfig.SetVertexAiApiProjectID("vertex_ai_api_project_id")
	testConfig.SetVertexAiApiRegionID("vertex_ai_api_region_id")
	testConfig.SetVertexAiApiModelID("vertex_ai_api_model_id")

	// Payments
	testConfig.SetStripeKeyPrivate("sk_test_yoursecretkey")
	testConfig.SetStripeKeyPublic("pk_test_yourpublickey")

	// Verify that the configuration was applied
	if testConfig.GetAppName() != "Test App" {
		t.Errorf("Expected AppName to be %s, got %s", "Test App", testConfig.GetAppName())
	}
}

func TestTestDBIntegration(t *testing.T) {
	// Create a test database
	db, err := test.NewTestDB(nil)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer func() {
		_ = test.CloseTestDB(db)
	}()

	// Create a test table
	err = test.CreateTestTable(db, "test_users", "id INTEGER PRIMARY KEY, name TEXT")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Insert test data
	err = test.ExecuteSQLWithArgs(db, "INSERT INTO test_users (name) VALUES (?)", "Test User")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Query the test data
	var name string
	err = db.QueryRow("SELECT name FROM test_users WHERE name = ?", "Test User").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to query test data: %v", err)
	}

	// Verify the result
	if name != "Test User" {
		t.Errorf("Expected name to be 'Test User', got '%s'", name)
	}
}

func TestSeedUser_NilStore(t *testing.T) {
	t.Parallel()
	user, err := SeedUser(nil, "test-id")
	if err == nil {
		t.Error("Expected error when userStore is nil")
	}
	if user != nil {
		t.Error("Expected nil user when userStore is nil")
	}
}

func TestSeedUser_EmptyUserID(t *testing.T) {
	t.Parallel()
	app := Setup(WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	user, err := SeedUser(app.GetUserStore(), "")
	if err == nil {
		t.Error("Expected error when userID is empty")
	}
	if user != nil {
		t.Error("Expected nil user when userID is empty")
	}
}

func TestSeedUser_NewUser(t *testing.T) {
	t.Parallel()
	app := Setup(WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	user, err := SeedUser(app.GetUserStore(), "new-user-id")
	if err != nil {
		t.Fatalf("Failed to seed user: %v", err)
	}
	if user == nil {
		t.Fatal("Expected non-nil user")
	}
	if user.GetID() != "new-user-id" {
		t.Errorf("Expected user ID to be 'new-user-id', got '%s'", user.GetID())
	}
	if user.GetStatus() != "active" {
		t.Errorf("Expected user status to be 'active', got '%s'", user.GetStatus())
	}
}

func TestSeedUser_ExistingUser(t *testing.T) {
	t.Parallel()
	app := Setup(WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	// Create user first
	user1, err := SeedUser(app.GetUserStore(), "existing-user-id")
	if err != nil {
		t.Fatalf("Failed to seed user: %v", err)
	}

	// Try to seed same user again
	user2, err := SeedUser(app.GetUserStore(), "existing-user-id")
	if err != nil {
		t.Fatalf("Failed to seed existing user: %v", err)
	}
	if user2.GetID() != user1.GetID() {
		t.Error("Expected same user ID")
	}
}

func TestSeedUser_UserRole(t *testing.T) {
	t.Parallel()
	app := Setup(WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	// Test USER_01 role
	user1, err := SeedUser(app.GetUserStore(), test.USER_01)
	if err != nil {
		t.Fatalf("Failed to seed user: %v", err)
	}
	if user1.GetRole() != "user" {
		t.Errorf("Expected USER_01 to have role 'user', got '%s'", user1.GetRole())
	}

	// Test ADMIN_01 role
	user2, err := SeedUser(app.GetUserStore(), test.ADMIN_01)
	if err != nil {
		t.Fatalf("Failed to seed user: %v", err)
	}
	if user2.GetRole() != "administrator" {
		t.Errorf("Expected ADMIN_01 to have role 'administrator', got '%s'", user2.GetRole())
	}
}

func TestSeedSession_NilStore(t *testing.T) {
	t.Parallel()
	app := Setup(WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	user, _ := SeedUser(app.GetUserStore(), "test-user")
	r, _ := NewRequest("GET", "/", NewRequestOptions{})

	session, err := SeedSession(nil, r, user, 10)
	if err == nil {
		t.Error("Expected error when sessionStore is nil")
	}
	if session != nil {
		t.Error("Expected nil session when sessionStore is nil")
	}
}

func TestSeedSession_Success(t *testing.T) {
	t.Parallel()
	app := Setup(WithUserStore(true), WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	user, _ := SeedUser(app.GetUserStore(), "test-user")
	r, _ := NewRequest("GET", "/", NewRequestOptions{})

	session, err := SeedSession(app.GetSessionStore(), r, user, 10)
	if err != nil {
		t.Fatalf("Failed to seed session: %v", err)
	}
	if session == nil {
		t.Fatal("Expected non-nil session")
	}
	if session.GetUserID() != user.GetID() {
		t.Errorf("Expected session user ID to be '%s', got '%s'", user.GetID(), session.GetUserID())
	}
}

func TestSeedUserAndSession_NilRequest(t *testing.T) {
	t.Parallel()
	app := Setup(WithUserStore(true), WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	user, session, err := SeedUserAndSession(app.GetUserStore(), app.GetSessionStore(), "test-user", nil, 10)
	if err == nil {
		t.Error("Expected error when request is nil")
	}
	if user != nil || session != nil {
		t.Error("Expected nil user and session when request is nil")
	}
}

func TestSeedUserAndSession_Success(t *testing.T) {
	t.Parallel()
	app := Setup(WithUserStore(true), WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	r, _ := NewRequest("GET", "/", NewRequestOptions{})
	user, session, err := SeedUserAndSession(app.GetUserStore(), app.GetSessionStore(), "test-user", r, 10)
	if err != nil {
		t.Fatalf("Failed to seed user and session: %v", err)
	}
	if user == nil {
		t.Fatal("Expected non-nil user")
	}
	if session == nil {
		t.Fatal("Expected non-nil session")
	}
	if session.GetUserID() != user.GetID() {
		t.Errorf("Expected session user ID to match user ID")
	}
}

func TestLoginAs_Success(t *testing.T) {
	t.Parallel()
	app := Setup(WithUserStore(true), WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	user, _ := SeedUser(app.GetUserStore(), "test-user")
	r, _ := NewRequest("GET", "/", NewRequestOptions{})

	authenticatedReq, err := LoginAs(app, r, user)
	if err != nil {
		t.Fatalf("Failed to login as user: %v", err)
	}
	if authenticatedReq == nil {
		t.Fatal("Expected non-nil request")
	}
}

func TestNewRequest_DefaultURL(t *testing.T) {
	t.Parallel()
	req, err := NewRequest("GET", "", NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	if req.URL.Path != "/" {
		t.Errorf("Expected default URL path to be '/', got '%s'", req.URL.Path)
	}
}

func TestNewRequest_WithBody(t *testing.T) {
	t.Parallel()
	body := "test body content"
	req, err := NewRequest("POST", "/test", NewRequestOptions{Body: body})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	if buf.String() != body {
		t.Errorf("Expected request body to be '%s', got '%s'", body, buf.String())
	}
}

func TestNewRequest_WithJSONData(t *testing.T) {
	t.Parallel()
	data := map[string]string{"key": "value"}
	req, err := NewRequest("POST", "/test", NewRequestOptions{JSONData: data})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type to be 'application/json', got '%s'", req.Header.Get("Content-Type"))
	}
}

func TestNewRequest_WithFormValues(t *testing.T) {
	t.Parallel()
	formValues := url.Values{}
	formValues.Set("field1", "value1")
	formValues.Set("field2", "value2")
	req, err := NewRequest("POST", "/test", NewRequestOptions{FormValues: formValues})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Errorf("Expected Content-Type to be 'application/x-www-form-urlencoded', got '%s'", req.Header.Get("Content-Type"))
	}
}

func TestNewRequest_WithQueryParams(t *testing.T) {
	t.Parallel()
	queryParams := url.Values{}
	queryParams.Set("param1", "value1")
	queryParams.Set("param2", "value2")
	req, err := NewRequest("GET", "/test", NewRequestOptions{QueryParams: queryParams})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	if req.URL.RawQuery != queryParams.Encode() {
		t.Errorf("Expected query params to match")
	}
}

func TestNewRequest_WithHeaders(t *testing.T) {
	t.Parallel()
	headers := map[string]string{
		"X-Custom-Header": "custom-value",
		"Authorization":   "Bearer token",
	}
	req, err := NewRequest("GET", "/test", NewRequestOptions{Headers: headers})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	if req.Header.Get("X-Custom-Header") != "custom-value" {
		t.Error("Expected custom header to be set")
	}
	if req.Header.Get("Authorization") != "Bearer token" {
		t.Error("Expected authorization header to be set")
	}
}

func TestNewRequest_WithContext(t *testing.T) {
	t.Parallel()
	ctxKey := "context-key"
	ctxValue := "context-value"
	ctx := map[any]any{ctxKey: ctxValue}
	req, err := NewRequest("GET", "/test", NewRequestOptions{Context: ctx})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	if req.Context().Value(ctxKey) != ctxValue {
		t.Error("Expected context value to be set")
	}
}

func TestNewRequest_RequestURI(t *testing.T) {
	t.Parallel()
	req, err := NewRequest("GET", "/test/path", NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	if req.RequestURI != "/test/path" {
		t.Errorf("Expected RequestURI to be '/test/path', got '%s'", req.RequestURI)
	}
}

func TestTestKey_WithConfig(t *testing.T) {
	t.Parallel()
	cfg := config.New()
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseHost("")
	cfg.SetDatabasePort("")
	cfg.SetDatabaseName("test.db")
	cfg.SetDatabaseUsername("")
	cfg.SetDatabasePassword("")

	key := TestKey(cfg)
	if key == "" {
		t.Error("Expected non-empty test key")
	}
}
