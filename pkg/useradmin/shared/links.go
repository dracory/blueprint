package shared

import (
	"net/url"
	"strings"
)

// Links generates URLs for user admin controllers
type Links struct {
	BaseURL string
}

// NewLinks creates a new Links instance
func NewLinks(baseURL string) *Links {
	if baseURL == "" {
		baseURL = "/admin/users"
	}
	return &Links{BaseURL: strings.TrimSuffix(baseURL, "/")}
}

// buildURL builds URL with controller parameter
func (l *Links) buildURL(controller string, params map[string]string) string {
	if params == nil {
		params = map[string]string{}
	}
	params["controller"] = controller

	q := url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}

	if len(q) > 0 {
		return l.BaseURL + "?" + q.Encode()
	}
	return l.BaseURL
}

// Home returns URL for user manager (default)
func (l *Links) Home(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_USER_MANAGER, p)
}

// UserManager returns URL for user manager
func (l *Links) UserManager(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_USER_MANAGER, p)
}

// UserCreate returns URL for user create
func (l *Links) UserCreate(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_USER_CREATE, p)
}

// UserDelete returns URL for user delete
func (l *Links) UserDelete(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_USER_DELETE, p)
}

// UserUpdate returns URL for user update
func (l *Links) UserUpdate(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_USER_UPDATE, p)
}

// UserImpersonate returns URL for user impersonate
func (l *Links) UserImpersonate(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_USER_IMPERSONATE, p)
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
