package links

import (
	"net/url"
	"os"

	baseurl "github.com/dracory/base/url"
)

var initialized bool

// initializeURLBuilder initializes the base URL builder with environment variable
// This should be called before any URL building operations
func initializeURLBuilder() {
	if !initialized {
		appURL := os.Getenv("APP_URL")
		if os.Getenv("APP_ENV") == "testing" {
			appURL = "http://localhost:8080" // Set a default URL for testing
		}
		baseurl.SetDefaultURL(appURL)
		initialized = true
	}
}

// RootURL returns a URL to the current website
func RootURL() string {
	initializeURLBuilder()
	return baseurl.RootURL()
}

func query(queryData map[string]string) string {
	return baseurl.BuildQuery(queryData)
}

func httpBuildQuery(queryData url.Values) string {
	return baseurl.HttpBuildQuery(queryData)
}
