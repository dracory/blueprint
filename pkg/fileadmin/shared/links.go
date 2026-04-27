package shared

import (
	"fmt"
	"strings"
)

// Links generates URLs for file admin controllers
type Links struct {
	BaseURL string
}

// NewLinks creates a new Links instance
func NewLinks(baseURL string) *Links {
	if baseURL == "" {
		baseURL = "/admin/file-manager"
	}
	return &Links{BaseURL: strings.TrimSuffix(baseURL, "/")}
}

// buildURL builds URL with controller parameter
func (l *Links) buildURL(controller string, params map[string]string) string {
	if params == nil {
		params = map[string]string{}
	}
	params["controller"] = controller

	queryParts := []string{}
	for k, v := range params {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", k, v))
	}

	if len(queryParts) > 0 {
		return l.BaseURL + "?" + strings.Join(queryParts, "&")
	}
	return l.BaseURL
}

// FileManager returns URL for file manager
func (l *Links) FileManager(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_FILE_MANAGER, p)
}

// mergeParams merges multiple param maps
func mergeParams(params ...map[string]string) map[string]string {
	result := map[string]string{}
	for _, p := range params {
		for k, v := range p {
			result[k] = v
		}
	}
	return result
}
