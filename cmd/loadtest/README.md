# Load Test Tool

A simple command-line load testing tool for testing the performance of your application endpoints.

## Usage

```bash
go run ./cmd/loadtest [flags]
```

## Flags

- `-url string` - Target URL to load test (default: loaded from `APP_URL` environment variable)
- `-c int` - Number of concurrent requests (default: 10)
- `-d duration` - Duration of the load test (default: 30s)
- `-t duration` - Request timeout (default: 10s)
- `-r int` - Rate limit in requests per second (default: 0 = unlimited)

## Examples

### Small Website Load (Realistic)
Simulates typical small website traffic: ~10 requests/second
**≈ 120-500 concurrent active users** or **1,000-3,000 daily visitors**
```bash
task loadtest:small
# or: go run ./cmd/loadtest -c=5 -r=10 -d=1m
```

### Medium Website Load (Realistic)
Simulates typical medium website traffic: ~50 requests/second
**≈ 600-2,500 concurrent active users** or **5,000-15,000 daily visitors**
```bash
task loadtest:medium
# or: go run ./cmd/loadtest -c=10 -r=50 -d=1m
```

### Large Website Load (Realistic)
Simulates high-traffic website: ~200 requests/second
**≈ 2,400-10,000 concurrent active users** or **20,000-60,000 daily visitors**
```bash
task loadtest:large
# or: go run ./cmd/loadtest -c=20 -r=200 -d=1m
```

### Stress Test (Find Breaking Point)
No rate limit - finds maximum capacity
```bash
task loadtest:stress
# or: go run ./cmd/loadtest -c=50 -d=1m
```

### Spike Test (Sudden Traffic Burst)
Simulates sudden traffic spike (e.g., viral post)
```bash
task loadtest:spike
# or: go run ./cmd/loadtest -c=100 -d=30s
```

### Test a Specific Endpoint
```bash
go run ./cmd/loadtest -url=http://localhost:8080/blog -c=10 -r=20 -d=1m
```

## Output

The tool provides the following metrics:

- **Total Requests** - Total number of requests made during the test
- **Successful** - Number of successful requests (status 200-299)
- **Failed** - Number of failed requests
- **Requests/sec** - Throughput (requests per second)
- **Avg Response Time** - Average response time across all requests
- **Min Response Time** - Fastest response time
- **Max Response Time** - Slowest response time
- **Total Test Duration** - Actual time taken to complete the test

## Example Output

```
Load Testing: http://localhost:8080
Concurrency: 50
Duration: 1m0s
Timeout: 10s

=== Load Test Results ===
Total Requests: 5000
Successful: 4950 (99.00%)
Failed: 50 (1.00%)
Requests/sec: 83.33
Avg Response Time: 600ms
Min Response Time: 50ms
Max Response Time: 2500ms
Total Test Duration: 1m0.123s
```

## Configuration

The tool automatically loads the default URL from your application's environment configuration:

- Reads from `.env` file if present
- Uses `APP_URL` environment variable
- Falls back to `http://localhost:8080` if not configured

## Understanding Load Test Parameters

### What's Realistic?

**Small Website** (blog, portfolio, small business)
- Traffic: 1-10 requests/second
- Daily visitors: 100-1,000
- Use: `task loadtest:small`

**Medium Website** (popular blog, e-commerce)
- Traffic: 10-100 requests/second  
- Daily visitors: 1,000-10,000
- Use: `task loadtest:medium`

**Large Website** (major platform, news site)
- Traffic: 100-1,000+ requests/second
- Daily visitors: 10,000-1,000,000+
- Use: `task loadtest:large`

### Connection Pooling

The tool uses HTTP keep-alive and connection pooling to reuse TCP connections, preventing port exhaustion on Windows. This mirrors real-world browser behavior.

### Rate Limiting vs Unlimited

- **Rate Limited** (`-r` flag): Simulates realistic user traffic patterns
- **Unlimited**: Stress test to find breaking point (not realistic traffic)

## Tips

- **Start realistic**: Use `task loadtest:small` for small websites
- **Use rate limiting**: The `-r` flag simulates real traffic patterns
- **Stress test separately**: Use `task loadtest:stress` to find limits
- **Test critical endpoints**: Homepage, API endpoints, authentication flows
- **Run multiple times**: Get consistent results before making conclusions
- **Monitor resources**: Watch CPU, memory, database connections during tests
- **Avoid production**: Never run unlimited stress tests against production
