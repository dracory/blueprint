package aiposteditor

import (
	"encoding/json"
	"testing"

	"project/pkg/blogai"
)

// TestNewAiPostEditorController tests the constructor
func TestNewAiPostEditorController(t *testing.T) {
	t.Parallel()

	controller := NewAiPostEditorController(nil)
	if controller == nil {
		t.Error("Expected controller to be non-nil")
	}
	if controller.registry != nil {
		t.Error("Expected registry to be nil when passed nil")
	}
}

// TestRecordFromJSON_ValidJSON tests parsing valid JSON
func TestRecordFromJSON_ValidJSON(t *testing.T) {
	t.Parallel()

	validJSON := `{
		"title": "Test Title",
		"subtitle": "Test Subtitle",
		"summary": "Test Summary",
		"introduction": {
			"title": "Intro Title",
			"paragraphs": ["Para 1", "Para 2"]
		},
		"sections": [
			{
				"title": "Section 1",
				"paragraphs": ["Section para 1"]
			}
		],
		"conclusion": {
			"title": "Conclusion Title",
			"paragraphs": ["Conclusion para"]
		},
		"keywords": ["key1", "key2"],
		"metaDescription": "Meta desc",
		"metaKeywords": ["meta1"],
		"metaTitle": "Meta Title",
		"image": "http://example.com/image.jpg"
	}`

	record, err := RecordFromJSON(validJSON)
	if err != nil {
		t.Errorf("RecordFromJSON() unexpected error: %v", err)
	}
	if record == nil {
		t.Error("RecordFromJSON() returned nil record")
		return
	}

	// Verify fields
	if record.Title != "Test Title" {
		t.Errorf("Expected Title 'Test Title', got: %s", record.Title)
	}
	if record.Subtitle != "Test Subtitle" {
		t.Errorf("Expected Subtitle 'Test Subtitle', got: %s", record.Subtitle)
	}
	if record.Summary != "Test Summary" {
		t.Errorf("Expected Summary 'Test Summary', got: %s", record.Summary)
	}
	if record.MetaDescription != "Meta desc" {
		t.Errorf("Expected MetaDescription 'Meta desc', got: %s", record.MetaDescription)
	}
	if record.Image != "http://example.com/image.jpg" {
		t.Errorf("Expected Image 'http://example.com/image.jpg', got: %s", record.Image)
	}
	if len(record.Keywords) != 2 {
		t.Errorf("Expected 2 keywords, got: %d", len(record.Keywords))
	}
	if len(record.Sections) != 1 {
		t.Errorf("Expected 1 section, got: %d", len(record.Sections))
	}
}

// TestRecordFromJSON_InvalidJSON tests parsing invalid JSON
func TestRecordFromJSON_InvalidJSON(t *testing.T) {
	t.Parallel()

	invalidJSON := `{
		"title": "Test",
		"invalid json here
	}`

	record, err := RecordFromJSON(invalidJSON)
	if err == nil {
		t.Error("RecordFromJSON() expected error for invalid JSON")
	}
	if record != nil {
		t.Error("RecordFromJSON() should return nil for invalid JSON")
	}
}

// TestRecordFromJSON_EmptyJSON tests parsing empty JSON object
func TestRecordFromJSON_EmptyJSON(t *testing.T) {
	t.Parallel()

	record, err := RecordFromJSON("{}")
	if err != nil {
		t.Errorf("RecordFromJSON() unexpected error for empty JSON: %v", err)
	}
	if record == nil {
		t.Error("RecordFromJSON() should return non-nil for empty JSON")
		return
	}
	// All fields should be zero values
	if record.Title != "" {
		t.Errorf("Expected empty Title, got: %s", record.Title)
	}
}

// TestRecordFromJSON_NilJSON tests parsing nil-like JSON
func TestRecordFromJSON_NilJSON(t *testing.T) {
	t.Parallel()

	record, err := RecordFromJSON("null")
	if err != nil {
		t.Errorf("RecordFromJSON() unexpected error for null JSON: %v", err)
	}
	// null unmarshals to an empty struct, not nil
	if record == nil {
		t.Error("RecordFromJSON() for null should not return nil")
	}
}

// TestRecordFromJSON_MissingOptionalFields tests JSON with only required fields
func TestRecordFromJSON_MissingOptionalFields(t *testing.T) {
	t.Parallel()

	minimalJSON := `{
		"title": "Minimal Post",
		"introduction": {
			"title": "",
			"paragraphs": []
		},
		"sections": [],
		"conclusion": {
			"title": "",
			"paragraphs": []
		}
	}`

	record, err := RecordFromJSON(minimalJSON)
	if err != nil {
		t.Errorf("RecordFromJSON() unexpected error: %v", err)
	}
	if record == nil {
		t.Fatal("RecordFromJSON() returned nil")
	}

	if record.Title != "Minimal Post" {
		t.Errorf("Expected Title 'Minimal Post', got: %s", record.Title)
	}
	// Optional fields should be empty
	if record.Subtitle != "" {
		t.Errorf("Expected empty Subtitle, got: %s", record.Subtitle)
	}
	if record.Image != "" {
		t.Errorf("Expected empty Image, got: %s", record.Image)
	}
}

