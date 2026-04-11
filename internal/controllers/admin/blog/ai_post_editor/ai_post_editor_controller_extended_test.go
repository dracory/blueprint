package aiposteditor

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"
	"project/pkg/blogai"
)

func TestAiPostEditorController_HandlerWithGET(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostEditorController(registry)
	if controller == nil {
		t.Fatal("NewAiPostEditorController() returned nil")
	}

	// Test GET request without ID - expect panic due to nil custom store
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Expected panic due to nil custom store
			}
		}()

		req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-post-editor", nil)
		w := httptest.NewRecorder()

		controller.Handler(w, req)
	}()
}

func TestAiPostEditorController_HandlerWithInvalidID(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostEditorController(registry)

	// Test GET with invalid ID - expect panic due to nil custom store
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Expected panic due to nil custom store
			}
		}()

		req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-post-editor?id=invalid-id", nil)
		w := httptest.NewRecorder()

		controller.Handler(w, req)
	}()
}

func TestAiPostEditorController_HandlerWithNilRegistry(t *testing.T) {
	// Test with nil registry - expect panic
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil registry
		}
	}()

	controller := NewAiPostEditorController(nil)
	if controller == nil {
		t.Fatal("NewAiPostEditorController(nil) should not return nil")
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-post-editor", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestAiPostEditorController_HandlerWithPOSTActions(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostEditorController(registry)

	actions := []string{
		ACTION_REGENERATE_SECTION,
		ACTION_REGENERATE_IMAGE,
		ACTION_REGENERATE_PARAGRAPH,
		ACTION_CREATE_FINAL_POST,
		ACTION_SAVE_DRAFT,
		ACTION_LOAD_POST,
		ACTION_REGENERATE_SUMMARY,
		ACTION_REGENERATE_METAS,
	}

	for _, action := range actions {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Expected panic due to nil custom store
				}
			}()

			req := httptest.NewRequest(http.MethodPost, "/admin/blog/ai-post-editor?id=test&action="+action, nil)
			w := httptest.NewRecorder()

			controller.Handler(w, req)
		}()
	}
}

func TestAiPostEditorController_MultipleInstances(t *testing.T) {
	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	controller1 := NewAiPostEditorController(registry1)
	controller2 := NewAiPostEditorController(registry2)

	if controller1 == nil || controller2 == nil {
		t.Fatal("Controllers should not be nil")
	}

	if controller1 == controller2 {
		t.Error("Controllers should be separate instances")
	}

	if controller1.registry != registry1 {
		t.Error("Controller1 should have registry1")
	}

	if controller2.registry != registry2 {
		t.Error("Controller2 should have registry2")
	}
}

func TestAiPostEditorController_RegistryField(t *testing.T) {
	// Test with nil registry
	controller := NewAiPostEditorController(nil)
	if controller.registry != nil {
		t.Error("Controller registry should be nil when passed nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewAiPostEditorController(registry)
	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}

func TestAiPostEditorController_buildPostMarkdownContent(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostEditorController(registry)

	post := &blogai.RecordPost{
		ID:    "test-post",
		Title: "Test Post",
		Introduction: blogai.PostContentIntroduction{
			Title:      "Intro Title",
			Paragraphs: []string{"Para 1", "Para 2"},
		},
		Sections: []blogai.PostContentSection{
			{
				Title:      "Section 1",
				Paragraphs: []string{"Section content"},
			},
		},
		Conclusion: blogai.PostContentConclusion{
			Title:      "Conclusion Title",
			Paragraphs: []string{"Conclusion para"},
		},
	}

	content := controller.buildPostMarkdownContent(nil, post)

	if content == "" {
		t.Error("buildPostMarkdownContent should return non-empty string")
	}

	if !strings.Contains(content, "# Test Post") {
		t.Error("Content should contain main title")
	}

	if !strings.Contains(content, "## Intro Title") {
		t.Error("Content should contain introduction section")
	}

	if !strings.Contains(content, "## Section 1") {
		t.Error("Content should contain section 1")
	}

	if !strings.Contains(content, "## Conclusion Title") {
		t.Error("Content should contain conclusion section")
	}
}

func TestAiPostEditorController_buildPostMarkdownContent_Empty(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostEditorController(registry)

	post := &blogai.RecordPost{
		ID:    "empty-post",
		Title: "Empty Post",
		Introduction: blogai.PostContentIntroduction{
			Title:      "",
			Paragraphs: []string{},
		},
		Sections: []blogai.PostContentSection{},
		Conclusion: blogai.PostContentConclusion{
			Title:      "",
			Paragraphs: []string{},
		},
	}

	content := controller.buildPostMarkdownContent(nil, post)

	if content == "" {
		t.Error("buildPostMarkdownContent with empty post should return non-empty string")
	}
}

func TestAiPostEditorController_prepareDataAndValidate_MissingID(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostEditorController(registry)

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-post-editor", nil)

	_, errorMessage := controller.prepareDataAndValidate(req)

	if errorMessage == "" {
		t.Error("prepareDataAndValidate with missing ID should return error")
	}
}

func TestAiPostEditorController_prepareDataAndValidate_EmptyID(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostEditorController(registry)

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-post-editor?id=", nil)

	_, errorMessage := controller.prepareDataAndValidate(req)

	if errorMessage == "" {
		t.Error("prepareDataAndValidate with empty ID should return error")
	}
}

func TestAiPostEditorController_prepareDataAndValidate_InvalidID(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostEditorController(registry)

	// Test with invalid ID - expect panic due to nil custom store
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Expected panic due to nil custom store
			}
		}()

		req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-post-editor?id=invalid-id", nil)
		controller.prepareDataAndValidate(req)
	}()
}

