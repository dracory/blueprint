package links

import (
	"net/url"
	"os"
)

// RootURL returns a URL to the current website
func RootURL() string {
	appURL := os.Getenv("APP_URL")
	if os.Getenv("APP_ENV") == "testing" {
		return ""
	}
	return appURL
}

func query(queryData map[string]string) string {
	queryString := ""

	if len(queryData) > 0 {
		v := url.Values{}
		for key, value := range queryData {
			v.Set(key, value)
		}
		queryString += "?" + httpBuildQuery(v)
	}

	return queryString
}

func httpBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}
