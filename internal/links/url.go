package links

import (
	baseurl "github.com/dracory/base/url"
)

// URL returns the full URL for a given path with optional query parameters
func URL(path string, params map[string]string) string {
	return baseurl.BuildURL(path, params)
}