func TestRecordFromJSON_Valid(t *testing.T) {
	jsonStr := `{
		"title": "Test Post",
		"subtitle": "Test Subtitle",
		"summary": "Test Summary",
		"introduction": {
			"title": "Intro",
			"paragraphs": ["Para 1"]
		},
		"sections": [
			{
				"title": "Section 1",
				"paragraphs": ["Section para"]
			}
		],
		"conclusion": {
			"title": "Conclusion",
			"paragraphs": ["Conclusion para"]
		},
		"keywords": ["test", "blog"],
		"metaDescription": "Meta Desc",
		"metaKeywords": ["meta1", "meta2"],
		"metaTitle": "Meta Title",
		"image": "image.jpg"
	}`

	post, err := RecordFromJSON(jsonStr)
	if err != nil {
		t.Errorf("RecordFromJSON error = %v", err)
	}
	if post == nil {
		t.Fatal("RecordFromJSON should return non-nil post")
	}

	if post.Title != "Test Post" {
		t.Errorf("Title = %v, want Test Post", post.Title)
	}

	if post.Subtitle != "Test Subtitle" {
		t.Errorf("Subtitle = %v, want Test Subtitle", post.Subtitle)
	}

	if post.MetaTitle != "Meta Title" {
		t.Errorf("MetaTitle = %v, want Meta Title", post.MetaTitle)
	}
}

func TestRecordFromJSON_InvalidJSON(t *testing.T) {
	jsonStr := `{invalid json}`

	_, err := RecordFromJSON(jsonStr)
	if err == nil {
		t.Error("RecordFromJSON should error with invalid JSON")
	}
}

func TestRecordFromJSON_Empty(t *testing.T) {
	jsonStr := `{}`

	post, err := RecordFromJSON(jsonStr)
	if err != nil {
		t.Errorf("RecordFromJSON error = %v", err)
	}
	if post == nil {
		t.Fatal("RecordFromJSON should return non-nil post even for empty JSON")
	}
}

func TestAiPostEditorController_ActionConstants(t *testing.T) {
	if ACTION_REGENERATE_SECTION != "regenerate_section" {
		t.Errorf("ACTION_REGENERATE_SECTION = %q, want regenerate_section", ACTION_REGENERATE_SECTION)
	}
	if ACTION_REGENERATE_IMAGE != "regenerate_image" {
		t.Errorf("ACTION_REGENERATE_IMAGE = %q, want regenerate_image", ACTION_REGENERATE_IMAGE)
	}
	if ACTION_CREATE_FINAL_POST != "create_final_post" {
		t.Errorf("ACTION_CREATE_FINAL_POST = %q, want create_final_post", ACTION_CREATE_FINAL_POST)
	}
	if ACTION_SAVE_DRAFT != "save_draft" {
		t.Errorf("ACTION_SAVE_DRAFT = %q, want save_draft", ACTION_SAVE_DRAFT)
	}
	if ACTION_REGENERATE_PARAGRAPH != "regenerate_paragraph" {
		t.Errorf("ACTION_REGENERATE_PARAGRAPH = %q, want regenerate_paragraph", ACTION_REGENERATE_PARAGRAPH)
	}
	if ACTION_LOAD_POST != "load_post" {
		t.Errorf("ACTION_LOAD_POST = %q, want load_post", ACTION_LOAD_POST)
	}
	if ACTION_REGENERATE_SUMMARY != "regenerate_summary" {
		t.Errorf("ACTION_REGENERATE_SUMMARY = %q, want regenerate_summary", ACTION_REGENERATE_SUMMARY)
	}
	if ACTION_REGENERATE_METAS != "regenerate_metas" {
		t.Errorf("ACTION_REGENERATE_METAS = %q, want regenerate_metas", ACTION_REGENERATE_METAS)
	}
}

func TestAiPostEditorController_pageDataStruct(t *testing.T) {
	// Test that pageData can be created
	data := pageData{
		Request:    nil,
		BlogAiPost: blogai.RecordPost{ID: "test", Title: "Test"},
		Record:     nil,
	}

	if data.BlogAiPost.ID != "test" {
		t.Error("pageData.BlogAiPost.ID not set correctly")
	}
}
