package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	RequestID int
	StatusCode int
	Duration   time.Duration
	RateLimitRemaining string
	Error      error
}

func main() {
	// Define command-line flags
	url := flag.String("url", "http://localhost:8080/service-a", "The URL to send requests to.")
	apiKey := flag.String("api-key", "super-secret-key", "The API key to include in the X-API-KEY header.")
	totalRequests := flag.Int("n", 50, "Total number of requests to send.")
	concurrency := flag.Int("c", 10, "Number of concurrent workers.")
	flag.Parse()

	if *totalRequests <= 0 || *concurrency <= 0 {
		log.Fatal("Number of requests and concurrency must be greater than 0.")
	}

	log.Printf("Starting load test with %d requests and %d concurrency on %s\n", *totalRequests, *concurrency, *url)

	jobs := make(chan int, *totalRequests)
	results := make(chan Result, *totalRequests)
	
	var wg sync.WaitGroup
	
	// Start workers
	for w := 1; w <= *concurrency; w++ {
		wg.Add(1)
		go worker(w, *url, *apiKey, jobs, results, &wg)
	}

	// Send jobs
	for j := 1; j <= *totalRequests; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process and print results
	processResults(results, *totalRequests)
}

func worker(id int, url string, apiKey string, jobs <-chan int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for jobID := range jobs {
		start := time.Now()

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			results <- Result{RequestID: jobID, Error: err}
			continue
		}

		if apiKey != "" {
			req.Header.Set("X-API-KEY", apiKey)
		}

		resp, err := client.Do(req)
		duration := time.Since(start)

		if err != nil {
			results <- Result{RequestID: jobID, Duration: duration, Error: err}
			continue
		}

		results <- Result{
			RequestID:          jobID,
			StatusCode:         resp.StatusCode,
			Duration:           duration,
			RateLimitRemaining: resp.Header.Get("X-RateLimit-Remaining"),
		}
		resp.Body.Close()
	}
}

func processResults(results <-chan Result, totalRequests int) {
	successCount := 0
	rateLimitedCount := 0
	errorCount := 0
	
	fmt.Println("-----------------------------------------------------------------")
	fmt.Printf("%-10s %-10s %-15s %-20s\n", "Req ID", "Status", "Duration", "X-RateLimit-Remaining")
	fmt.Println("-----------------------------------------------------------------")

	startTime := time.Now()

	for result := range results {
		if result.Error != nil {
			errorCount++
			log.Printf("Request %d failed: %v", result.RequestID, result.Error)
			continue
		}

		if result.StatusCode == http.StatusOK {
			successCount++
		} else if result.StatusCode == http.StatusTooManyRequests {
			rateLimitedCount++
		} else {
			errorCount++
		}
		
		fmt.Printf("%-10d %-10d %-15s %-20s\n", result.RequestID, result.StatusCode, result.Duration.Round(time.Millisecond), result.RateLimitRemaining)
	}

	totalDuration := time.Since(startTime)
	rps := float64(totalRequests) / totalDuration.Seconds()

	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("Summary:")
	fmt.Printf("Total time:        %s\n", totalDuration.Round(time.Second))
	fmt.Printf("Total requests:    %d\n", totalRequests)
	fmt.Printf("Requests per sec:  %.2f\n", rps)
	fmt.Printf("Successful (200):  %d\n", successCount)
	fmt.Printf("Rate limited (429):%d\n", rateLimitedCount)
	fmt.Printf("Other errors:      %d\n", errorCount)
	fmt.Println("-----------------------------------------------------------------")
}

