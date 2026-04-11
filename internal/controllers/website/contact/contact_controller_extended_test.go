package contact

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestContactController_AnyIndex(t *testing.T) {
	registry := testutils.Setup()
	controller := NewContactController(registry)

	// Create a test request
	req := httptest.NewRequest("GET", "/contact", nil)
	w := httptest.NewRecorder()

	// Test the handler
	html := controller.AnyIndex(w, req)

	// Verify the response
	if html == "" {
		t.Error("AnyIndex() should return non-empty HTML")
	}
	if !strings.Contains(html, "Contact") {
		t.Error("AnyIndex() should contain 'Contact' in the HTML")
	}
	if !strings.Contains(html, "container") {
		t.Error("AnyIndex() should contain container class")
	}
}

func TestFormContact_GetKind(t *testing.T) {
	registry := testutils.Setup()
	component := NewFormContact(registry)

	if component == nil {
		t.Fatal("NewFormContact() should not return nil")
	}

	// Test GetKind via type assertion
	fc, ok := component.(*formContact)
	if !ok {
		t.Fatalf("Component is not *formContact type, got %T", component)
	}
	kind := fc.GetKind()
	if kind != "website_contact_form" {
		t.Errorf("GetKind() = %q, want %q", kind, "website_contact_form")
	}
}

func TestFormContact_Mount(t *testing.T) {
	registry := testutils.Setup()
	component := NewFormContact(registry)

	if component == nil {
		t.Fatal("NewFormContact() should not return nil")
	}

	// Test Mount with empty params
	fc, ok := component.(*formContact)
	if !ok {
		t.Fatalf("Component is not *formContact type, got %T", component)
	}
	err := fc.Mount(context.TODO(), map[string]string{})
	// Should not error with empty params
	if err != nil {
		t.Errorf("Mount() with empty params returned error: %v", err)
	}

	// Verify CSRF token was generated
	if fc.CsrfToken == "" {
		t.Error("Mount() should generate CSRF token")
	}

	// Verify captcha was initialized
	if fc.CaptchaQuestion == "" {
		t.Error("Mount() should generate captcha question")
	}
	if fc.CaptchaExpected == "" {
		t.Error("Mount() should generate captcha expected value")
	}
}

func TestFormContact_MountWithUserID(t *testing.T) {
	registry := testutils.Setup()
	component := NewFormContact(registry)

	if component == nil {
		t.Fatal("NewFormContact() should not return nil")
	}

	// Test Mount with user_id param (but no actual user in store)
	fc, ok := component.(*formContact)
	if !ok {
		t.Fatalf("Component is not *formContact type, got %T", component)
	}
	err := fc.Mount(context.TODO(), map[string]string{"user_id": "test-user-id"})
	// Should not error even with non-existent user
	if err != nil {
		t.Errorf("Mount() with user_id returned error: %v", err)
	}
}

func TestFormContact_Fields(t *testing.T) {
	registry := testutils.Setup()
	component := NewFormContact(registry)

	if component == nil {
		t.Fatal("NewFormContact() should not return nil")
	}

	// Test that component was properly initialized
	fc, ok := component.(*formContact)
	if !ok {
		t.Fatalf("Component is not *formContact type, got %T", component)
	}
	// Test setting fields
	fc.Email = "test@example.com"
	fc.FirstName = "John"
	fc.LastName = "Doe"
	fc.Text = "Test message"

	if fc.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", fc.Email, "test@example.com")
	}
	if fc.FirstName != "John" {
		t.Errorf("FirstName = %q, want %q", fc.FirstName, "John")
	}
	if fc.LastName != "Doe" {
		t.Errorf("LastName = %q, want %q", fc.LastName, "Doe")
	}
	if fc.Text != "Test message" {
		t.Errorf("Text = %q, want %q", fc.Text, "Test message")
	}
}

func TestFormContact_CanUpdateFields(t *testing.T) {
	registry := testutils.Setup()
	component := NewFormContact(registry)

	if component == nil {
		t.Fatal("NewFormContact() should not return nil")
	}

	// Test that CanUpdate fields are properly set after Mount
	fc, ok := component.(*formContact)
	if !ok {
		t.Fatalf("Component is not *formContact type, got %T", component)
	}
	// Initially these should be false/zero values
	if fc.CanUpdateEmail {
		t.Error("CanUpdateEmail should be false initially")
	}
	if fc.CanUpdateFirst {
		t.Error("CanUpdateFirst should be false initially")
	}
	if fc.CanUpdateLast {
		t.Error("CanUpdateLast should be false initially")
	}

	// After mount, they should be true (for new form)
	err := fc.Mount(context.TODO(), map[string]string{})
	if err != nil {
		t.Errorf("Mount() returned error: %v", err)
	}

	if !fc.CanUpdateEmail {
		t.Error("CanUpdateEmail should be true after Mount")
	}
	if !fc.CanUpdateFirst {
		t.Error("CanUpdateFirst should be true after Mount")
	}
	if !fc.CanUpdateLast {
		t.Error("CanUpdateLast should be true after Mount")
	}
}

func TestFormContact_CaptchaFields(t *testing.T) {
	registry := testutils.Setup()
	component := NewFormContact(registry)

	if component == nil {
		t.Fatal("NewFormContact() should not return nil")
	}

	// Test captcha fields after Mount
	fc, ok := component.(*formContact)
	if !ok {
		t.Fatalf("Component is not *formContact type, got %T", component)
	}
	err := fc.Mount(context.TODO(), map[string]string{})
	if err != nil {
		t.Errorf("Mount() returned error: %v", err)
	}

	// Verify captcha question format
	if fc.CaptchaQuestion == "" {
		t.Error("CaptchaQuestion should not be empty")
	}
	// Should contain "+" for addition
	if !strings.Contains(fc.CaptchaQuestion, "+") {
		t.Error("CaptchaQuestion should contain '+' for addition")
	}
	// Should end with "="
	if !strings.HasSuffix(fc.CaptchaQuestion, "=") {
		t.Error("CaptchaQuestion should end with '='")
	}

	// Verify expected hash is set
	if fc.CaptchaExpected == "" {
		t.Error("CaptchaExpected should not be empty")
	}
}
