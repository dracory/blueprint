package links

import (
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Admin Links Tests
// ============================================================================

func TestAdminLinks_Home(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.Home()
	if result == "" {
		t.Error("Home() should return non-empty string")
	}
	if !strings.Contains(result, "/admin") {
		t.Errorf("Home() = %q, should contain /admin", result)
	}
}

func TestAdminLinks_Blog(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.Blog()
	if result == "" {
		t.Error("Blog() should return non-empty string")
	}
	if !strings.Contains(result, "/admin/blog") {
		t.Errorf("Blog() = %q, should contain /admin/blog", result)
	}
}

func TestAdminLinks_Cms(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.Cms()
	if result == "" {
		t.Error("Cms() should return non-empty string")
	}
	if !strings.Contains(result, "/admin/cms") {
		t.Errorf("Cms() = %q, should contain /admin/cms", result)
	}
}

func TestAdminLinks_CmsOld(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.CmsOld()
	if result == "" {
		t.Error("CmsOld() should return non-empty string")
	}
	if !strings.Contains(result, "/admin/cmsold") {
		t.Errorf("CmsOld() = %q, should contain /admin/cmsold", result)
	}
}

func TestAdminLinks_FileManager(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.FileManager()
	if result == "" {
		t.Error("FileManager() should return non-empty string")
	}
	if !strings.Contains(result, "/admin/file-manager") {
		t.Errorf("FileManager() = %q, should contain /admin/file-manager", result)
	}
}

func TestAdminLinks_Logs(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.Logs()
	if result == "" {
		t.Error("Logs() should return non-empty string")
	}
	if !strings.Contains(result, "/admin/logs") {
		t.Errorf("Logs() = %q, should contain /admin/logs", result)
	}
}

func TestAdminLinks_MediaManager(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.MediaManager()
	if result == "" {
		t.Error("MediaManager() should return non-empty string")
	}
	if !strings.Contains(result, "/admin/media") {
		t.Errorf("MediaManager() = %q, should contain /admin/media", result)
	}
}

func TestAdminLinks_Shop(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.Shop()
	if result == "" {
		t.Error("Shop() should return non-empty string")
	}
	if !strings.Contains(result, "/admin/shop") {
		t.Errorf("Shop() = %q, should contain /admin/shop", result)
	}
}

func TestAdminLinks_Stats(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.Stats()
	if result == "" {
		t.Error("Stats() should return non-empty string")
	}
	if !strings.Contains(result, "/admin/stats") {
		t.Errorf("Stats() = %q, should contain /admin/stats", result)
	}
}

func TestAdminLinks_Tasks(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.Tasks()
	if result == "" {
		t.Error("Tasks() should return non-empty string")
	}
	if !strings.Contains(result, "/admin/tasks") {
		t.Errorf("Tasks() = %q, should contain /admin/tasks", result)
	}
}

func TestAdminLinks_Users(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.Users()
	if result == "" {
		t.Error("Users() should return non-empty string")
	}
	if !strings.Contains(result, "/admin/users") {
		t.Errorf("Users() = %q, should contain /admin/users", result)
	}
}

func TestAdminLinks_WithParams_Home(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.Home(map[string]string{"key": "value"})
	if result == "" {
		t.Error("Home with params should return non-empty string")
	}
}

func TestAdminLinks_WithParams_Blog(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.Blog(map[string]string{"page": "1"})
	if result == "" {
		t.Error("Blog with params should return non-empty string")
	}
}

func TestAdminLinks_WithParams_Users(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	admin := Admin()
	result := admin.Users(map[string]string{"sort": "name"})
	if result == "" {
		t.Error("Users with params should return non-empty string")
	}
}

// ============================================================================
// Website Links Tests
// ============================================================================

func TestWebsiteLinks_Home(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	website := Website()
	result := website.Home()
	if result == "" {
		t.Error("Home() should return non-empty string")
	}
	if !strings.Contains(result, "/") {
		t.Errorf("Home() = %q, should contain /", result)
	}
}

func TestWebsiteLinks_Blog(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	website := Website()
	result := website.Blog()
	if result == "" {
		t.Error("Blog() should return non-empty string")
	}
	if !strings.Contains(result, "/blog") {
		t.Errorf("Blog() = %q, should contain /blog", result)
	}
}

func TestWebsiteLinks_Chat(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	website := Website()
	result := website.Chat()
	if result == "" {
		t.Error("Chat() should return non-empty string")
	}
	if !strings.Contains(result, "/chat") {
		t.Errorf("Chat() = %q, should contain /chat", result)
	}
}

func TestWebsiteLinks_Contact(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	website := Website()
	result := website.Contact()
	if result == "" {
		t.Error("Contact() should return non-empty string")
	}
	if !strings.Contains(result, "/contact") {
		t.Errorf("Contact() = %q, should contain /contact", result)
	}
}

func TestWebsiteLinks_Shop(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	website := Website()
	result := website.Shop()
	if result == "" {
		t.Error("Shop() should return non-empty string")
	}
	if !strings.Contains(result, "/shop") {
		t.Errorf("Shop() = %q, should contain /shop", result)
	}
}

func TestWebsiteLinks_SitemapXml(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	website := Website()
	result := website.SitemapXml()
	if result == "" {
		t.Error("SitemapXml() should return non-empty string")
	}
	if !strings.Contains(result, "/sitemap.xml") {
		t.Errorf("SitemapXml() = %q, should contain /sitemap.xml", result)
	}
}

func TestWebsiteLinks_WithParams_Blog(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	website := Website()
	result := website.Blog(map[string]string{"page": "1"})
	if result == "" {
		t.Error("Blog with params should return non-empty string")
	}
}

func TestWebsiteLinks_WithParams_Chat(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	website := Website()
	result := website.Chat(map[string]string{"room": "general"})
	if result == "" {
		t.Error("Chat with params should return non-empty string")
	}
}

func TestWebsiteLinks_WithParams_Shop(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	website := Website()
	result := website.Shop(map[string]string{"category": "all"})
	if result == "" {
		t.Error("Shop with params should return non-empty string")
	}
}

func TestWebsiteLinks_BlogPost(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	website := Website()

	result := website.BlogPost("123", "test-post")
	if result == "" {
		t.Error("BlogPost() should return non-empty string")
	}
	if !strings.Contains(result, "123") {
		t.Error("BlogPost() should contain post ID")
	}
	if !strings.Contains(result, "test-post") {
		t.Error("BlogPost() should contain post slug")
	}
}

func TestWebsiteLinks_File(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	website := Website()

	result := website.File("uploads/test.jpg")
	if result == "" {
		t.Error("File() should return non-empty string")
	}
	if !strings.Contains(result, "uploads") {
		t.Error("File() should contain path")
	}
}

func TestWebsiteLinks_Flash(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	website := Website()

	result := website.Flash()
	if result == "" {
		t.Error("Flash() should return non-empty string")
	}

	result = website.Flash(map[string]string{"type": "success"})
	if result == "" {
		t.Error("Flash(params) should return non-empty string")
	}
}

func TestWebsiteLinks_ShopProduct(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	website := Website()

	// With slug
	result := website.ShopProduct("123", "test-product", nil)
	if result == "" {
		t.Error("ShopProduct() should return non-empty string")
	}
	if !strings.Contains(result, "123") {
		t.Error("ShopProduct() should contain product ID")
	}
	if !strings.Contains(result, "test-product") {
		t.Error("ShopProduct() should contain product slug")
	}

	// Without slug
	result = website.ShopProduct("456", "", nil)
	if !strings.Contains(result, "456") {
		t.Error("ShopProduct() should contain product ID without slug")
	}

	// With params
	result = website.ShopProduct("789", "product", map[string]string{"ref": "ad"})
	if result == "" {
		t.Error("ShopProduct(params) should return non-empty string")
	}
}

func TestWebsiteLinks_Payment(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	website := Website()

	canceled := website.PaymentCanceled("payment-key-123")
	if canceled == "" {
		t.Error("PaymentCanceled() should return non-empty string")
	}
	if !strings.Contains(canceled, "payment-key-123") {
		t.Error("PaymentCanceled() should contain payment key")
	}

	success := website.PaymentSuccess("payment-key-456")
	if success == "" {
		t.Error("PaymentSuccess() should return non-empty string")
	}
	if !strings.Contains(success, "payment-key-456") {
		t.Error("PaymentSuccess() should contain payment key")
	}
}

func TestWebsiteLinks_Resource(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	website := Website()

	// Empty path
	result := website.Resource("")
	if result != "" {
		t.Error("Resource(\"\") should return empty string")
	}

	// With leading slash
	result = website.Resource("/css/style.css")
	if result == "" {
		t.Error("Resource(\"/css/style.css\") should return non-empty string")
	}

	// Without leading slash
	result = website.Resource("js/app.js")
	if result == "" {
		t.Error("Resource(\"js/app.js\") should return non-empty string")
	}

	// With params
	result = website.Resource("/img/logo.png", map[string]string{"v": "2"})
	if result == "" {
		t.Error("Resource with params should return non-empty string")
	}
}

func TestWebsiteLinks_Theme(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	website := Website()

	result := website.Theme(map[string]string{"name": "default"})
	if result == "" {
		t.Error("Theme() should return non-empty string")
	}
}

func TestWebsiteLinks_Widget(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	website := Website()

	result := website.Widget("sidebar", map[string]string{"position": "left"})
	if result == "" {
		t.Error("Widget() should return non-empty string")
	}
	if !strings.Contains(result, "sidebar") {
		t.Error("Widget() should contain alias")
	}
}

func TestWebsiteLinks_Thumbnail(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "http://localhost:8080")

	website := Website()

	// Basic thumbnail
	result := website.Thumbnail("jpg", "200", "100", "90", "uploads/test.jpg", nil)
	if result == "" {
		t.Error("Thumbnail() should return non-empty string")
	}
	if !strings.Contains(result, "/th/") {
		t.Error("Thumbnail() should contain /th/ prefix")
	}
	if !strings.Contains(result, "200x100") {
		t.Error("Thumbnail() should contain dimensions")
	}

	// With defaults
	result = website.Thumbnail("", "", "", "", "test.png", nil)
	if !strings.Contains(result, "100x100") {
		t.Error("Thumbnail() should use default dimensions")
	}
	if !strings.Contains(result, "png") {
		t.Error("Thumbnail() should use default extension")
	}

	// With HTTP URL in path
	result = website.Thumbnail("jpg", "300", "200", "80", "http://example.com/image.jpg", nil)
	if !strings.Contains(result, "http/") {
		t.Error("Thumbnail() should replace http:// with http/")
	}

	// With HTTPS URL in path
	result = website.Thumbnail("jpg", "300", "200", "80", "https://example.com/image.jpg", nil)
	if !strings.Contains(result, "https/") {
		t.Error("Thumbnail() should replace https:// with https/")
	}
}

