package details

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gobnb/trace"
)

func getFromRoomURL(roomURL string, proxyURL *url.URL) (Data, PriceDependencyInput, []*http.Cookie, error) {
	// Add a small random delay before making the request to avoid rate limiting
	// This helps distribute requests over time
	minDelay := 500  // 0.5 second minimum delay
	maxDelay := 2000 // 2 seconds maximum delay
	randomDelay := minDelay + rand.Intn(maxDelay-minDelay)
	time.Sleep(time.Duration(randomDelay) * time.Millisecond)
	
	// List of user agents to rotate through
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0",
	}
	
	// Retry configuration
	maxRetries := 5
	var lastErr error
	var resp *http.Response
	var body []byte
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Add a delay before each retry (except the first attempt)
		if attempt > 0 {
			// Exponential backoff with jitter
			// Base delay: 2^attempt seconds + random jitter (0-1000ms)
			backoffTime := time.Duration(1<<uint(attempt)) * time.Second
			jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
			delay := backoffTime + jitter
			
			fmt.Printf("Rate limit encountered for details. Retrying in %v (attempt %d/%d)...\n", delay, attempt+1, maxRetries)
			time.Sleep(delay)
		}
		
		// Create a new request for each attempt
		req, err := http.NewRequest("GET", roomURL, nil)
		if err != nil {
			lastErr = trace.NewOrAdd(1, "main", "getFromRoomURL", err, "")
			continue
		}
		
		// Set headers
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Add("Accept-Language", "en")
		req.Header.Add("Cache-Control", "no-cache")
		req.Header.Add("Pragma", "no-cache")
		req.Header.Add("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
		req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Add("Sec-Ch-Ua-Platform", `"Windows"`)
		req.Header.Add("Sec-Fetch-Dest", "document")
		req.Header.Add("Sec-Fetch-Mode", "navigate")
		req.Header.Add("Sec-Fetch-Site", "none")
		req.Header.Add("Sec-Fetch-User", "?1")
		req.Header.Add("Upgrade-Insecure-Requests", "1")
		
		// Rotate user agents to appear more human-like
		userAgentIndex := rand.Intn(len(userAgents))
		req.Header.Add("User-Agent", userAgents[userAgentIndex])
		
		// Configure transport
		transport := &http.Transport{
			MaxIdleConnsPerHost: 30,
			DisableKeepAlives:   true,
		}
		if proxyURL != nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
		
		// Configure client
		client := &http.Client{
			Timeout: time.Minute,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: transport,
		}
		
		// Execute request
		resp, err = client.Do(req)
		if err != nil {
			lastErr = trace.NewOrAdd(2, "main", "getFromRoomURL", err, "")
			continue
		}
		
		// Read response body
		body, err = io.ReadAll(resp.Body)
		resp.Body.Close() // Always close the body
		
		if err != nil {
			lastErr = trace.NewOrAdd(3, "main", "getFromRoomURL", err, "")
			continue
		}
		
		// Check for rate limiting (429)
		if resp.StatusCode == 429 {
			errData := fmt.Sprintf("status: %d headers: %+v", resp.StatusCode, resp.Header)
			lastErr = trace.NewOrAdd(4, "main", "getFromRoomURL", trace.ErrStatusCode, errData)
			// Continue to retry
			continue
		}
		
		// Check for other errors
		if resp.StatusCode != 200 {
			errData := fmt.Sprintf("status: %d headers: %+v", resp.StatusCode, resp.Header)
			lastErr = trace.NewOrAdd(4, "main", "getFromRoomURL", trace.ErrStatusCode, errData)
			// For non-429 errors, we might want to break here depending on the status code
			// For now, we'll retry all non-200 responses
			continue
		}
		
		// If we got here, we succeeded
		break
	}
	
	// If we exhausted all retries, return the last error
	if resp == nil || resp.StatusCode != 200 {
		return Data{}, PriceDependencyInput{}, nil, lastErr
	}
	data, priceDependencyInput, err := ParseBodyDetails(body)

	dataUrl := roomURL
	//remove https://www.airbnb.com/rooms/ from url
	roomID := strings.ReplaceAll(dataUrl, "https://www.airbnb.com/rooms/", "")
	roomIDParsed, _ := strconv.ParseInt(roomID, 10, 64)
	data.RoomID = roomIDParsed
	if err != nil {
		return Data{}, PriceDependencyInput{}, nil, trace.NewOrAdd(5, "main", "getFromRoomURL", err, "")
	}
	data.URL = roomURL
	return data, priceDependencyInput, resp.Cookies(), nil
}
