package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"project/internal/config"
)

func main() {
	cfg, err := config.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	defaultURL := cfg.GetAppUrl()
	url := flag.String("url", defaultURL, "Target URL to load test")
	concurrency := flag.Int("c", 10, "Number of concurrent requests")
	duration := flag.Duration("d", 30*time.Second, "Duration of the load test")
	timeout := flag.Duration("t", 10*time.Second, "Request timeout")
	rateLimit := flag.Int("r", 0, "Rate limit (requests per second, 0 = unlimited)")
	flag.Parse()

	fmt.Printf("Load Testing: %s\n", *url)
	fmt.Printf("Concurrency: %d\n", *concurrency)
	fmt.Printf("Duration: %v\n", *duration)
	fmt.Printf("Timeout: %v\n", *timeout)
	if *rateLimit > 0 {
		fmt.Printf("Rate Limit: %d req/sec\n", *rateLimit)
	}
	fmt.Printf("Connection Pooling: Enabled (MaxIdle=%d)\n\n", *concurrency)

	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: *concurrency,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	}

	client := &http.Client{
		Timeout:   *timeout,
		Transport: transport,
	}

	var (
		totalRequests   int64
		successRequests int64
		failedRequests  int64
		totalDuration   int64
		minDuration     int64 = 1<<63 - 1
		maxDuration     int64
		mu              sync.Mutex
		errorCounts     = make(map[string]int64)
	)

	ctx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()

	var wg sync.WaitGroup
	startTime := time.Now()

	var rateLimiter <-chan time.Time
	if *rateLimit > 0 {
		ticker := time.NewTicker(time.Second / time.Duration(*rateLimit))
		defer ticker.Stop()
		rateLimiter = ticker.C
	}

	progressTicker := time.NewTicker(2 * time.Second)
	defer progressTicker.Stop()
	go func() {
		for {
			select {
			case <-progressTicker.C:
				fmt.Print(".")
			case <-ctx.Done():
				return
			}
		}
	}()

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					if rateLimiter != nil {
						select {
						case <-rateLimiter:
						case <-ctx.Done():
							return
						}
					}
					reqStart := time.Now()
					resp, err := client.Get(*url)
					reqDuration := time.Since(reqStart).Nanoseconds()

					atomic.AddInt64(&totalRequests, 1)

					if err != nil {
						atomic.AddInt64(&failedRequests, 1)
						mu.Lock()
						errorCounts[err.Error()]++
						mu.Unlock()
					} else {
						// Validate HTTP status code (2xx = success)
						if resp.StatusCode >= 200 && resp.StatusCode < 300 {
							atomic.AddInt64(&successRequests, 1)
						} else {
							atomic.AddInt64(&failedRequests, 1)
							mu.Lock()
							errorCounts[fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)]++
							mu.Unlock()
						}
						resp.Body.Close()
					}

					atomic.AddInt64(&totalDuration, reqDuration)

					mu.Lock()
					if reqDuration < minDuration {
						minDuration = reqDuration
					}
					if reqDuration > maxDuration {
						maxDuration = reqDuration
					}
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	elapsedTime := time.Since(startTime)

	total := atomic.LoadInt64(&totalRequests)
	success := atomic.LoadInt64(&successRequests)
	failed := atomic.LoadInt64(&failedRequests)

	var avgDuration time.Duration
	if total > 0 {
		avgDuration = time.Duration(atomic.LoadInt64(&totalDuration) / total)
	} else {
		avgDuration = 0
	}

	fmt.Println("\n\n=== Load Test Results ===")
	fmt.Printf("Total Requests: %d\n", total)
	if total > 0 {
		fmt.Printf("Successful: %d (%.2f%%)\n", success, float64(success)/float64(total)*100)
		fmt.Printf("Failed: %d (%.2f%%)\n", failed, float64(failed)/float64(total)*100)
	} else {
		fmt.Printf("Successful: %d (0.00%%)\n", success)
		fmt.Printf("Failed: %d (0.00%%)\n", failed)
	}
	if total > 0 {
		fmt.Printf("Requests/sec: %.2f\n", float64(total)/elapsedTime.Seconds())
	} else {
		fmt.Printf("Requests/sec: 0\n")
	}
	fmt.Printf("Avg Response Time: %v\n", avgDuration)
	fmt.Printf("Min Response Time: %v\n", time.Duration(minDuration))
	fmt.Printf("Max Response Time: %v\n", time.Duration(maxDuration))
	fmt.Printf("Total Test Duration: %v\n", elapsedTime)

	if len(errorCounts) > 0 {
		fmt.Println("\n=== Error Breakdown ===")
		for errMsg, count := range errorCounts {
			fmt.Printf("%d: %s\n", count, errMsg)
		}
	}
}