func TestWebsiteLinks_Thumbnail_WithDataURI(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "http://localhost:8080")

	website := Website()

	// Data URI should trigger cache
	dataURI := "data:image/png;base64,iVBORw0KGgo="
	result := website.Thumbnail("png", "100", "100", "80", dataURI, nil)

	// Result should be generated (may use cache or return error URL)
	if result == "" {
		t.Error("Thumbnail() with data URI should return non-empty string")
	}
}

// ============================================================================
// User Links Tests
// ============================================================================

func TestUserLinks_Subscriptions_PlanSelect(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	user := User()
	result := user.SubscriptionsPlanSelect()
	if result == "" {
		t.Error("SubscriptionsPlanSelect() should return non-empty string")
	}
	if !strings.Contains(result, "/subscription") {
		t.Errorf("SubscriptionsPlanSelect() = %q, should contain /subscription", result)
	}
}

func TestUserLinks_Subscriptions_PlanSelectAjax(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	user := User()
	result := user.SubscriptionsPlanSelectAjax()
	if result == "" {
		t.Error("SubscriptionsPlanSelectAjax() should return non-empty string")
	}
	if !strings.Contains(result, "/plan-select-ajax") {
		t.Errorf("SubscriptionsPlanSelectAjax() = %q, should contain /plan-select-ajax", result)
	}
}

