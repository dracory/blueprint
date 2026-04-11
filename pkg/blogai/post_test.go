package blogai

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewRecordPostFromMap(t *testing.T) {
	// Test with complete data
	data := map[string]any{
		"id":       "post-123",
		"title":    "Test Post",
		"status":   POST_STATUS_PUBLISHED,
		"subtitle": "A test subtitle",
		"summary":  "Test summary",
		"introduction": map[string]any{
			"title":      "Introduction",
			"paragraphs": []any{"Paragraph 1", "Paragraph 2"},
		},
		"sections": []any{
			map[string]any{
				"title":      "Section 1",
				"paragraphs": []any{"Section content"},
			},
		},
		"conclusion": map[string]any{
			"title":      "Conclusion",
			"paragraphs": []any{"Concluding thoughts"},
		},
		"keywords":         []any{"test", "blog", "ai"},
		"created_at":       "2024-01-01T00:00:00Z",
		"updated_at":       "2024-01-02T00:00:00Z",
		"meta_title":       "Meta Title",
		"meta_description": "Meta Description",
		"meta_keywords":    []any{"meta", "keywords"},
		"image":            "https://example.com/image.jpg",
	}

	post, err := newRecordPostFromMap(data)
	if err != nil {
		t.Errorf("newRecordPostFromMap() error = %v", err)
	}

	if post.ID != "post-123" {
		t.Errorf("ID = %v, want post-123", post.ID)
	}
	if post.Title != "Test Post" {
		t.Errorf("Title = %v, want Test Post", post.Title)
	}
	if post.Status != POST_STATUS_PUBLISHED {
		t.Errorf("Status = %v, want %s", post.Status, POST_STATUS_PUBLISHED)
	}
	if post.Subtitle != "A test subtitle" {
		t.Errorf("Subtitle = %v, want A test subtitle", post.Subtitle)
	}
	if post.Summary != "Test summary" {
		t.Errorf("Summary = %v, want Test summary", post.Summary)
	}
	if post.Introduction.Title != "Introduction" {
		t.Errorf("Introduction.Title = %v, want Introduction", post.Introduction.Title)
	}
	if len(post.Introduction.Paragraphs) != 2 {
		t.Errorf("Introduction.Paragraphs length = %v, want 2", len(post.Introduction.Paragraphs))
	}
	if len(post.Sections) != 1 {
		t.Errorf("Sections length = %v, want 1", len(post.Sections))
	}
	if post.Conclusion.Title != "Conclusion" {
		t.Errorf("Conclusion.Title = %v, want Conclusion", post.Conclusion.Title)
	}
	if len(post.Keywords) != 3 {
		t.Errorf("Keywords length = %v, want 3", len(post.Keywords))
	}
	if post.Image != "https://example.com/image.jpg" {
		t.Errorf("Image = %v, want https://example.com/image.jpg", post.Image)
	}
}

func TestNewRecordPostFromMap_MinimalData(t *testing.T) {
	// Test with minimal required fields only
	data := map[string]any{
		"id":    "minimal-post",
		"title": "Minimal Post",
	}

	post, err := newRecordPostFromMap(data)
	if err != nil {
		t.Errorf("newRecordPostFromMap() error = %v", err)
	}

	if post.ID != "minimal-post" {
		t.Errorf("ID = %v, want minimal-post", post.ID)
	}
	if post.Title != "Minimal Post" {
		t.Errorf("Title = %v, want Minimal Post", post.Title)
	}
	// Check defaults
	if post.Status != "" {
		t.Errorf("Status = %v, want empty", post.Status)
	}
	if post.MetaTitle != "Minimal Post" {
		t.Errorf("MetaTitle = %v, want Minimal Post (default to title)", post.MetaTitle)
	}
	if post.CreatedAt == "" {
		t.Error("CreatedAt should not be empty")
	}
	if post.UpdatedAt == "" {
		t.Error("UpdatedAt should not be empty")
	}
}

func TestNewRecordPostFromMap_MissingID(t *testing.T) {
	// Test with missing ID - should error
	data := map[string]any{
		"title": "No ID Post",
	}

	_, err := newRecordPostFromMap(data)
	if err == nil {
		t.Error("newRecordPostFromMap() should error when id is missing")
	}
}

func TestNewRecordPostFromMap_MissingTitle(t *testing.T) {
	// Test with missing title - should error
	data := map[string]any{
		"id": "no-title-post",
	}

	_, err := newRecordPostFromMap(data)
	if err == nil {
		t.Error("newRecordPostFromMap() should error when title is missing")
	}
}

