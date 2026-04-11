package social

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	postURL := "https://example.com/blog/post/123/test-blog-post-with-special-chars-&-symbols"
	postTitle := "Test Post with Special Characters & Symbols!"
	postImageURL := "https://example.com/image.jpg"

	shareLinks := NewQuick(postURL, postTitle, postImageURL)

	// Test that all URLs are generated
	facebookURL := shareLinks.GetFacebookShareUrl()
	twitterURL := shareLinks.GetTwitterShareUrl()
	linkedinURL := shareLinks.GetLinkedInShareUrl()
	pinterestURL := shareLinks.GetPinterestShareUrl()

	if facebookURL == "" {
		t.Error("Facebook URL should not be empty")
	}
	if twitterURL == "" {
		t.Error("Twitter URL should not be empty")
	}
	if linkedinURL == "" {
		t.Error("LinkedIn URL should not be empty")
	}
	if pinterestURL == "" {
		t.Error("Pinterest URL should not be empty")
	}

	// Test that URLs contain the expected base URLs and parameters
	if !strings.Contains(facebookURL, "https://www.facebook.com/sharer.php?") {
		t.Errorf("Facebook URL should contain base URL, got %s", facebookURL)
	}
	if !strings.Contains(facebookURL, "u=") {
		t.Error("Facebook URL should contain u parameter")
	}

	if !strings.Contains(twitterURL, "https://x.com/intent/tweet?") {
		t.Errorf("Twitter URL should contain base URL, got %s", twitterURL)
	}
	if !strings.Contains(twitterURL, "url=") {
		t.Error("Twitter URL should contain url parameter")
	}

	if !strings.Contains(linkedinURL, "https://www.linkedin.com/sharing/share-offsite/?") {
		t.Errorf("LinkedIn URL should contain base URL, got %s", linkedinURL)
	}
	if !strings.Contains(linkedinURL, "url=") {
		t.Error("LinkedIn URL should contain url parameter")
	}

	if !strings.Contains(pinterestURL, "https://pinterest.com/pin/create/button/?") {
		t.Errorf("Pinterest URL should contain base URL, got %s", pinterestURL)
	}
	if !strings.Contains(pinterestURL, "url=") {
		t.Error("Pinterest URL should contain url parameter")
	}

	// Test that URLs contain encoded characters (should contain %26 for &)
	if !strings.Contains(facebookURL, "%26") {
		t.Error("Facebook URL should contain encoded & character (%26)")
	}
	if !strings.Contains(twitterURL, "%26") {
		t.Error("Twitter URL should contain encoded & character (%26)")
	}
	if !strings.Contains(linkedinURL, "%26") {
		t.Error("LinkedIn URL should contain encoded & character (%26)")
	}
	if !strings.Contains(pinterestURL, "%26") {
		t.Error("Pinterest URL should contain encoded & character (%26)")
	}

	// Test that Pinterest URL contains the image URL
	if !strings.Contains(pinterestURL, "media=") {
		t.Error("Pinterest URL should contain media parameter for image")
	}
}
