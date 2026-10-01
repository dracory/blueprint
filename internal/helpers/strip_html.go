package helpers

import (
	"html"
	"regexp"
	"strings"
)

// StripHTMLTags removes HTML tags from a string and decodes HTML entities.
// Used to sanitize user-supplied plain-text fields (names, etc.) before
// they are stored, preventing stored XSS when the values are rendered in
// pages, emails, or exports.
func StripHTMLTags(input string) string {
	// Remove script and style tags with content
	reScript := regexp.MustCompile(`(?s)<script.*?</script>`)
	reStyle := regexp.MustCompile(`(?s)<style.*?</style>`)
	text := reScript.ReplaceAllString(input, " ")
	text = reStyle.ReplaceAllString(text, " ")

	// Remove HTML tags
	reTag := regexp.MustCompile(`<[^>]+>`)
	text = reTag.ReplaceAllString(text, " ")

	// Replace multiple spaces with single space
	reSpace := regexp.MustCompile(`\s+`)
	text = reSpace.ReplaceAllString(text, " ")

	// Decode HTML entities (&nbsp; -> space, &amp; -> &, etc.)
	text = html.UnescapeString(text)

	return strings.TrimSpace(text)
}
