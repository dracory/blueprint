package social

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	// emailRegex matches most common email formats
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// phoneRegex matches international phone numbers with optional +, spaces, dashes, dots, and parentheses
	// Supports formats like: +1-555-123-4567, +44 20 7946 0958, (555) 123-4567, 555.123.4567
	phoneRegex = regexp.MustCompile(`^[\+]?[0-9]{0,4}[-\s\.]?[(]?[0-9]{1,4}[)]?[-\s\.]?[0-9]{1,4}[-\s\.]?[0-9]{1,9}$`)
)

// ValidateURL checks if the given string is a valid absolute URL
// Returns true for valid http:// or https:// URLs
func ValidateURL(input string) bool {
	if input == "" {
		return false
	}
	u, err := url.Parse(input)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

// ValidateEmail checks if the given string is a valid email address
func ValidateEmail(input string) bool {
	if input == "" {
		return false
	}
	return emailRegex.MatchString(input)
}

// ValidatePhone checks if the given string is a valid phone number
// Accepts international formats with optional +, spaces, dashes, and parentheses
func ValidatePhone(input string) bool {
	if input == "" {
		return false
	}
	// Remove common formatting characters for additional validation
	cleaned := strings.ReplaceAll(input, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// Must have at least 7 digits and no more than 15 (E.164 standard)
	if len(cleaned) < 7 || len(cleaned) > 15 {
		return false
	}

	// Must start with + or digit
	if !strings.HasPrefix(cleaned, "+") && (cleaned[0] < '0' || cleaned[0] > '9') {
		return false
	}

	return phoneRegex.MatchString(input)
}
