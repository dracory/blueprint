package emails

import (
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestEmailToAdminOnNewContactFormSubmitted(t *testing.T) {
	// Test with nil app
	email := NewEmailToAdminOnNewContactFormSubmitted(nil)
	if email == nil {
		t.Fatal("NewEmailToAdminOnNewContactFormSubmitted(nil) should return non-nil")
	}
	if email.app != nil {
		t.Error("email.app should be nil when passed nil")
	}

	// Test with valid app
	app := testutils.Setup()
	email = NewEmailToAdminOnNewContactFormSubmitted(app)
	if email == nil {
		t.Fatal("NewEmailToAdminOnNewContactFormSubmitted(app) should return non-nil")
	}
	if email.app == nil {
		t.Error("email.app should not be nil when passed valid app")
	}
}

func TestEmailToAdminOnNewContactFormSubmitted_Template(t *testing.T) {
	app := testutils.Setup()
	email := NewEmailToAdminOnNewContactFormSubmitted(app)

	// Test template generation
	html := email.template("TestApp")

	if html == "" {
		t.Error("template() should return non-empty HTML")
	}
	if !strings.Contains(html, "New Contact Form Submitted") {
		t.Error("template() should contain the heading")
	}
	if !strings.Contains(html, "TestApp") {
		t.Error("template() should contain the app name")
	}
}

func TestEmailToAdminOnNewContactFormSubmitted_Send(t *testing.T) {
	// Save and restore global state
	senderMu.Lock()
	originalSender := emailSender
	emailSender = nil
	senderMu.Unlock()
	defer func() {
		senderMu.Lock()
		emailSender = originalSender
		senderMu.Unlock()
	}()

	// Test with uninitialized sender
	email := NewEmailToAdminOnNewContactFormSubmitted(nil)
	err := email.Send()
	if err == nil {
		t.Error("Send() with uninitialized sender should return error")
	}

	// Test with valid app but uninitialized sender
	app := testutils.Setup()
	email = NewEmailToAdminOnNewContactFormSubmitted(app)
	err = email.Send()
	if err == nil {
		t.Error("Send() with uninitialized sender should return error")
	}
}

func TestEmailToAdminOnNewUserRegistered(t *testing.T) {
	// Test with nil app
	email := NewEmailToAdminOnNewUserRegistered(nil)
	if email == nil {
		t.Fatal("NewEmailToAdminOnNewUserRegistered(nil) should return non-nil")
	}
	if email.app != nil {
		t.Error("email.app should be nil when passed nil")
	}

	// Test with valid app
	app := testutils.Setup()
	email = NewEmailToAdminOnNewUserRegistered(app)
	if email == nil {
		t.Fatal("NewEmailToAdminOnNewUserRegistered(app) should return non-nil")
	}
	if email.app == nil {
		t.Error("email.app should not be nil when passed valid app")
	}
}

func TestEmailToAdminOnNewUserRegistered_Template(t *testing.T) {
	app := testutils.Setup()
	email := NewEmailToAdminOnNewUserRegistered(app)

	// Test template generation
	html := email.template("TestApp", "user-123")

	if html == "" {
		t.Error("template() should return non-empty HTML")
	}
	if !strings.Contains(html, "user-123") {
		t.Error("template() should contain the user ID")
	}
	if !strings.Contains(html, "TestApp") {
		t.Error("template() should contain the app name")
	}
}

func TestEmailToAdminOnNewUserRegistered_Send(t *testing.T) {
	// Save and restore global state
	senderMu.Lock()
	originalSender := emailSender
	emailSender = nil
	senderMu.Unlock()
	defer func() {
		senderMu.Lock()
		emailSender = originalSender
		senderMu.Unlock()
	}()

	// Test with uninitialized sender
	email := NewEmailToAdminOnNewUserRegistered(nil)
	err := email.Send("user-123")
	if err == nil {
		t.Error("Send() with uninitialized sender should return error")
	}

	// Test with valid app but uninitialized sender
	app := testutils.Setup()
	email = NewEmailToAdminOnNewUserRegistered(app)
	err = email.Send("user-123")
	if err == nil {
		t.Error("Send() with uninitialized sender should return error")
	}
}

func TestEmailNotifyAdmin(t *testing.T) {
	// Test with nil app
	email := NewEmailNotifyAdmin(nil)
	if email == nil {
		t.Fatal("NewEmailNotifyAdmin(nil) should return non-nil")
	}
	if email.app != nil {
		t.Error("email.app should be nil when passed nil")
	}

	// Test with valid app
	app := testutils.Setup()
	email = NewEmailNotifyAdmin(app)
	if email == nil {
		t.Fatal("NewEmailNotifyAdmin(app) should return non-nil")
	}
	if email.app == nil {
		t.Error("email.app should not be nil when passed valid app")
	}
}

func TestEmailNotifyAdmin_Send(t *testing.T) {
	// Save and restore global state
	senderMu.Lock()
	originalSender := emailSender
	emailSender = nil
	senderMu.Unlock()
	defer func() {
		senderMu.Lock()
		emailSender = originalSender
		senderMu.Unlock()
	}()

	// Test with uninitialized sender
	email := NewEmailNotifyAdmin(nil)
	err := email.Send("<p>Test HTML</p>")
	if err == nil {
		t.Error("Send() with uninitialized sender should return error")
	}

	// Test with valid app but uninitialized sender
	app := testutils.Setup()
	email = NewEmailNotifyAdmin(app)
	err = email.Send("<p>Test HTML</p>")
	if err == nil {
		t.Error("Send() with uninitialized sender should return error")
	}
}

func TestInviteFriendEmail(t *testing.T) {
	// Test with nil app and nil user store
	email := NewInviteFriendEmail(nil, nil)
	if email == nil {
		t.Fatal("NewInviteFriendEmail(nil, nil) should return non-nil")
	}
	if email.app != nil {
		t.Error("email.app should be nil when passed nil")
	}
	if email.userStore != nil {
		t.Error("email.userStore should be nil when passed nil")
	}

	// Test with valid app
	app := testutils.Setup()
	email = NewInviteFriendEmail(app, nil)
	if email == nil {
		t.Fatal("NewInviteFriendEmail(app, nil) should return non-nil")
	}
	if email.app == nil {
		t.Error("email.app should not be nil when passed valid app")
	}
}

func TestInviteFriendEmail_Template(t *testing.T) {
	app := testutils.Setup()
	email := NewInviteFriendEmail(app, nil)

	// Test template generation
	html := email.template("TestApp", "John", "Hello friend!", "Jane")

	if html == "" {
		t.Error("template() should return non-empty HTML")
	}
	if !strings.Contains(html, "John") {
		t.Error("template() should contain the user name")
	}
	if !strings.Contains(html, "Jane") {
		t.Error("template() should contain the recipient name")
	}
	if !strings.Contains(html, "Hello friend!") {
		t.Error("template() should contain the user note")
	}
	if !strings.Contains(html, "TestApp") {
		t.Error("template() should contain the app name")
	}
}

func TestInviteFriendEmail_Send(t *testing.T) {
	// Save and restore global state
	senderMu.Lock()
	originalSender := emailSender
	emailSender = nil
	senderMu.Unlock()
	defer func() {
		senderMu.Lock()
		emailSender = originalSender
		senderMu.Unlock()
	}()

	// Test with nil user store (sender is also nil)
	email := NewInviteFriendEmail(nil, nil)
	err := email.Send("user-123", "Hello!", "friend@example.com", "Friend")
	if err == nil {
		t.Error("Send() with nil userStore should return error")
	}

	// Test with valid app but nil user store
	app := testutils.Setup()
	email = NewInviteFriendEmail(app, nil)
	err = email.Send("user-123", "Hello!", "friend@example.com", "Friend")
	if err == nil {
		t.Error("Send() with nil userStore should return error")
	}
	if err.Error() != "user store not configured" {
		t.Errorf("Send() error = %q, want %q", err.Error(), "user store not configured")
	}
}
