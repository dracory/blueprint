package middlewares

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"project/internal/app"

	"github.com/dracory/req"
	"github.com/dracory/rtr"
)

// MaintenanceState represents the JSON structure of maintenance_mode_state.json
type MaintenanceState struct {
	Message           string   `json:"message"`
	RetryAfterSeconds int      `json:"retry_after_seconds"`
	ExcludeIPs        []string `json:"exclude_ips"`
	ExcludePaths      []string `json:"exclude_paths"`
	CreatedAt         string   `json:"created_at"`
}

// maintenanceMiddleware holds the internal state for the maintenance middleware
type maintenanceMiddleware struct {
	app          app.AppInterface
	filePath     string
	cacheDur     time.Duration
	mu           sync.Mutex
	cachedState  *MaintenanceState
	cachedExists bool
	lastCheck    time.Time
}

// NewMaintenanceMiddleware creates a middleware that checks for a maintenance mode
// file and returns 503 Service Unavailable when maintenance is active.
func NewMaintenanceMiddleware(app app.AppInterface) rtr.MiddlewareInterface {
	filePath := "maintenance_mode_state.json"
	cacheDur := 30 * time.Second

	if app != nil && app.GetConfig() != nil {
		if fp := app.GetConfig().GetAppMaintenanceFilePath(); fp != "" {
			filePath = fp
		}
	}

	m := &maintenanceMiddleware{
		app:      app,
		filePath: filePath,
		cacheDur: cacheDur,
	}

	return rtr.NewMiddleware().
		SetName("Maintenance Mode Middleware").
		SetHandler(func(next http.Handler) http.Handler {
			return m.handler(next)
		})
}

func (m *maintenanceMiddleware) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var state *MaintenanceState
		var exists bool

		// Check env var override first (for containerized deployments)
		if m.app != nil && m.app.GetConfig() != nil && m.app.GetConfig().GetAppMaintenanceEnabled() {
			state, exists = m.getMaintenanceState()
			if !exists {
				state = &MaintenanceState{Message: "We'll be right back."}
			}
		} else {
			state, exists = m.getMaintenanceState()
			if !exists {
				next.ServeHTTP(w, r)
				return
			}
		}

		ip := req.GetIP(r)
		if isIPExcluded(ip, state.ExcludeIPs) {
			next.ServeHTTP(w, r)
			return
		}

		if isPathExcluded(r.URL.Path, state.ExcludePaths) {
			next.ServeHTTP(w, r)
			return
		}

		m.writeMaintenanceResponse(w, state)
	})
}

func (m *maintenanceMiddleware) getMaintenanceState() (*MaintenanceState, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	if now.Sub(m.lastCheck) < m.cacheDur {
		return m.cachedState, m.cachedExists
	}

	info, err := os.Stat(m.filePath)
	if err != nil {
		m.cachedState = nil
		m.cachedExists = false
		m.lastCheck = now
		return nil, false
	}

	if m.cachedState != nil && m.cachedExists && !info.ModTime().After(m.lastCheck) {
		return m.cachedState, m.cachedExists
	}

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		m.cachedState = nil
		m.cachedExists = false
		m.lastCheck = now
		return nil, false
	}

	var state MaintenanceState
	if err := json.Unmarshal(data, &state); err != nil {
		m.cachedState = nil
		m.cachedExists = false
		m.lastCheck = now
		return nil, false
	}

	m.cachedState = &state
	m.cachedExists = true
	m.lastCheck = now

	return &state, true
}

func (m *maintenanceMiddleware) writeMaintenanceResponse(w http.ResponseWriter, state *MaintenanceState) {
	if state.RetryAfterSeconds > 0 {
		w.Header().Set("Retry-After", fmt.Sprintf("%d", state.RetryAfterSeconds))
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.WriteHeader(http.StatusServiceUnavailable)

	message := state.Message
	if message == "" {
		message = "We'll be right back."
	}

	htmlContent := buildMaintenanceHTML(message)
	_, _ = w.Write([]byte(htmlContent))
}

func buildMaintenanceHTML(message string) string {
	message = html.EscapeString(message)
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Maintenance</title>
    <style>
        body { font-family: system-ui, sans-serif; text-align: center; padding: 50px; }
        h1 { font-size: 2em; color: #333; }
        p { color: #666; }
    </style>
</head>
<body>
    <h1>Undergoing Maintenance</h1>
    <p>%s</p>
</body>
</html>`, message)
}

// isIPExcluded checks if the given IP matches any entry in excludeIPs
func isIPExcluded(ip string, excludeIPs []string) bool {
	for _, excluded := range excludeIPs {
		if strings.TrimSpace(excluded) == ip {
			return true
		}
	}
	return false
}

// isPathExcluded checks if the given path matches any exclude pattern.
// Supports wildcard matching: "/admin/*" matches "/admin/anything"
func isPathExcluded(path string, excludePaths []string) bool {
	for _, pattern := range excludePaths {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		if strings.HasSuffix(pattern, "/*") {
			prefix := strings.TrimSuffix(pattern, "/*")
			if strings.HasPrefix(path, prefix+"/") || path == prefix {
				return true
			}
		} else if pattern == path {
			return true
		} else if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}
	}
	return false
}
