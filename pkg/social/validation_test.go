package social

import (
	"strings"
	"testing"
)

func TestWidget_XSSProtection(t *testing.T) {
	malicious := `<script>alert('xss')</script>`
	shareLinks := NewQuick("https://example.com", malicious, "")

	html := shareLinks.Widget(WidgetOptions{
		Platforms: []string{PlatformFacebook},
	})

	// The malicious script should be URL-encoded in the href, not executed as HTML
	if strings.Contains(html, "<script>alert") {
		t.Error("Widget should not contain unescaped script tags in HTML")
	}

	// Verify the HTML structure is intact
	if !strings.Contains(html, `<div id="social-links">`) {
		t.Error("Widget should contain proper HTML structure")
	}
}

func TestWidget_UniqueIDs(t *testing.T) {
	shareLinks := NewQuick("https://example.com", "Test", "")

	html := shareLinks.Widget(WidgetOptions{
		Platforms: []string{PlatformFacebook, PlatformTwitter},
	})

	if !strings.Contains(html, `id="social-facebook"`) {
		t.Error("Widget should contain unique ID for Facebook")
	}

	if !strings.Contains(html, `id="social-twitter"`) {
		t.Error("Widget should contain unique ID for Twitter")
	}

	if strings.Contains(html, `id=""`) {
		t.Error("Widget should not contain empty ID attributes")
	}
}

func TestEmptyParameters(t *testing.T) {
	shareLinks := NewQuick("", "", "")

	facebookURL := shareLinks.GetFacebookShareUrl()
	if facebookURL == "" {
		t.Error("Should generate URL even with empty params")
	}

	if !strings.Contains(facebookURL, "facebook.com") {
		t.Error("URL should still contain base Facebook URL")
	}
}

func TestTo_SkipsEmptyValues(t *testing.T) {
	result := to("https://example.com", map[string]string{
		"key1": "value1",
		"key2": "",
		"key3": "value3",
	})

	if strings.Contains(result, "key2=") {
		t.Error("Should not include empty parameters in URL")
	}

	if !strings.Contains(result, "key1=value1") {
		t.Error("Should include non-empty parameters")
	}

	if !strings.Contains(result, "key3=value3") {
		t.Error("Should include non-empty parameters")
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid HTTPS URL", "https://example.com", true},
		{"valid HTTP URL", "http://example.com", true},
		{"valid URL with path", "https://example.com/path/to/page", true},
		{"valid URL with query", "https://example.com?foo=bar", true},
		{"empty string", "", false},
		{"missing scheme", "example.com", false},
		{"invalid scheme", "ftp://example.com", false},
		{"javascript scheme", "javascript:alert('xss')", false},
		{"data URI", "data:text/html,<script>alert('xss')</script>", false},
		{"relative URL", "/path/to/page", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateURL(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateURL(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid email", "user@example.com", true},
		{"valid email with subdomain", "user@mail.example.com", true},
		{"valid email with plus", "user+tag@example.com", true},
		{"valid email with dots", "first.last@example.com", true},
		{"valid email with numbers", "user123@example.com", true},
		{"empty string", "", false},
		{"missing @", "userexample.com", false},
		{"missing domain", "user@", false},
		{"missing local part", "@example.com", false},
		{"spaces in email", "user @example.com", false},
		{"multiple @ signs", "user@@example.com", false},
		{"invalid TLD", "user@example", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateEmail(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"US number with dashes", "555-123-4567", true},
		{"US number with spaces", "555 123 4567", true},
		{"international with plus", "+1-555-123-4567", true},
		{"international UK format", "+44 20 7946 0958", true},
		{"with parentheses", "(555) 123-4567", true},
		{"with dots", "555.123.4567", true},
		{"empty string", "", false},
		{"too short", "123", false},
		{"too long", "+1234567890123456", false},
		{"letters not allowed", "555-abc-1234", false},
		{"special chars not allowed", "555@123#4567", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePhone(tt.input)
			if result != tt.expected {
				t.Errorf("ValidatePhone(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