func TestUserLinks_Subscriptions_PaymentSuccess(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	user := User()
	result := user.SubscriptionsPaymentSuccess()
	if result == "" {
		t.Error("SubscriptionsPaymentSuccess() should return non-empty string")
	}
	if !strings.Contains(result, "/payment-success") {
		t.Errorf("SubscriptionsPaymentSuccess() = %q, should contain /payment-success", result)
	}
}

func TestUserLinks_Subscriptions_PaymentCanceled(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")
	user := User()
	result := user.SubscriptionsPaymentCanceled()
	if result == "" {
		t.Error("SubscriptionsPaymentCanceled() should return non-empty string")
	}
	if !strings.Contains(result, "/payment-canceled") {
		t.Errorf("SubscriptionsPaymentCanceled() = %q, should contain /payment-canceled", result)
	}
}

func TestUserLinks_WithParams(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	user := User()

	result := user.SubscriptionsPlanSelect(map[string]string{"plan": "premium"})
	if result == "" {
		t.Error("SubscriptionsPlanSelect(params) should return non-empty string")
	}

	result = user.SubscriptionsPaymentSuccess(map[string]string{"order": "123"})
	if result == "" {
		t.Error("SubscriptionsPaymentSuccess(params) should return non-empty string")
	}
}

// ============================================================================
// Constants Tests
// ============================================================================

