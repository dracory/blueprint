package links

import (
	"net/url"
	"os"
	"sync"

	baseurl "github.com/dracory/base/url"
)

var initOnce sync.Once

// initializeURLBuilder initializes the base URL builder with environment variable
// This should be called before any URL building operations
func initializeURLBuilder() {
	initOnce.Do(func() {
		appURL := os.Getenv("APP_URL")
		if os.Getenv("APP_ENV") == "testing" {
			appURL = "http://localhost:8080" // Set a default URL for testing
		}
		baseurl.SetDefaultURL(appURL)
	})
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
