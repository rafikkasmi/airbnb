package reviews

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"gobnb/trace"
)

// getFromRoomID fetches review data for a specific room ID
func getFromRoomID(roomID int64, proxyURL *url.URL) (ReviewData, []*http.Cookie, error) {
	// First, fetch the room page to get cookies
	roomPageURL := fmt.Sprintf("https://www.airbnb.com/rooms/%d", roomID)
	
	// Create a request to the room page first to get cookies and establish a session
	pageReq, err := http.NewRequest("GET", roomPageURL, nil)
	if err != nil {
		return ReviewData{}, nil, trace.NewOrAdd(1, "reviews", "getFromRoomID", err, "")
	}
	
	// Set headers for the room page request (similar to details package)
	pageReq.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	pageReq.Header.Add("Accept-Language", "en-US,en;q=0.9")
	pageReq.Header.Add("Cache-Control", "no-cache")
	pageReq.Header.Add("Pragma", "no-cache")
	pageReq.Header.Add("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	pageReq.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	pageReq.Header.Add("Sec-Ch-Ua-Platform", `"Windows"`)
	pageReq.Header.Add("Sec-Fetch-Dest", "document")
	pageReq.Header.Add("Sec-Fetch-Mode", "navigate")
	pageReq.Header.Add("Sec-Fetch-Site", "none")
	pageReq.Header.Add("Sec-Fetch-User", "?1")
	pageReq.Header.Add("Upgrade-Insecure-Requests", "1")
	pageReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	
	transport := &http.Transport{
		MaxIdleConnsPerHost: 30,
		DisableKeepAlives:   true,
	}
	if proxyURL != nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	}
	
	client := &http.Client{
		Timeout: time.Minute,
		Transport: transport,
	}
	
	// Execute the page request to get cookies
	pageResp, err := client.Do(pageReq)
	if err != nil {
		return ReviewData{}, nil, trace.NewOrAdd(2, "reviews", "getFromRoomID", err, "")
	}
	// Don't close the body yet, we need to read it to extract the API key
	
	if pageResp.StatusCode != 200 {
		errData := fmt.Sprintf("room page status: %d headers: %+v", pageResp.StatusCode, pageResp.Header)
		return ReviewData{}, nil, trace.NewOrAdd(3, "reviews", "getFromRoomID", trace.ErrStatusCode, errData)
	}
	
	// Read the page body to extract API key
	pageBody, err := io.ReadAll(pageResp.Body)
	if err != nil {
		return ReviewData{}, nil, trace.NewOrAdd(3, "reviews", "getFromRoomID", err, "")
	}
	pageResp.Body.Close()
	
	// Extract the API key from the page content
	apiKeyRegex := regexp.MustCompile(`"key":"(.+?)"`) // Similar to the one in details package
	apiKeyMatch := apiKeyRegex.FindStringSubmatch(string(pageBody))
	if len(apiKeyMatch) < 2 {
		return ReviewData{}, nil, trace.NewOrAdd(4, "reviews", "getFromRoomID", fmt.Errorf("could not extract API key"), "")
	}
	apiKey := apiKeyMatch[1]
	fmt.Printf("Found API key: %s\n", apiKey)
	
	// Now construct the URL for the reviews API endpoint using the exact format from the user
	// Format the room ID as a base64 encoded string with prefix
	listingID := fmt.Sprintf("U3RheUxpc3Rpbmc6%d", roomID)
	
	// Construct the URL using the exact format provided
	reviewURL := fmt.Sprintf(
		"https://www.airbnb.com/api/v3/StaysPdpReviewsQuery/dec1c8061483e78373602047450322fd474e79ba9afa8d3dbbc27f504030f91d?operationName=StaysPdpReviewsQuery&locale=en&currency=USD&variables=%%7B%%22id%%22%%3A%%22%s%%22%%2C%%22pdpReviewsRequest%%22%%3A%%7B%%22fieldSelector%%22%%3A%%22for_p3_translation_only%%22%%2C%%22forPreview%%22%%3Afalse%%2C%%22limit%%22%%3A24%%2C%%22offset%%22%%3A%%220%%22%%2C%%22showingTranslationButton%%22%%3Afalse%%2C%%22first%%22%%3A24%%2C%%22sortingPreference%%22%%3A%%22RATING_DESC%%22%%2C%%22numberOfAdults%%22%%3A%%221%%22%%2C%%22numberOfChildren%%22%%3A%%220%%22%%2C%%22numberOfInfants%%22%%3A%%220%%22%%2C%%22numberOfPets%%22%%3A%%220%%22%%7D%%7D&extensions=%%7B%%22persistedQuery%%22%%3A%%7B%%22version%%22%%3A1%%2C%%22sha256Hash%%22%%3A%%22dec1c8061483e78373602047450322fd474e79ba9afa8d3dbbc27f504030f91d%%22%%7D%%7D",
		listingID,
	)
	
	// Debug information
	fmt.Printf("Requesting reviews with URL: %s\n", reviewURL)

	// Create a new request for the reviews API
	req, err := http.NewRequest("GET", reviewURL, nil)
	if err != nil {
		return ReviewData{}, nil, trace.NewOrAdd(4, "reviews", "getFromRoomID", err, "")
	}

	// Set headers for the API request
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Referer", roomPageURL)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Add("X-Airbnb-Api-Key", apiKey)
	
	// Add cookies from the page request
	for _, cookie := range pageResp.Cookies() {
		req.AddCookie(cookie)
	}
	
	// Reuse the same transport for the API request
	// We already defined the transport and client above, so we can reuse them
	resp, err := client.Do(req)
	if err != nil {
		return ReviewData{}, nil, trace.NewOrAdd(5, "reviews", "getFromRoomID", err, "")
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ReviewData{}, nil, trace.NewOrAdd(6, "reviews", "getFromRoomID", err, "")
	}
	
	if resp.StatusCode != 200 {
		// Print the response body for debugging
		fmt.Printf("Response body: %s\n", string(body))
		errData := fmt.Sprintf("status: %d headers: %+v", resp.StatusCode, resp.Header)
		return ReviewData{}, nil, trace.NewOrAdd(7, "reviews", "getFromRoomID", trace.ErrStatusCode, errData)
	}
	
	// Parse the response body to extract review data
	reviewData, err := ParseReviewResponse(body)
	if err != nil {
		return ReviewData{}, nil, trace.NewOrAdd(8, "reviews", "getFromRoomID", err, "")
	}
	
	// Set the room ID in the review data
	reviewData.RoomID = roomID
	
	return reviewData, resp.Cookies(), nil
}

// getFromRoomURL extracts the room ID from a URL and fetches the reviews
func getFromRoomURL(roomURL string, proxyURL *url.URL) (ReviewData, []*http.Cookie, error) {
	// Extract room ID from URL
	roomID, err := extractRoomIDFromURL(roomURL)
	if err != nil {
		return ReviewData{}, nil, trace.NewOrAdd(1, "reviews", "getFromRoomURL", err, "")
	}
	
	return getFromRoomID(roomID, proxyURL)
}

// Helper function to extract room ID from URL
func extractRoomIDFromURL(roomURL string) (int64, error) {
	// This is a simplified version - you may need more robust parsing
	var roomID int64
	_, err := fmt.Sscanf(roomURL, "https://www.airbnb.com/rooms/%d", &roomID)
	if err != nil {
		return 0, fmt.Errorf("failed to extract room ID from URL: %w", err)
	}
	return roomID, nil
}
