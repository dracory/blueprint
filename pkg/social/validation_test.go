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

func TestValidateURL_ValidHTTPS(t *testing.T) {
	result := ValidateURL("https://example.com")
	if result != true {
		t.Errorf("ValidateURL(\"https://example.com\") = %v, want true", result)
	}
}

func TestValidateURL_ValidHTTP(t *testing.T) {
	result := ValidateURL("http://example.com")
	if result != true {
		t.Errorf("ValidateURL(\"http://example.com\") = %v, want true", result)
	}
}

func TestValidateURL_ValidURLWithPath(t *testing.T) {
	result := ValidateURL("https://example.com/path/to/page")
	if result != true {
		t.Errorf("ValidateURL(\"https://example.com/path/to/page\") = %v, want true", result)
	}
}

func TestValidateURL_ValidURLWithQuery(t *testing.T) {
	result := ValidateURL("https://example.com?foo=bar")
	if result != true {
		t.Errorf("ValidateURL(\"https://example.com?foo=bar\") = %v, want true", result)
	}
}

func TestValidateURL_EmptyString(t *testing.T) {
	result := ValidateURL("")
	if result != false {
		t.Errorf("ValidateURL(\"\") = %v, want false", result)
	}
}

func TestValidateURL_MissingScheme(t *testing.T) {
	result := ValidateURL("example.com")
	if result != false {
		t.Errorf("ValidateURL(\"example.com\") = %v, want false", result)
	}
}

func TestValidateURL_InvalidScheme(t *testing.T) {
	result := ValidateURL("ftp://example.com")
	if result != false {
		t.Errorf("ValidateURL(\"ftp://example.com\") = %v, want false", result)
	}
}

func TestValidateURL_JavascriptScheme(t *testing.T) {
	result := ValidateURL("javascript:alert('xss')")
	if result != false {
		t.Errorf("ValidateURL(\"javascript:alert('xss')\") = %v, want false", result)
	}
}

func TestValidateURL_DataURI(t *testing.T) {
	result := ValidateURL("data:text/html,<script>alert('xss')</script>")
	if result != false {
		t.Errorf("ValidateURL(\"data:text/html,<script>alert('xss')</script>\") = %v, want false", result)
	}
}

func TestValidateURL_RelativeURL(t *testing.T) {
	result := ValidateURL("/path/to/page")
	if result != false {
		t.Errorf("ValidateURL(\"/path/to/page\") = %v, want false", result)
	}
}

func TestValidateEmail_ValidEmail(t *testing.T) {
	result := ValidateEmail("user@example.com")
	if result != true {
		t.Errorf("ValidateEmail(\"user@example.com\") = %v, want true", result)
	}
}

func TestValidateEmail_ValidEmailWithSubdomain(t *testing.T) {
	result := ValidateEmail("user@mail.example.com")
	if result != true {
		t.Errorf("ValidateEmail(\"user@mail.example.com\") = %v, want true", result)
	}
}

func TestValidateEmail_ValidEmailWithPlus(t *testing.T) {
	result := ValidateEmail("user+tag@example.com")
	if result != true {
		t.Errorf("ValidateEmail(\"user+tag@example.com\") = %v, want true", result)
	}
}

func TestValidateEmail_ValidEmailWithDots(t *testing.T) {
	result := ValidateEmail("first.last@example.com")
	if result != true {
		t.Errorf("ValidateEmail(\"first.last@example.com\") = %v, want true", result)
	}
}

func TestValidateEmail_ValidEmailWithNumbers(t *testing.T) {
	result := ValidateEmail("user123@example.com")
	if result != true {
		t.Errorf("ValidateEmail(\"user123@example.com\") = %v, want true", result)
	}
}

func TestValidateEmail_EmptyString(t *testing.T) {
	result := ValidateEmail("")
	if result != false {
		t.Errorf("ValidateEmail(\"\") = %v, want false", result)
	}
}

