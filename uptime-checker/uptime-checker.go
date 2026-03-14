package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// main checks the availability of URLs passed via command-line arguments.
// It performs DNS lookup and HTTP requests concurrently using goroutines.
func main() {
	// Get URLs from command-line arguments (skip program name)
	urls := os.Args[1:]

	// Create HTTP client with a timeout to avoid hanging requests
	client := &http.Client{Timeout: 5 * time.Second}

	// WaitGroup ensures the program waits for all goroutines to finish
	var wg sync.WaitGroup

	for _, argName := range urls {
		wg.Add(1)

		// Launch a goroutine for each URL so checks run concurrently
		go func(rawUrl string) {
			// Mark goroutine as complete when function exits
			defer wg.Done()

			// Parse the URL so we can extract the hostname for DNS lookup
			parsedUrl, err := url.Parse(rawUrl)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Extract hostname (removes scheme like https://)
			host := parsedUrl.Hostname()

			// Perform DNS lookup to resolve hostname to IP addresses
			ip, err := net.LookupHost(host)
			if err != nil {
				log.Fatalf("Error: %s %v", rawUrl, err)
			}

			// Send HTTP GET request to the URL
			res, err := client.Get(rawUrl)
			if err != nil {
				log.Fatalf("Error: %s %v", rawUrl, err)
			}

			// Ensure response body is closed to prevent resource leaks
			defer res.Body.Close()

			// Consider any status code <= 299 as "UP"
			if res.StatusCode <= 299 {
				fmt.Printf("%-6s status:%-9d url:%-30s ip:%v\n", "[UP]", res.StatusCode, rawUrl, ip)
			} else {
				fmt.Printf("%-6s status:%-9d url:%-30s\n", "[DOWN]", res.StatusCode, rawUrl)
			}

		}(argName) // Pass loop variable explicitly to avoid goroutine capture issues
	}

	// Wait until all goroutines finish execution
	wg.Wait()
}
