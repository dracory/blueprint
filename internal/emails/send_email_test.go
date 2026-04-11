package emails

import (
	"testing"

	"project/internal/testutils"
)

func TestInitEmailSender(t *testing.T) {
	// Save and restore global state to avoid test pollution
	senderMu.Lock()
	originalSender := emailSender
	emailSender = nil
	senderMu.Unlock()
	defer func() {
		senderMu.Lock()
		emailSender = originalSender
		senderMu.Unlock()
	}()

	// Test with nil registry
	InitEmailSender(nil)

	senderMu.RLock()
	if emailSender != nil {
		senderMu.RUnlock()
		t.Error("emailSender should remain nil after InitEmailSender(nil)")
	} else {
		senderMu.RUnlock()
	}

	// Test with valid registry
	registry := testutils.Setup()
	InitEmailSender(registry)

	senderMu.RLock()
	if emailSender == nil {
		senderMu.RUnlock()
		t.Error("emailSender should be initialized after InitEmailSender with valid registry")
	} else {
		senderMu.RUnlock()
	}
}

func TestSendEmail(t *testing.T) {
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
	err := SendEmail(SendOptions{To: []string{"test@example.com"}})
	if err == nil {
		t.Error("SendEmail should return error when sender is not initialized")
	}
	if err.Error() != "email sender is not initialized" {
		t.Errorf("SendEmail error = %q, want %q", err.Error(), "email sender is not initialized")
	}
}

func TestSendOptions(t *testing.T) {
	// Test SendOptions struct initialization
	options := SendOptions{
		From:     "from@example.com",
		FromName: "From Name",
		To:       []string{"to@example.com"},
		Bcc:      []string{"bcc@example.com"},
		Cc:       []string{"cc@example.com"},
		Subject:  "Test Subject",
		HtmlBody: "<p>HTML Body</p>",
		TextBody: "Text Body",
	}

	if options.From != "from@example.com" {
		t.Errorf("From = %q, want %q", options.From, "from@example.com")
	}
	if options.FromName != "From Name" {
		t.Errorf("FromName = %q, want %q", options.FromName, "From Name")
	}
	if len(options.To) != 1 || options.To[0] != "to@example.com" {
		t.Errorf("To = %v, want [to@example.com]", options.To)
	}
	if len(options.Bcc) != 1 || options.Bcc[0] != "bcc@example.com" {
		t.Errorf("Bcc = %v, want [bcc@example.com]", options.Bcc)
	}
	if len(options.Cc) != 1 || options.Cc[0] != "cc@example.com" {
		t.Errorf("Cc = %v, want [cc@example.com]", options.Cc)
	}
	if options.Subject != "Test Subject" {
		t.Errorf("Subject = %q, want %q", options.Subject, "Test Subject")
	}
	if options.HtmlBody != "<p>HTML Body</p>" {
		t.Errorf("HtmlBody = %q, want %q", options.HtmlBody, "<p>HTML Body</p>")
	}
	if options.TextBody != "Text Body" {
		t.Errorf("TextBody = %q, want %q", options.TextBody, "Text Body")
	}
}

func TestMailerStruct(t *testing.T) {
	// Test Mailer struct initialization
	mailer := Mailer{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user",
		Password: "pass",
	}

	if mailer.Host != "smtp.example.com" {
		t.Errorf("Host = %q, want %q", mailer.Host, "smtp.example.com")
	}
	if mailer.Port != 587 {
		t.Errorf("Port = %d, want 587", mailer.Port)
	}
	if mailer.Username != "user" {
		t.Errorf("Username = %q, want %q", mailer.Username, "user")
	}
	if mailer.Password != "pass" {
		t.Errorf("Password = %q, want %q", mailer.Password, "pass")
	}
}
