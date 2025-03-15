package middleware

import (
	"net/http"
	"time"

	"project/app/config"
)

// Middleware represents a middleware handler
type Middleware func(http.Handler) http.Handler

// Chain chains multiple middleware together
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Logger logs information about each request
func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a custom response writer to capture the status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Call the next handler
			next.ServeHTTP(rw, r)

			// Log the request
			duration := time.Since(start)

			// In a real application, you would use a proper logger
			// For simplicity, we're using the standard log package
			// log.Printf("%s %s %d %s", r.Method, r.URL.Path, rw.statusCode, duration)

			// For now, we'll just print to stdout
			println(r.Method, r.URL.Path, rw.statusCode, duration.String())
		})
	}
}

// CORS adds Cross-Origin Resource Sharing headers
func CORS(cfg *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers based on the environment
			if cfg.IsDevelopment() {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				w.Header().Set("Access-Control-Allow-Origin", cfg.AppURL)
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Auth middleware for authentication
func Auth(cfg *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the token from the Authorization header
			token := r.Header.Get("Authorization")

			// In a real application, you would validate the token
			// For simplicity, we're just checking if it exists
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// Recover middleware to recover from panics
func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// In a real application, you would log the error
					// For simplicity, we're just printing to stdout
					println("Recovered from panic:", err)

					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// Custom response writer to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write writes the response and sets the status code to 200 if not already set
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}
