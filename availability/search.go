package availability

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"gobnb/api"
	"gobnb/trace"
	"gobnb/utils"
)

// Endpoint is defined in variables.go

func (input InputData) GetAvailabilityCalendar(currency string, proxyURL *url.URL) (AvailabilityData, []CalendarDay, error) {
	apiKey, err := api.Get(proxyURL)
	if err != nil {
		// return AvailabilityData{}, nil, trace.NewOrAdd(1, "availability", "GetAvailabilityCalendar", err, "")
		apiKey = "d306zoyjsyarp7ifhu67rjxn52tv0t20"
	}

	// Fetch availability calendar data
	resultsRaw, err := input.search(currency, apiKey, proxyURL)
	if err != nil {
		errData := trace.NewOrAdd(2, "availability", "GetAvailabilityCalendar", err, "")
		log.Println(errData)
		return AvailabilityData{}, nil, errData
	}
	var daysData []CalendarDay

	// Convert the raw data to our clean format
	availabilityData := resultsRaw.standardize(input.RoomId)

	// Log some basic info about the availability data
	totalMonths := len(availabilityData.CalendarMonths)
	totalDays := 0
	availableDays := 0

	for _, month := range availabilityData.CalendarMonths {
		for _, day := range month.Days {
			totalDays++
			daysData = append(daysData, day)
			if day.Available {
				availableDays++
			}
		}
	}

	fmt.Printf("Availability calendar: %d months, %d total days, %d available days\n",
		totalMonths, totalDays, availableDays)

	return availabilityData, daysData, nil
}

// Implementation of the search function to fetch availability calendar data
func (input InputData) search(currency, apiKey string, proxyURL *url.URL) (root, error) {
	urlParsed, err := url.Parse(Endpoint)
	if err != nil {
		return root{}, trace.NewOrAdd(1, "availability", "search", err, "")
	}
	query := url.Values{}

	extensions := extensions{
		PersistedQuery: persistedQuery{
			Version:    1,
			Sha256Hash: "8f08e03c7bd16fcad3c92a3592c19a8b559a0d0855a84028d1163d4733ed9ade",
		},
	}

	variables := variables{
		Request: request{
			Count:     6, // Number of months to fetch
			ListingId: fmt.Sprintf("%d", input.RoomId),
			Month:     input.StartMonth,
			Year:      input.StartYear,
		},
	}
	extensionsJson, _ := json.Marshal(extensions)
	variablesJson, _ := json.Marshal(variables)

	query.Add("extensions", string(extensionsJson))
	query.Add("variables", string(variablesJson))
	query.Add("locale", "en")
	query.Add("currency", currency)
	urlParsed.RawQuery = query.Encode()
	urlToUse := urlParsed.String()

	// Retry configuration
	maxRetries := 5
	var lastErr error
	var data root

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Add a delay before each retry (except the first attempt)
		if attempt > 0 {
			// Exponential backoff with jitter
			// Base delay: 2^attempt seconds + random jitter (0-1000ms)
			backoffTime := time.Duration(1<<uint(attempt)) * time.Second
			jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
			delay := backoffTime + jitter

			log.Printf("Rate limit encountered. Retrying in %v (attempt %d/%d)...", delay, attempt+1, maxRetries)
			time.Sleep(delay)
		}

		req, err := http.NewRequest("GET", urlToUse, nil)
		if err != nil {
			lastErr = trace.NewOrAdd(3, "availability", "search", err, "")
			continue
		}
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Add("Accept-Language", "en")
		req.Header.Add("Cache-Control", "no-cache")
		req.Header.Add("Pragma", "no-cache")
		req.Header.Add("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
		req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Add("X-Airbnb-Api-Key", apiKey)
		req.Header.Add("Sec-Ch-Ua-Platform", `"Windows"`)
		req.Header.Add("Sec-Fetch-Dest", "document")
		req.Header.Add("Sec-Fetch-Mode", "navigate")
		req.Header.Add("Sec-Fetch-Site", "none")
		req.Header.Add("Sec-Fetch-User", "?1")
		req.Header.Add("Upgrade-Insecure-Requests", "1")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		// Vary the user agent slightly to appear more human-like
		if attempt > 0 {
			userAgents := []string{
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0",
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15",
			}
			req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
		}

		transport := &http.Transport{
			MaxIdleConnsPerHost: 30,
			DisableKeepAlives:   true,
		}
		if proxyURL != nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
		client := &http.Client{
			Timeout: time.Minute,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: transport,
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = trace.NewOrAdd(4, "availability", "search", err, "")
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Always close the body

		if err != nil {
			lastErr = trace.NewOrAdd(5, "availability", "search", err, "")
			continue
		}

		// Check for rate limiting (429)
		if resp.StatusCode == 429 {
			errData := fmt.Sprintf("status: %d headers: %+v", resp.StatusCode, resp.Header)
			lastErr = trace.NewOrAdd(6, "availability", "search", trace.ErrStatusCode, errData)
			// Continue to retry
			continue
		}

		// Check for other errors
		if resp.StatusCode != 200 {
			errData := fmt.Sprintf("status: %d headers: %+v", resp.StatusCode, resp.Header)
			lastErr = trace.NewOrAdd(6, "availability", "search", trace.ErrStatusCode, errData)
			// If it's not a rate limit error, we might want to break here depending on the status code
			// For now, we'll retry all non-200 responses
			continue
		}

		body = utils.RemoveSpaceByte(body) //some values are returned with weird empty values

		if err := json.Unmarshal(body, &data); err != nil {
			lastErr = trace.NewOrAdd(7, "availability", "search", err, "")
			continue
		}

		// If we got here, we succeeded
		return data, nil
	}

	// If we exhausted all retries, return the last error
	return root{}, lastErr
}

// CursorData represents the structure to be encoded as base64
type CursorData struct {
	SectionOffset int `json:"section_offset"`
	ItemsOffset   int `json:"items_offset"`
	Version       int `json:"version"`
}

// GenerateCursor creates a base64 encoded cursor string based on the provided multiplier
// The items_offset will be calculated as x * 40
func GenerateCursor(x int) (string, error) {
	// Create the cursor data structure
	data := CursorData{
		SectionOffset: 0,
		ItemsOffset:   x * 40,
		Version:       1,
	}

	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshaling cursor data: %w", err)
	}

	// Encode the JSON to base64
	encodedCursor := base64.StdEncoding.EncodeToString(jsonData)

	return encodedCursor, nil
}
