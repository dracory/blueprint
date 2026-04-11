package social

import (
	"strings"
	"testing"
)

func TestTo(t *testing.T) {
	// Test with no parameters
	result := to("https://example.com", map[string]string{})
	if result != "https://example.com" {
		t.Errorf("Expected https://example.com, got %s", result)
	}

	// Test with parameters
	result = to("https://example.com", map[string]string{
		"key1": "value1",
		"key2": "value2",
	})
	if result == "" {
		t.Error("Result should not be empty")
	}
	if !strings.Contains(result, "key1=") {
		t.Error("Result should contain key1 parameter")
	}
	if !strings.Contains(result, "key2=") {
		t.Error("Result should contain key2 parameter")
	}

	// Test with special characters that need encoding
	result = to("https://example.com", map[string]string{
		"title": "Test & Special",
	})
	if !strings.Contains(result, "Test+") || !strings.Contains(result, "%26") {
		t.Error("Result should have encoded special characters")
	}
}

func TestSocialMediaColors(t *testing.T) {
	colors := SocialMediaColors()

	// Test that major platforms have colors
	if colors[PlatformFacebook] != ColorFacebook {
		t.Errorf("Facebook color should be %s, got %s", ColorFacebook, colors[PlatformFacebook])
	}

	if colors[PlatformTwitter] != ColorTwitter {
		t.Errorf("Twitter color should be %s, got %s", ColorTwitter, colors[PlatformTwitter])
	}

	if colors[PlatformLinkedIn] != ColorLinkedIn {
		t.Errorf("LinkedIn color should be %s, got %s", ColorLinkedIn, colors[PlatformLinkedIn])
	}

	// Test that colors are valid hex values (start with # and have 7 characters)
	for platform, color := range colors {
		if len(color) != 7 || color[0] != '#' {
			t.Errorf("Color for %s should be a valid hex value (e.g., #RRGGBB), got %s", platform, color)
		}
	}

	// Test that default color is set for utility platforms
	if colors[PlatformEmail] != ColorDefault {
		t.Errorf("Email color should be default %s, got %s", ColorDefault, colors[PlatformEmail])
	}

	if colors[PlatformPrint] != ColorDefault {
		t.Errorf("Print color should be default %s, got %s", ColorDefault, colors[PlatformPrint])
	}
}
