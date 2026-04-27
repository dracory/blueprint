package shared

import (
	"fmt"
	"strings"
)

// Links generates URLs for blog admin controllers
type Links struct {
	BaseURL string
}

// NewLinks creates a new Links instance
func NewLinks(baseURL string) *Links {
	if baseURL == "" {
		baseURL = "/admin/blog"
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

// Home returns URL for home/post manager
func (l *Links) Home(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_POST_MANAGER, p)
}

// PostCreate returns URL for post create
func (l *Links) PostCreate(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_POST_CREATE, p)
}

// PostDelete returns URL for post delete
func (l *Links) PostDelete(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_POST_DELETE, p)
}

// PostManager returns URL for post manager
func (l *Links) PostManager(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_POST_MANAGER, p)
}

// PostUpdate returns URL for post update
func (l *Links) PostUpdate(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_POST_UPDATE, p)
}

// PostUpdateV1 returns URL for post update v1
func (l *Links) PostUpdateV1(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_POST_UPDATE_V1, p)
}

// BlogSettings returns URL for blog settings
func (l *Links) BlogSettings(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_BLOG_SETTINGS, p)
}

// AiTools returns URL for AI tools
func (l *Links) AiTools(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_AI_TOOLS, p)
}

// AiPostContentUpdate returns URL for AI post content update
func (l *Links) AiPostContentUpdate(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_AI_POST_CONTENT_UPDATE, p)
}

// AiPostGenerator returns URL for AI post generator
func (l *Links) AiPostGenerator(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_AI_POST_GENERATOR, p)
}

// AiTitleGenerator returns URL for AI title generator
func (l *Links) AiTitleGenerator(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_AI_TITLE_GENERATOR, p)
}

// AiPostEditor returns URL for AI post editor
func (l *Links) AiPostEditor(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_AI_POST_EDITOR, p)
}

// AiTest returns URL for AI test
func (l *Links) AiTest(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_AI_TEST, p)
}

// Dashboard returns URL for dashboard
func (l *Links) Dashboard(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_DASHBOARD, p)
}

// CategoryManager returns URL for category manager
func (l *Links) CategoryManager(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_CATEGORY_MANAGER, p)
}

// TagManager returns URL for tag manager
func (l *Links) TagManager(params ...map[string]string) string {
	p := mergeParams(params...)
	return l.buildURL(CONTROLLER_TAG_MANAGER, p)
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