func TestNewRecordPostFromMap_BackwardCompatibility(t *testing.T) {
	// Test backward compatibility with old format using "content" field
	data := map[string]any{
		"id":    "compat-post",
		"title": "Compat Post",
		"introduction": map[string]any{
			"content": "Old style introduction",
		},
		"sections": []any{
			map[string]any{
				"title":   "Section",
				"content": "Old style section content",
			},
		},
		"conclusion": map[string]any{
			"content": "Old style conclusion",
		},
	}

	post, err := newRecordPostFromMap(data)
	if err != nil {
		t.Errorf("newRecordPostFromMap() error = %v", err)
	}

	if len(post.Introduction.Paragraphs) != 1 || post.Introduction.Paragraphs[0] != "Old style introduction" {
		t.Error("Introduction paragraphs not parsed from old format")
	}
	if len(post.Sections) != 1 || len(post.Sections[0].Paragraphs) != 1 {
		t.Error("Section paragraphs not parsed from old format")
	}
	if len(post.Conclusion.Paragraphs) != 1 || post.Conclusion.Paragraphs[0] != "Old style conclusion" {
		t.Error("Conclusion paragraphs not parsed from old format")
	}
}

func TestRecordPostToJSON(t *testing.T) {
	post := RecordPost{
		ID:              "json-post",
		Title:           "JSON Post",
		Status:          POST_STATUS_DRAFT,
		Subtitle:        "Subtitle",
		Summary:         "Summary",
		Keywords:        []string{"test", "json"},
		MetaDescription: "Meta desc",
		MetaKeywords:    []string{"meta1", "meta2"},
		MetaTitle:       "Meta Title",
		Image:           "image.jpg",
		CreatedAt:       time.Now().Format(time.RFC3339),
		UpdatedAt:       time.Now().Format(time.RFC3339),
		Introduction:    PostContentIntroduction{Title: "Intro", Paragraphs: []string{"p1"}},
		Sections:        []PostContentSection{{Title: "Section", Paragraphs: []string{"p2"}}},
		Conclusion:      PostContentConclusion{Title: "Outro", Paragraphs: []string{"p3"}},
	}

	jsonStr := post.ToJSON()
	if jsonStr == "" {
		t.Error("ToJSON() should not return empty string")
	}

	// Verify it's valid JSON
	var decoded map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &decoded); err != nil {
		t.Errorf("ToJSON() returned invalid JSON: %v", err)
	}

	// Verify some fields
	if decoded["id"] != "json-post" {
		t.Error("JSON id field mismatch")
	}
	if decoded["title"] != "JSON Post" {
		t.Error("JSON title field mismatch")
	}
}

func TestRecordPostToJSON_ValidObject(t *testing.T) {
	// Test ToJSON returns valid JSON for a normal object
	post := RecordPost{
		ID:    "test",
		Title: "Test",
	}

	jsonStr := post.ToJSON()
	if jsonStr == "" {
		t.Error("ToJSON() should return valid JSON for normal objects")
	}
}

func TestPostStructs(t *testing.T) {
	// Test that structs can be created and fields set
	intro := PostContentIntroduction{
		Title:      "Intro Title",
		Paragraphs: []string{"para1", "para2"},
	}
	if intro.Title != "Intro Title" {
		t.Error("PostContentIntroduction.Title not set correctly")
	}
	if len(intro.Paragraphs) != 2 {
		t.Error("PostContentIntroduction.Paragraphs not set correctly")
	}

	section := PostContentSection{
		Title:      "Section Title",
		Paragraphs: []string{"section para"},
	}
	if section.Title != "Section Title" {
		t.Error("PostContentSection.Title not set correctly")
	}

	conclusion := PostContentConclusion{
		Title:      "Conclusion Title",
		Paragraphs: []string{"conclusion para"},
	}
	if conclusion.Title != "Conclusion Title" {
		t.Error("PostContentConclusion.Title not set correctly")
	}
}

func TestNewRecordPostFromMap_EmptySlices(t *testing.T) {
	// Test with empty or nil slices
	data := map[string]any{
		"id":       "empty-post",
		"title":    "Empty Post",
		"keywords": []any{},
	}

	post, err := newRecordPostFromMap(data)
	if err != nil {
		t.Errorf("newRecordPostFromMap() error = %v", err)
	}

	if len(post.Keywords) != 0 {
		t.Errorf("Keywords should be empty, got %v", len(post.Keywords))
	}
}
