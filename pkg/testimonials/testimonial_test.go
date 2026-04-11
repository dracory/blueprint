package testimonials

import "testing"

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

func TestTestimonialSettersAndGetters(t *testing.T) {
	testimonial := NewTestimonial()

	// Test Date
	testimonial.SetDate("2024-01-01")
	if testimonial.Date() != "2024-01-01" {
		t.Errorf("Date() = %q, want %q", testimonial.Date(), "2024-01-01")
	}

	// Test FirstName
	testimonial.SetFirstName("John")
	if testimonial.FirstName() != "John" {
		t.Errorf("FirstName() = %q, want %q", testimonial.FirstName(), "John")
	}

	// Test ID
	testimonial.SetID("123")
	if testimonial.ID() != "123" {
		t.Errorf("ID() = %q, want %q", testimonial.ID(), "123")
	}

	// Test ImageUrl
	testimonial.SetImageUrl("https://example.com/image.jpg")
	if testimonial.ImageUrl() != "https://example.com/image.jpg" {
		t.Errorf("ImageUrl() = %q, want %q", testimonial.ImageUrl(), "https://example.com/image.jpg")
	}

	// Test JobTitle
	testimonial.SetJobTitle("Developer")
	if testimonial.JobTitle() != "Developer" {
		t.Errorf("JobTitle() = %q, want %q", testimonial.JobTitle(), "Developer")
	}

	// Test LastName
	testimonial.SetLastName("Doe")
	if testimonial.LastName() != "Doe" {
		t.Errorf("LastName() = %q, want %q", testimonial.LastName(), "Doe")
	}

	// Test Quote
	testimonial.SetQuote("This is a testimonial")
	if testimonial.Quote() != "This is a testimonial" {
		t.Errorf("Quote() = %q, want %q", testimonial.Quote(), "This is a testimonial")
	}

	// Test Status
	testimonial.SetStatus("approved")
	if testimonial.Status() != "approved" {
		t.Errorf("Status() = %q, want %q", testimonial.Status(), "approved")
	}

	// Test CreatedAt
	testimonial.SetCreatedAt("2024-01-01T00:00:00Z")
	if testimonial.CreatedAt() != "2024-01-01T00:00:00Z" {
		t.Errorf("CreatedAt() = %q, want %q", testimonial.CreatedAt(), "2024-01-01T00:00:00Z")
	}

	// Test UpdatedAt
	testimonial.SetUpdatedAt("2024-01-02T00:00:00Z")
	if testimonial.UpdatedAt() != "2024-01-02T00:00:00Z" {
		t.Errorf("UpdatedAt() = %q, want %q", testimonial.UpdatedAt(), "2024-01-02T00:00:00Z")
	}
}

func TestNewTestimonialFromEntity(t *testing.T) {
	// Test with nil store and nil entity
	_, err := NewTestimonialFromEntity(nil, nil)
	if err == nil {
		t.Error("NewTestimonialFromEntity() with nil store should return error")
	}
	if err.Error() != "store cannot be nil" {
		t.Errorf("NewTestimonialFromEntity() with nil store returned wrong error: %v", err)
	}
}
