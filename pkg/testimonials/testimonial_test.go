package testimonials

import (
	"testing"
)

func TestNewTestimonial(t *testing.T) {
	testimonial := NewTestimonial()
	if testimonial == nil {
		t.Error("NewTestimonial() returned nil")
	}
}

func TestConstants(t *testing.T) {
	if ENTITY_TYPE != "testimonial" {
		t.Errorf("ENTITY_TYPE = %q, want %q", ENTITY_TYPE, "testimonial")
	}

	if FIELD_DATE != "date" {
		t.Errorf("FIELD_DATE = %q, want %q", FIELD_DATE, "date")
	}

	if FIELD_FIRST_NAME != "first_name" {
		t.Errorf("FIELD_FIRST_NAME = %q, want %q", FIELD_FIRST_NAME, "first_name")
	}

	if FIELD_ID != "id" {
		t.Errorf("FIELD_ID = %q, want %q", FIELD_ID, "id")
	}

	if FIELD_IMAGE_URL != "image_url" {
		t.Errorf("FIELD_IMAGE_URL = %q, want %q", FIELD_IMAGE_URL, "image_url")
	}

	if FIELD_JOB_TITLE != "job_title" {
		t.Errorf("FIELD_JOB_TITLE = %q, want %q", FIELD_JOB_TITLE, "job_title")
	}

	if FIELD_LAST_NAME != "last_name" {
		t.Errorf("FIELD_LAST_NAME = %q, want %q", FIELD_LAST_NAME, "last_name")
	}

	if FIELD_QUOTE != "quote" {
		t.Errorf("FIELD_QUOTE = %q, want %q", FIELD_QUOTE, "quote")
	}

	if FIELD_STATUS != "status" {
		t.Errorf("FIELD_STATUS = %q, want %q", FIELD_STATUS, "status")
	}
}
