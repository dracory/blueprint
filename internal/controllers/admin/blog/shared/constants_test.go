package shared

import (
	"testing"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"CONTROLLER_HOME", CONTROLLER_HOME},
		{"CONTROLLER_POST_CREATE", CONTROLLER_POST_CREATE},
		{"CONTROLLER_POST_DELETE", CONTROLLER_POST_DELETE},
		{"CONTROLLER_POST_MANAGER", CONTROLLER_POST_MANAGER},
		{"CONTROLLER_POST_UPDATE", CONTROLLER_POST_UPDATE},
		{"CONTROLLER_POST_UPDATE_V1", CONTROLLER_POST_UPDATE_V1},
		{"CONTROLLER_BLOG_SETTINGS", CONTROLLER_BLOG_SETTINGS},
		{"CONTROLLER_AI_TOOLS", CONTROLLER_AI_TOOLS},
		{"CONTROLLER_AI_POST_CONTENT_UPDATE", CONTROLLER_AI_POST_CONTENT_UPDATE},
		{"CONTROLLER_AI_POST_GENERATOR", CONTROLLER_AI_POST_GENERATOR},
		{"CONTROLLER_AI_TITLE_GENERATOR", CONTROLLER_AI_TITLE_GENERATOR},
		{"CONTROLLER_AI_POST_EDITOR", CONTROLLER_AI_POST_EDITOR},
		{"CONTROLLER_AI_TEST", CONTROLLER_AI_TEST},
		{"CONTROLLER_DASHBOARD", CONTROLLER_DASHBOARD},
		{"CONTROLLER_CATEGORY_MANAGER", CONTROLLER_CATEGORY_MANAGER},
		{"CONTROLLER_TAG_MANAGER", CONTROLLER_TAG_MANAGER},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("%s should not be empty", tt.name)
			}
		})
	}
}

func TestNewLinks(t *testing.T) {
	links := NewLinks()
	if links == nil {
		t.Error("NewLinks() should not return nil")
	}
}

func TestLinksMethods(t *testing.T) {
	links := NewLinks()

	tests := []struct {
		name   string
		method func() string
	}{
		{"Home", func() string { return links.Home() }},
		{"PostCreate", func() string { return links.PostCreate() }},
		{"PostDelete", func() string { return links.PostDelete() }},
		{"PostManager", func() string { return links.PostManager() }},
		{"PostUpdate", func() string { return links.PostUpdate() }},
		{"PostUpdateV1", func() string { return links.PostUpdateV1() }},
		{"BlogSettings", func() string { return links.BlogSettings() }},
		{"AiTools", func() string { return links.AiTools() }},
		{"AiPostContentUpdate", func() string { return links.AiPostContentUpdate() }},
		{"AiPostGenerator", func() string { return links.AiPostGenerator() }},
		{"AiTitleGenerator", func() string { return links.AiTitleGenerator() }},
		{"AiPostEditor", func() string { return links.AiPostEditor() }},
		{"AiTest", func() string { return links.AiTest() }},
		{"Dashboard", func() string { return links.Dashboard() }},
		{"CategoryManager", func() string { return links.CategoryManager() }},
		{"TagManager", func() string { return links.TagManager() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()
			if result == "" {
				t.Errorf("Links.%s() should not return empty string", tt.name)
			}
		})
	}
}
