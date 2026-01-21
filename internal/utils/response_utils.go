package utils

import (
	"io"
	"log/slog"
)

// SafeCloseResponseBody safely closes an HTTP response body with proper error handling.
// This utility function ensures consistent error handling and logging across the application.
func SafeCloseResponseBody(body io.Closer) {
	if body == nil {
		return
	}

	if err := body.Close(); err != nil {
		// Log the error but don't panic - the body might already be closed
		// or there might be network issues preventing proper cleanup.
		slog.Error("failed to close response body", "error", err)
	}
}