// TestActionConstants tests action constant values
func TestActionConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"ACTION_REGENERATE_SECTION", ACTION_REGENERATE_SECTION, "regenerate_section"},
		{"ACTION_REGENERATE_IMAGE", ACTION_REGENERATE_IMAGE, "regenerate_image"},
		{"ACTION_CREATE_FINAL_POST", ACTION_CREATE_FINAL_POST, "create_final_post"},
		{"ACTION_SAVE_DRAFT", ACTION_SAVE_DRAFT, "save_draft"},
		{"ACTION_REGENERATE_PARAGRAPH", ACTION_REGENERATE_PARAGRAPH, "regenerate_paragraph"},
		{"ACTION_LOAD_POST", ACTION_LOAD_POST, "load_post"},
		{"ACTION_REGENERATE_SUMMARY", ACTION_REGENERATE_SUMMARY, "regenerate_summary"},
		{"ACTION_REGENERATE_METAS", ACTION_REGENERATE_METAS, "regenerate_metas"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s to be %q, got: %q", tt.name, tt.expected, tt.constant)
			}
		})
	}
}

// TestRecordFromJSON_ArrayFields tests parsing arrays correctly
func TestRecordFromJSON_ArrayFields(t *testing.T) {
	t.Parallel()

	jsonWithArrays := `{
		"title": "Array Test",
		"keywords": ["a", "b", "c", "d"],
		"metaKeywords": ["meta-a", "meta-b"],
		"introduction": {"title": "", "paragraphs": ["p1", "p2", "p3"]},
		"sections": [
			{"title": "S1", "paragraphs": ["s1p1"]},
			{"title": "S2", "paragraphs": ["s2p1", "s2p2"]}
		],
		"conclusion": {"title": "", "paragraphs": []}
	}`

	record, err := RecordFromJSON(jsonWithArrays)
	if err != nil {
		t.Fatalf("RecordFromJSON() unexpected error: %v", err)
	}
	if record == nil {
		t.Fatal("RecordFromJSON() returned nil")
	}

	// Check arrays
	if len(record.Keywords) != 4 {
		t.Errorf("Expected 4 keywords, got: %d", len(record.Keywords))
	}
	if len(record.MetaKeywords) != 2 {
		t.Errorf("Expected 2 metaKeywords, got: %d", len(record.MetaKeywords))
	}
	if len(record.Introduction.Paragraphs) != 3 {
		t.Errorf("Expected 3 intro paragraphs, got: %d", len(record.Introduction.Paragraphs))
	}
	if len(record.Sections) != 2 {
		t.Errorf("Expected 2 sections, got: %d", len(record.Sections))
	}
}

// TestRecordFromJSON_ComplexContent tests deeply nested structures
func TestRecordFromJSON_ComplexContent(t *testing.T) {
	t.Parallel()

	complexJSON := `{
		"title": "Complex Post",
		"subtitle": "A detailed subtitle",
		"summary": "Brief summary",
		"introduction": {
			"title": "Introduction",
			"paragraphs": [
				"First paragraph of introduction.",
				"Second paragraph with more details."
			]
		},
		"sections": [
			{
				"title": "First Section",
				"paragraphs": [
					"Content for first section."
				]
			},
			{
				"title": "Second Section",
				"paragraphs": [
					"First para of second section.",
					"Second para of second section."
				]
			},
			{
				"title": "Third Section",
				"paragraphs": [
					"Content for third section."
				]
			}
		],
		"conclusion": {
			"title": "Final Thoughts",
			"paragraphs": [
				"Concluding paragraph one.",
				"Concluding paragraph two."
			]
		},
		"keywords": ["test", "complex", "json"],
		"metaDescription": "A complex test post description",
		"metaKeywords": ["testing", "go"],
		"metaTitle": "Complex Test Post Meta Title",
		"image": "https://example.com/complex.jpg"
	}`

	record, err := RecordFromJSON(complexJSON)
	if err != nil {
		t.Fatalf("RecordFromJSON() unexpected error: %v", err)
	}
	if record == nil {
		t.Fatal("RecordFromJSON() returned nil")
	}

	// Verify structure
	if record.Title != "Complex Post" {
		t.Errorf("Expected title 'Complex Post', got: %s", record.Title)
	}

	// Verify nested structures
	if record.Introduction.Title != "Introduction" {
		t.Errorf("Expected intro title 'Introduction', got: %s", record.Introduction.Title)
	}
	if len(record.Introduction.Paragraphs) != 2 {
		t.Errorf("Expected 2 intro paragraphs, got: %d", len(record.Introduction.Paragraphs))
	}

	// Verify sections
	if len(record.Sections) != 3 {
		t.Errorf("Expected 3 sections, got: %d", len(record.Sections))
	}

	// Verify conclusion
	if record.Conclusion.Title != "Final Thoughts" {
		t.Errorf("Expected conclusion title 'Final Thoughts', got: %s", record.Conclusion.Title)
	}
	if len(record.Conclusion.Paragraphs) != 2 {
		t.Errorf("Expected 2 conclusion paragraphs, got: %d", len(record.Conclusion.Paragraphs))
	}
}