func TestConstants(t *testing.T) {
	// Test that constants are defined
	if CATCHALL != "/*" {
		t.Errorf("CATCHALL = %q, want /*", CATCHALL)
	}

	// Auth constants
	if AUTH_AUTH != "/auth/auth" {
		t.Error("AUTH_AUTH constant incorrect")
	}
	if AUTH_LOGIN != "/auth/login" {
		t.Error("AUTH_LOGIN constant incorrect")
	}
	if AUTH_LOGOUT != "/auth/logout" {
		t.Error("AUTH_LOGOUT constant incorrect")
	}
	if AUTH_REGISTER != "/auth/register" {
		t.Error("AUTH_REGISTER constant incorrect")
	}

	// Admin constants
	if ADMIN_HOME != "/admin" {
		t.Error("ADMIN_HOME constant incorrect")
	}
	if ADMIN_BLOG != "/admin/blog" {
		t.Error("ADMIN_BLOG constant incorrect")
	}

	// User constants
	if USER_HOME != "/user" {
		t.Error("USER_HOME constant incorrect")
	}
	if USER_PROFILE != "/user/profile" {
		t.Error("USER_PROFILE constant incorrect")
	}
}

// ============================================================================
// Multiple Instances Tests
// ============================================================================

func TestMultipleInstances(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	// Test that multiple instances work independently
	admin1 := Admin()
	admin2 := Admin()

	if admin1 == admin2 {
		t.Error("Admin() should return different instances")
	}

	website1 := Website()
	website2 := Website()

	if website1 == website2 {
		t.Error("Website() should return different instances")
	}

	user1 := User()
	user2 := User()

	if user1 == user2 {
		t.Error("User() should return different instances")
	}

	auth1 := Auth()
	auth2 := Auth()

	if auth1 == auth2 {
		t.Error("Auth() should return different instances")
	}
}

// ============================================================================
// URL Builder Edge Cases
// ============================================================================

func TestURL_EdgeCases(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "")

	// Empty path
	result := URL("", nil)
	if result == "" {
		t.Error("URL(\"\", nil) should return non-empty string (root URL)")
	}

	// Path with special characters
	result = URL("/path with spaces", nil)
	if result == "" {
		t.Error("URL with spaces should return non-empty string")
	}

	// Multiple calls to initializeURLBuilder should work
	initializeURLBuilder()
	initializeURLBuilder()
	result = RootURL()
	if result == "" {
		t.Error("RootURL after multiple init calls should work")
	}
}

// MockCache for testing Thumbnail with cache
type mockCache struct{}

func (m *mockCache) Save(key string, value string, lifetime time.Duration) error {
	return nil
}

func (m *mockCache) Fetch(key string) (string, error) {
	return "", nil
}

func (m *mockCache) FetchMulti(keys []string) map[string]string {
	return map[string]string{}
}

func (m *mockCache) Flush() error {
	return nil
}

func (m *mockCache) Delete(key string) error {
	return nil
}

func (m *mockCache) Contains(key string) bool {
	return false
}

func TestWebsiteLinks_Thumbnail_WithMockCache(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_URL", "http://localhost:8080")

	website := Website()
	mock := &mockCache{}

	// Data URI with mock cache
	dataURI := "data:image/png;base64,test"
	result := website.Thumbnail("png", "100", "100", "80", dataURI, mock)

	if result == "" {
		t.Error("Thumbnail() with mock cache should return non-empty string")
	}
}

// mockCache is a minimal cache implementation for testing
