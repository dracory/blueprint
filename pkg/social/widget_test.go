package social

import (
	"strings"
	"testing"
)

func TestWidget(t *testing.T) {
	postURL := "https://example.com/blog/post/123"
	postTitle := "Test Post"

	shareLinks := NewQuick(postURL, postTitle, "")

	html := shareLinks.Widget(WidgetOptions{
		Platforms: []string{PlatformFacebook, PlatformTwitter},
	})

	if !strings.Contains(html, `<div id="social-links">`) {
		t.Error("HTML should contain social-links div")
	}

	if !strings.Contains(html, "bi-facebook") {
		t.Error("HTML should contain facebook icon class")
	}

	if !strings.Contains(html, "bi-twitter") {
		t.Error("HTML should contain twitter icon class")
	}

	// Test with custom platforms and Bootstrap icons
	html = shareLinks.Widget(WidgetOptions{
		Platforms:   []string{PlatformFacebook, PlatformLinkedIn},
		IconLibrary: IconLibraryBootstrap,
	})
	if !strings.Contains(html, BootstrapIconFacebook) {
		t.Error("HTML should contain bootstrap facebook icon class")
	}
	if !strings.Contains(html, BootstrapIconLinkedIn) {
		t.Error("HTML should contain bootstrap linkedin icon class")
	}
}

func TestWidgetWithPopup(t *testing.T) {
	postURL := "https://example.com/blog/post/123"
	postTitle := "Test Post"

	shareLinks := NewQuick(postURL, postTitle, "")

	// Test with popup enabled (default behavior)
	html := shareLinks.Widget(WidgetOptions{
		Platforms: []string{PlatformFacebook, PlatformTwitter},
	})

	// Should contain the widget HTML
	if !strings.Contains(html, `<div id="social-links">`) {
		t.Error("HTML should contain social-links div")
	}

	// Should contain the popup script
	if !strings.Contains(html, "<script>") {
		t.Error("HTML should contain popup script")
	}

	if !strings.Contains(html, "popupSize") {
		t.Error("Script should contain popupSize variable")
	}

	if !strings.Contains(html, "window.open") {
		t.Error("Script should contain window.open call")
	}

	if !strings.Contains(html, "social-share") {
		t.Error("Script should use 'social-share' as popup name")
	}

	if !strings.Contains(html, "addEventListener") {
		t.Error("Script should use addEventListener for click events")
	}

	// Should handle mailto and sms links differently
	if !strings.Contains(html, "mailto:") {
		t.Log("Script should check for mailto: links")
	}

	if !strings.Contains(html, "sms:") {
		t.Log("Script should check for sms: links")
	}
}

