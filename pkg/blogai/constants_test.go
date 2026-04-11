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