func TestValidateEmail_MissingAtSign(t *testing.T) {
	result := ValidateEmail("userexample.com")
	if result != false {
		t.Errorf("ValidateEmail(\"userexample.com\") = %v, want false", result)
	}
}

func TestValidateEmail_MissingDomain(t *testing.T) {
	result := ValidateEmail("user@")
	if result != false {
		t.Errorf("ValidateEmail(\"user@\") = %v, want false", result)
	}
}

func TestValidateEmail_MissingLocalPart(t *testing.T) {
	result := ValidateEmail("@example.com")
	if result != false {
		t.Errorf("ValidateEmail(\"@example.com\") = %v, want false", result)
	}
}

func TestValidateEmail_SpacesInEmail(t *testing.T) {
	result := ValidateEmail("user @example.com")
	if result != false {
		t.Errorf("ValidateEmail(\"user @example.com\") = %v, want false", result)
	}
}

func TestValidateEmail_MultipleAtSigns(t *testing.T) {
	result := ValidateEmail("user@@example.com")
	if result != false {
		t.Errorf("ValidateEmail(\"user@@example.com\") = %v, want false", result)
	}
}

func TestValidateEmail_InvalidTLD(t *testing.T) {
	result := ValidateEmail("user@example")
	if result != false {
		t.Errorf("ValidateEmail(\"user@example\") = %v, want false", result)
	}
}

func TestValidatePhone_USNumberWithDashes(t *testing.T) {
	result := ValidatePhone("555-123-4567")
	if result != true {
		t.Errorf("ValidatePhone(\"555-123-4567\") = %v, want true", result)
	}
}

func TestValidatePhone_USNumberWithSpaces(t *testing.T) {
	result := ValidatePhone("555 123 4567")
	if result != true {
		t.Errorf("ValidatePhone(\"555 123 4567\") = %v, want true", result)
	}
}

func TestValidatePhone_InternationalWithPlus(t *testing.T) {
	result := ValidatePhone("+1-555-123-4567")
	if result != true {
		t.Errorf("ValidatePhone(\"+1-555-123-4567\") = %v, want true", result)
	}
}

func TestValidatePhone_InternationalUKFormat(t *testing.T) {
	result := ValidatePhone("+44 20 7946 0958")
	if result != true {
		t.Errorf("ValidatePhone(\"+44 20 7946 0958\") = %v, want true", result)
	}
}

func TestValidatePhone_WithParentheses(t *testing.T) {
	result := ValidatePhone("(555) 123-4567")
	if result != true {
		t.Errorf("ValidatePhone(\"(555) 123-4567\") = %v, want true", result)
	}
}

func TestValidatePhone_WithDots(t *testing.T) {
	result := ValidatePhone("555.123.4567")
	if result != true {
		t.Errorf("ValidatePhone(\"555.123.4567\") = %v, want true", result)
	}
}

func TestValidatePhone_EmptyString(t *testing.T) {
	result := ValidatePhone("")
	if result != false {
		t.Errorf("ValidatePhone(\"\") = %v, want false", result)
	}
}

func TestValidatePhone_TooShort(t *testing.T) {
	result := ValidatePhone("123")
	if result != false {
		t.Errorf("ValidatePhone(\"123\") = %v, want false", result)
	}
}

func TestValidatePhone_TooLong(t *testing.T) {
	result := ValidatePhone("+1234567890123456")
	if result != false {
		t.Errorf("ValidatePhone(\"+1234567890123456\") = %v, want false", result)
	}
}

func TestValidatePhone_LettersNotAllowed(t *testing.T) {
	result := ValidatePhone("555-abc-1234")
	if result != false {
		t.Errorf("ValidatePhone(\"555-abc-1234\") = %v, want false", result)
	}
}

func TestValidatePhone_SpecialCharsNotAllowed(t *testing.T) {
	result := ValidatePhone("555@123#4567")
	if result != false {
		t.Errorf("ValidatePhone(\"555@123#4567\") = %v, want false", result)
	}
}