// TestRecordFromJSON_UnicodeContent tests unicode handling
func TestRecordFromJSON_UnicodeContent(t *testing.T) {
	t.Parallel()

	unicodeJSON := `{
		"title": "Unicode Test: 你好世界 🌍 émojis",
		"introduction": {
			"title": "Intro: ñoño",
			"paragraphs": ["段落1", "🎉 Party time 🎉"]
		},
		"sections": [],
		"conclusion": {
			"title": "Conclusión",
			"paragraphs": []
		},
		"keywords": ["unicode", "你好", "🌍"]
	}`

	record, err := RecordFromJSON(unicodeJSON)
	if err != nil {
		t.Fatalf("RecordFromJSON() unexpected error with unicode: %v", err)
	}
	if record == nil {
		t.Fatal("RecordFromJSON() returned nil")
	}

	if record.Title != "Unicode Test: 你好世界 🌍 émojis" {
		t.Errorf("Unicode title not preserved correctly: %s", record.Title)
	}

	if len(record.Keywords) != 3 {
		t.Errorf("Expected 3 unicode keywords, got: %d", len(record.Keywords))
	}
}

// TestRecordFromJSON_SpecialCharacters tests special characters in content
func TestRecordFromJSON_SpecialCharacters(t *testing.T) {
	t.Parallel()

	specialJSON := `{
		"title": "Special <>&\"' Characters",
		"metaDescription": "Description with \"quotes\" and 'apostrophes'",
		"introduction": {"title": "", "paragraphs": ["Line 1\nLine 2\tTabbed"]},
		"sections": [],
		"conclusion": {"title": "", "paragraphs": []}
	}`

	record, err := RecordFromJSON(specialJSON)
	if err != nil {
		t.Fatalf("RecordFromJSON() unexpected error: %v", err)
	}
	if record == nil {
		t.Fatal("RecordFromJSON() returned nil")
	}

	// Special characters should be preserved
	if record.Title != "Special <>&\"' Characters" {
		t.Errorf("Special characters not preserved: %s", record.Title)
	}
}

// TestAiPostEditorController_StructFields tests controller structure
func TestAiPostEditorController_StructFields(t *testing.T) {
	t.Parallel()

	// Verify struct fields exist and are accessible
	controller := NewAiPostEditorController(nil)

	// The struct should have a registry field
	// We can't directly test private fields, but we can verify the type
	if controller == nil {
		t.Fatal("NewAiPostEditorController() returned nil")
	}
}

// TestPageData_Struct tests pageData structure
func TestPageData_Struct(t *testing.T) {
	t.Parallel()

	// Test that pageData can be created
	data := pageData{}

	// Verify default zero values
	if data.Request != nil {
		t.Error("Request should be nil by default")
	}

	// Test setting values
	record := &blogai.RecordPost{}
	record.Title = "Test"

	// We can't directly set BlogAiPost since it's not exported
	// but we verified the struct can be instantiated
	_ = data
	_ = record
}

// TestRecordFromJSON_TypeCompatibility tests that the parsed types match blogai.RecordPost
func TestRecordFromJSON_TypeCompatibility(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"title": "Type Test",
		"introduction": {"title": "Intro", "paragraphs": ["p1"]},
		"sections": [{"title": "S1", "paragraphs": ["s1p1"]}],
		"conclusion": {"title": "Conclusion", "paragraphs": ["c1"]}
	}`

	record, err := RecordFromJSON(jsonData)
	if err != nil {
		t.Fatalf("RecordFromJSON() error: %v", err)
	}

	// Verify the returned type is *blogai.RecordPost
	var _ *blogai.RecordPost = record

	// Verify we can call methods on it if any
	// RecordPost should have ToJSON method
	jsonOutput := record.ToJSON()
	if jsonOutput == "" {
		t.Error("ToJSON() returned empty string")
	}

	// Verify the JSON can be parsed back
	var verify map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOutput), &verify); err != nil {
		t.Errorf("ToJSON() output is not valid JSON: %v", err)
	}

	if verify["title"] != "Type Test" {
		t.Error("ToJSON() output doesn't contain expected data")
	}
}
