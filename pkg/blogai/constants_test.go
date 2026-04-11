package blogai

import (
	"testing"
)

func TestConstants(t *testing.T) {
	if POST_RECORD_TYPE != "blogai_post" {
		t.Errorf("POST_RECORD_TYPE = %q, want %q", POST_RECORD_TYPE, "blogai_post")
	}

	if POST_STATUS_PENDING != "pending" {
		t.Errorf("POST_STATUS_PENDING = %q, want %q", POST_STATUS_PENDING, "pending")
	}

	if POST_STATUS_APPROVED != "approved" {
		t.Errorf("POST_STATUS_APPROVED = %q, want %q", POST_STATUS_APPROVED, "approved")
	}

	if POST_STATUS_REJECTED != "rejected" {
		t.Errorf("POST_STATUS_REJECTED = %q, want %q", POST_STATUS_REJECTED, "rejected")
	}

	if POST_STATUS_DRAFT != "draft" {
		t.Errorf("POST_STATUS_DRAFT = %q, want %q", POST_STATUS_DRAFT, "draft")
	}

	if POST_STATUS_PUBLISHED != "published" {
		t.Errorf("POST_STATUS_PUBLISHED = %q, want %q", POST_STATUS_PUBLISHED, "published")
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
