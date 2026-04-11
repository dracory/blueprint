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

	// Test Home method
	result := links.Home()
	if result == "" {
		t.Error("Links.Home() should not return empty string")
	}

	// Test PostCreate method
	result = links.PostCreate()
	if result == "" {
		t.Error("Links.PostCreate() should not return empty string")
	}

	// Test PostDelete method
	result = links.PostDelete()
	if result == "" {
		t.Error("Links.PostDelete() should not return empty string")
	}

	// Test PostManager method
	result = links.PostManager()
	if result == "" {
		t.Error("Links.PostManager() should not return empty string")
	}

	// Test PostUpdate method
	result = links.PostUpdate()
	if result == "" {
		t.Error("Links.PostUpdate() should not return empty string")
	}

	// Test PostUpdateV1 method
	result = links.PostUpdateV1()
	if result == "" {
		t.Error("Links.PostUpdateV1() should not return empty string")
	}
}
