package blogai

import (
	"testing"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "POST_RECORD_TYPE", got: POST_RECORD_TYPE, want: "blogai_post"},
		{name: "POST_STATUS_PENDING", got: POST_STATUS_PENDING, want: "pending"},
		{name: "POST_STATUS_APPROVED", got: POST_STATUS_APPROVED, want: "approved"},
		{name: "POST_STATUS_REJECTED", got: POST_STATUS_REJECTED, want: "rejected"},
		{name: "POST_STATUS_DRAFT", got: POST_STATUS_DRAFT, want: "draft"},
		{name: "POST_STATUS_PUBLISHED", got: POST_STATUS_PUBLISHED, want: "published"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestRecordPostStruct(t *testing.T) {
	// Test RecordPost struct initialization
	post := RecordPost{
		Title:     "Test Title",
		Status:    POST_STATUS_DRAFT,
		Subtitle:  "Test Subtitle",
		Summary:   "Test Summary",
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-02T00:00:00Z",
		Image:     "https://example.com/image.jpg",
	}

	if post.Title != "Test Title" {
		t.Errorf("Title = %q, want %q", post.Title, "Test Title")
	}
	if post.Status != POST_STATUS_DRAFT {
		t.Errorf("Status = %q, want %q", post.Status, POST_STATUS_DRAFT)
	}
	if post.Subtitle != "Test Subtitle" {
		t.Errorf("Subtitle = %q, want %q", post.Subtitle, "Test Subtitle")
	}
}

func TestPostContentStructs(t *testing.T) {
	// Test PostContentIntroduction
	intro := PostContentIntroduction{
		Title:      "Introduction",
		Paragraphs: []string{"Para 1", "Para 2"},
	}

	if intro.Title != "Introduction" {
		t.Errorf("Introduction Title = %q, want %q", intro.Title, "Introduction")
	}
	if len(intro.Paragraphs) != 2 {
		t.Errorf("Introduction Paragraphs length = %d, want 2", len(intro.Paragraphs))
	}

	// Test PostContentSection
	section := PostContentSection{
		Title:      "Section 1",
		Paragraphs: []string{"Section para"},
	}

	if section.Title != "Section 1" {
		t.Errorf("Section Title = %q, want %q", section.Title, "Section 1")
	}

	// Test PostContentConclusion
	conclusion := PostContentConclusion{
		Title:      "Conclusion",
		Paragraphs: []string{"Final thought"},
	}

	if conclusion.Title != "Conclusion" {
		t.Errorf("Conclusion Title = %q, want %q", conclusion.Title, "Conclusion")
	}
}

func TestNewTitleGeneratorAgent(t *testing.T) {
	agent := NewTitleGeneratorAgent()
	if agent == nil {
		t.Error("NewTitleGeneratorAgent() should not return nil")
	}
}

func TestNewBlogWriterAgent(t *testing.T) {
	agent := NewBlogWriterAgent(nil)
	if agent == nil {
		t.Error("NewBlogWriterAgent() should not return nil")
	}
}

func TestNewTitleGeneratorAgentV1(t *testing.T) {
	agent := NewTitleGeneratorAgentV1()
	if agent == nil {
		t.Error("NewTitleGeneratorAgentV1() should not return nil")
	}
}
