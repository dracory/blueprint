package testutils

import (
	"log"
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock"
)

// SetupMailServer creates a mock SMTP server for testing email functionality.
// Returns the server instance and a cleanup function to stop the server.
func SetupMailServer(t *testing.T) (*smtpmock.Server, func()) {
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       false, // Set to true for debugging
		LogServerActivity: true,
		// Use port 0 so the OS assigns a free ephemeral port, avoiding collisions
		PortNumber:  0,
		HostAddress: "127.0.0.1",
	})

	if err := server.Start(); err != nil {
		t.Fatalf("Failed to start mock SMTP server: %v", err)
	}

	return server, func() {
		if err := server.Stop(); err != nil {
			log.Printf("Warning: failed to stop mock SMTP server: %v", err)
		}
	}
}