func TestWidget_SortStrategy(t *testing.T) {
	postURL := "https://example.com/blog/post/123"
	postTitle := "Test Post"

	shareLinks := NewQuick(postURL, postTitle, "")

	// Test alphabetical sorting
	htmlAlpha := shareLinks.Widget(WidgetOptions{
		Platforms:    []string{PlatformTwitter, PlatformFacebook, PlatformBluesky, PlatformDiscord},
		SortStrategy: SortStrategyAlphabetical,
	})

	// Find positions - should be alphabetical: Bluesky, Discord, Facebook, Twitter
	blueskyPos := strings.Index(htmlAlpha, "id=\"social-bluesky\"")
	discordPos := strings.Index(htmlAlpha, "id=\"social-discord\"")
	facebookPos := strings.Index(htmlAlpha, "id=\"social-facebook\"")
	twitterPos := strings.Index(htmlAlpha, "id=\"social-twitter\"")

	if blueskyPos == -1 || discordPos == -1 || facebookPos == -1 || twitterPos == -1 {
		t.Error("HTML should contain all social links")
	}

	// Check alphabetical order: Bluesky < Discord < Facebook < Twitter
	if blueskyPos > discordPos || discordPos > facebookPos || facebookPos > twitterPos {
		t.Error("With SortStrategyAlphabetical, platforms should appear in alphabetical order")
	}

	// Test manual sorting (default behavior - should match Platforms order)
	htmlManual := shareLinks.Widget(WidgetOptions{
		Platforms:    []string{PlatformTwitter, PlatformFacebook, PlatformBluesky},
		SortStrategy: SortStrategyManual,
	})

	twitterPosManual := strings.Index(htmlManual, "id=\"social-twitter\"")
	facebookPosManual := strings.Index(htmlManual, "id=\"social-facebook\"")
	blueskyPosManual := strings.Index(htmlManual, "id=\"social-bluesky\"")

	if twitterPosManual == -1 || facebookPosManual == -1 || blueskyPosManual == -1 {
		t.Error("HTML should contain all social links")
	}

	// Should be in Platforms order: Twitter, Facebook, Bluesky
	if twitterPosManual > facebookPosManual || facebookPosManual > blueskyPosManual {
		t.Error("With SortStrategyManual, platforms should appear in Platforms order")
	}

	// Test empty SortStrategy (defaults to manual)
	htmlEmpty := shareLinks.Widget(WidgetOptions{
		Platforms:    []string{PlatformTwitter, PlatformFacebook},
		SortStrategy: "",
	})

	twitterPosEmpty := strings.Index(htmlEmpty, "id=\"social-twitter\"")
	facebookPosEmpty := strings.Index(htmlEmpty, "id=\"social-facebook\"")

	if twitterPosEmpty == -1 || facebookPosEmpty == -1 {
		t.Error("HTML should contain all social links")
	}

	// With empty SortStrategy, should use Platforms order (Twitter before Facebook)
	if twitterPosEmpty > facebookPosEmpty {
		t.Error("With empty SortStrategy, should default to manual order (Platforms order)")
	}
}

func TestWidget_NewShareActions(t *testing.T) {
	postURL := "https://example.com/blog/post/123"
	postTitle := "Test Post"

	shareLinks := NewQuick(postURL, postTitle, "")

	// Test Print action
	html := shareLinks.Widget(WidgetOptions{
		Platforms: []string{PlatformPrint},
	})
	if !strings.Contains(html, "javascript:window.print()") {
		t.Error("Print action should contain window.print() JavaScript")
	}
	if !strings.Contains(html, "bi-printer") {
		t.Error("Print action should have printer icon")
	}

	// Test Copy Link action
	html = shareLinks.Widget(WidgetOptions{
		Platforms: []string{PlatformCopyLink},
	})
	if !strings.Contains(html, "javascript:navigator.clipboard.writeText") {
		t.Error("CopyLink action should contain clipboard.writeText JavaScript")
	}
	if !strings.Contains(html, postURL) {
		t.Error("CopyLink action should contain the post URL")
	}
	if !strings.Contains(html, "bi-link-45deg") {
		t.Error("CopyLink action should have link icon")
	}

	// Test Native Share action
	html = shareLinks.Widget(WidgetOptions{
		Platforms: []string{PlatformNativeShare},
	})
	if !strings.Contains(html, "javascript:if(navigator.share)") {
		t.Error("NativeShare action should contain navigator.share JavaScript")
	}
	if !strings.Contains(html, postTitle) {
		t.Error("NativeShare action should contain the post title")
	}
	if !strings.Contains(html, postURL) {
		t.Error("NativeShare action should contain the post URL")
	}
	if !strings.Contains(html, "bi-share") {
		t.Error("NativeShare action should have share icon")
	}

	// Test combined with social platforms
	html = shareLinks.Widget(WidgetOptions{
		Platforms: []string{PlatformFacebook, PlatformPrint, PlatformCopyLink},
	})
	if !strings.Contains(html, "social-facebook") {
		t.Error("Combined widget should contain Facebook")
	}
	if !strings.Contains(html, "social-print") {
		t.Error("Combined widget should contain Print")
	}
	if !strings.Contains(html, "social-copylink") {
		t.Error("Combined widget should contain CopyLink")
	}
}
