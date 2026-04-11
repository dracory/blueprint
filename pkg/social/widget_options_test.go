package social

import (
	"strings"
	"testing"
)

func TestWidget_PopupEnabledByDefault(t *testing.T) {
	shareLinks := NewQuick("https://example.com", "Test", "")

	html := shareLinks.Widget(WidgetOptions{
		Platforms: []string{PlatformFacebook},
	})

	if !strings.Contains(html, "<script>") {
		t.Error("Widget should include popup script by default")
	}

	if !strings.Contains(html, "popupSize") {
		t.Error("Widget should include popup functionality by default")
	}
}

func TestWidget_PopupDisabled(t *testing.T) {
	shareLinks := NewQuick("https://example.com", "Test", "")

	disablePopup := false
	html := shareLinks.Widget(WidgetOptions{
		Platforms:   []string{PlatformFacebook},
		EnablePopup: &disablePopup,
	})

	if strings.Contains(html, "<script>") {
		t.Error("Widget should not include popup script when disabled")
	}

	if strings.Contains(html, "popupSize") {
		t.Error("Widget should not include popup functionality when disabled")
	}

	// Should still contain the widget HTML
	if !strings.Contains(html, `<div id="social-links">`) {
		t.Error("Widget should still contain HTML structure when popup disabled")
	}
}

func TestWidget_PopupExplicitlyEnabled(t *testing.T) {
	shareLinks := NewQuick("https://example.com", "Test", "")

	enablePopup := true
	html := shareLinks.Widget(WidgetOptions{
		Platforms:   []string{PlatformFacebook, PlatformTwitter},
		EnablePopup: &enablePopup,
	})

	if !strings.Contains(html, "<script>") {
		t.Error("Widget should include popup script when explicitly enabled")
	}

	if !strings.Contains(html, "window.open") {
		t.Error("Widget should include window.open call")
	}
}

func TestWidget_WithShareText(t *testing.T) {
	shareLinks := NewQuick("https://example.com", "Test", "")

	html := shareLinks.Widget(WidgetOptions{
		Platforms: []string{PlatformFacebook, PlatformTwitter},
		ShareText: "Share",
	})

	if !strings.Contains(html, `<span class="share-text">`) {
		t.Error("Widget should contain share text span")
	}

	if !strings.Contains(html, `Share`) {
		t.Error("Widget should contain the word Share")
	}

	if !strings.Contains(html, `<div id="social-links">`) {
		t.Error("Widget should still contain social-links div")
	}
}

func TestWidget_WithoutShareText(t *testing.T) {
	shareLinks := NewQuick("https://example.com", "Test", "")

	html := shareLinks.Widget(WidgetOptions{
		Platforms: []string{PlatformFacebook},
	})

	if strings.Contains(html, `<span class="share-text">`) {
		t.Error("Widget should not contain share text span when not provided")
	}
}
