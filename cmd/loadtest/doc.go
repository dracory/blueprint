// Package main provides a command-line load testing tool for testing application performance.
//
// The loadtest tool simulates concurrent HTTP requests to measure application throughput,
// response times, and reliability under load. It's useful for identifying performance
// bottlenecks and capacity limits before production deployment.
//
// # Usage
//
// Run the load test with default settings (URL from APP_URL environment variable):
//
//	go run ./cmd/loadtest
//
// Or with custom parameters:
//
//	go run ./cmd/loadtest -url=http://localhost:8080 -c=100 -d=1m -t=10s
//
// # Flags
//
//   - url: Target URL to load test (default: from APP_URL environment variable)
//   - c: Number of concurrent requests (default: 10)
//   - d: Duration of the load test (default: 30s)
//   - t: Request timeout (default: 10s)
//
// # Output
//
// The tool displays:
//   - Total requests made
//   - Success and failure rates
//   - Requests per second (throughput)
//   - Average, minimum, and maximum response times
//   - Total test duration
//
// # Configuration
//
// The tool automatically loads the default URL from your application's environment
// configuration by calling config.NewFromEnv(). This reads from:
//   - .env file if present
//   - APP_URL environment variable
//   - Falls back to http://localhost:8080 if not configured
//
// # Examples
//
// Test homepage with 50 concurrent requests for 1 minute:
//
//	go run ./cmd/loadtest -url=http://localhost:8080 -c=50 -d=1m
//
// Test API endpoint with 200 concurrent requests for 2 minutes:
//
//	go run ./cmd/loadtest -url=http://localhost:8080/api/posts -c=200 -d=2m
//
// Test with custom timeout:
//
//	go run ./cmd/loadtest -url=http://localhost:8080/blog -c=100 -d=30s -t=15s
package main
