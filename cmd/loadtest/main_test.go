package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoadTestWithMockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}))
	defer server.Close()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var totalRequests int64
	var successRequests int64

	for i := 0; i < 100; i++ {
		resp, err := client.Get(server.URL)
		totalRequests++

		if err != nil {
			t.Logf("Request failed: %v", err)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			successRequests++
		}
		resp.Body.Close()
	}

	if totalRequests != 100 {
		t.Errorf("expected 100 requests, got %d", totalRequests)
	}

	if successRequests != 100 {
		t.Errorf("expected 100 successful requests, got %d", successRequests)
	}
}

func TestLoadTestWithSlowServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}))
	defer server.Close()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	startTime := time.Now()
	for i := 0; i < 10; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		resp.Body.Close()
	}
	duration := time.Since(startTime)

	expectedMinDuration := 1000 * time.Millisecond
	if duration < expectedMinDuration {
		t.Errorf("expected duration >= %v, got %v", expectedMinDuration, duration)
	}
}

func TestLoadTestWithFailingServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error")
	}))
	defer server.Close()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var totalRequests int64
	var failedRequests int64

	for i := 0; i < 50; i++ {
		resp, err := client.Get(server.URL)
		totalRequests++

		if err != nil {
			failedRequests++
			continue
		}

		if resp.StatusCode >= 400 {
			failedRequests++
		}
		resp.Body.Close()
	}

	if totalRequests != 50 {
		t.Errorf("expected 50 requests, got %d", totalRequests)
	}

	if failedRequests != 50 {
		t.Errorf("expected 50 failed requests, got %d", failedRequests)
	}
}

func TestLoadTestWithTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	var failedRequests int64

	for i := 0; i < 10; i++ {
		_, err := client.Get(server.URL)
		if err != nil {
			failedRequests++
		}
	}

	if failedRequests == 0 {
		t.Error("expected some requests to timeout, but all succeeded")
	}
}

func TestResponseTimeTracking(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}))
	defer server.Close()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var minDuration int64 = 1<<63 - 1
	var maxDuration int64

	for i := 0; i < 20; i++ {
		reqStart := time.Now()
		resp, err := client.Get(server.URL)
		reqDuration := time.Since(reqStart).Nanoseconds()

		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		resp.Body.Close()

		if reqDuration < minDuration {
			minDuration = reqDuration
		}
		if reqDuration > maxDuration {
			maxDuration = reqDuration
		}
	}

	if minDuration > maxDuration {
		t.Error("min duration should not be greater than max duration")
	}

	if minDuration <= 0 {
		t.Error("min duration should be greater than 0")
	}

	if maxDuration <= 0 {
		t.Error("max duration should be greater than 0")
	}
}

func TestConcurrentRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}))
	defer server.Close()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	done := make(chan bool, 5)
	var successCount int64

	for i := 0; i < 5; i++ {
		go func() {
			resp, err := client.Get(server.URL)
			if err == nil && resp.StatusCode == http.StatusOK {
				successCount++
				resp.Body.Close()
			}
			done <- true
		}()
	}

	for i := 0; i < 5; i++ {
		<-done
	}

	if successCount != 5 {
		t.Errorf("expected 5 successful concurrent requests, got %d", successCount)
	}
}
