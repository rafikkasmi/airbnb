package reviews

import (
	// "encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"gobnb/trace"
)

// PaginationParams contains parameters for paginated review requests
type PaginationParams struct {
	Offset int
	Count  int
	Cursor string
}

// GetReviewsPage fetches a specific page of reviews using offset-based pagination
func GetReviewsPage(roomID int64, params PaginationParams, proxyURL *url.URL) (ReviewData, error) {
	// First, fetch the room page to get cookies
	roomPageURL := fmt.Sprintf("https://www.airbnb.com/rooms/%d", roomID)
	
	// Create a request to the room page first to get cookies and establish a session
	pageReq, err := http.NewRequest("GET", roomPageURL, nil)
	if err != nil {
		return ReviewData{}, trace.NewOrAdd(1, "reviews", "GetReviewsPage", err, "")
	}
	
	// Set headers for the room page request
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
		return ReviewData{}, trace.NewOrAdd(2, "reviews", "GetReviewsPage", err, "")
	}
	pageResp.Body.Close() // We don't need the body content, just the cookies
	
	if pageResp.StatusCode != 200 {
		errData := fmt.Sprintf("room page status: %d headers: %+v", pageResp.StatusCode, pageResp.Header)
		return ReviewData{}, trace.NewOrAdd(3, "reviews", "GetReviewsPage", trace.ErrStatusCode, errData)
	}
	
	// Now construct the URL for the reviews API endpoint with pagination using the exact format from the user
	// Format the room ID as a base64 encoded string with prefix
	listingID := fmt.Sprintf("U3RheUxpc3Rpbmc6%d", roomID)
	
	// Construct the URL using the exact format provided
	var reviewURL string
	if params.Cursor != "" {
		// Use cursor-based pagination if a cursor is provided
		reviewURL = fmt.Sprintf(
			"https://www.airbnb.com/api/v3/StaysPdpReviewsQuery/dec1c8061483e78373602047450322fd474e79ba9afa8d3dbbc27f504030f91d?operationName=StaysPdpReviewsQuery&locale=en&currency=USD&variables=%%7B%%22id%%22%%3A%%22%s%%22%%2C%%22pdpReviewsRequest%%22%%3A%%7B%%22fieldSelector%%22%%3A%%22for_p3_translation_only%%22%%2C%%22forPreview%%22%%3Afalse%%2C%%22limit%%22%%3A%d%%2C%%22cursor%%22%%3A%%22%s%%22%%2C%%22showingTranslationButton%%22%%3Afalse%%2C%%22first%%22%%3A%d%%2C%%22sortingPreference%%22%%3A%%22RATING_DESC%%22%%2C%%22numberOfAdults%%22%%3A%%221%%22%%2C%%22numberOfChildren%%22%%3A%%220%%22%%2C%%22numberOfInfants%%22%%3A%%220%%22%%2C%%22numberOfPets%%22%%3A%%220%%22%%7D%%7D&extensions=%%7B%%22persistedQuery%%22%%3A%%7B%%22version%%22%%3A1%%2C%%22sha256Hash%%22%%3A%%22dec1c8061483e78373602047450322fd474e79ba9afa8d3dbbc27f504030f91d%%22%%7D%%7D",
			listingID, params.Count, params.Cursor, params.Count,
		)
	} else {
		// Use offset-based pagination otherwise
		reviewURL = fmt.Sprintf(
			"https://www.airbnb.com/api/v3/StaysPdpReviewsQuery/dec1c8061483e78373602047450322fd474e79ba9afa8d3dbbc27f504030f91d?operationName=StaysPdpReviewsQuery&locale=en&currency=USD&variables=%%7B%%22id%%22%%3A%%22%s%%22%%2C%%22pdpReviewsRequest%%22%%3A%%7B%%22fieldSelector%%22%%3A%%22for_p3_translation_only%%22%%2C%%22forPreview%%22%%3Afalse%%2C%%22limit%%22%%3A%d%%2C%%22offset%%22%%3A%%22%d%%22%%2C%%22showingTranslationButton%%22%%3Afalse%%2C%%22first%%22%%3A%d%%2C%%22sortingPreference%%22%%3A%%22RATING_DESC%%22%%2C%%22numberOfAdults%%22%%3A%%221%%22%%2C%%22numberOfChildren%%22%%3A%%220%%22%%2C%%22numberOfInfants%%22%%3A%%220%%22%%2C%%22numberOfPets%%22%%3A%%220%%22%%7D%%7D&extensions=%%7B%%22persistedQuery%%22%%3A%%7B%%22version%%22%%3A1%%2C%%22sha256Hash%%22%%3A%%22dec1c8061483e78373602047450322fd474e79ba9afa8d3dbbc27f504030f91d%%22%%7D%%7D",
			listingID, params.Count, params.Offset, params.Count,
		)
	}
	
	// Debug information
	fmt.Printf("Requesting paginated reviews with URL: %s\n", reviewURL)

	// Create a new request for the reviews API
	req, err := http.NewRequest("GET", reviewURL, nil)
	if err != nil {
		return ReviewData{}, trace.NewOrAdd(4, "reviews", "GetReviewsPage", err, "")
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
	
	// Add cookies from the page request
	for _, cookie := range pageResp.Cookies() {
		req.AddCookie(cookie)
	}
	
	// Execute the API request
	resp, err := client.Do(req)
	if err != nil {
		return ReviewData{}, trace.NewOrAdd(5, "reviews", "GetReviewsPage", err, "")
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ReviewData{}, trace.NewOrAdd(6, "reviews", "GetReviewsPage", err, "")
	}
	
	if resp.StatusCode != 200 {
		// Print the response body for debugging
		fmt.Printf("Response body: %s\n", string(body))
		errData := fmt.Sprintf("status: %d headers: %+v", resp.StatusCode, resp.Header)
		return ReviewData{}, trace.NewOrAdd(7, "reviews", "GetReviewsPage", trace.ErrStatusCode, errData)
	}
	
	// Parse the response body to extract review data
	reviewData, err := ParseReviewResponse(body)
	if err != nil {
		return ReviewData{}, trace.NewOrAdd(8, "reviews", "GetReviewsPage", err, "")
	}
	
	// Set the room ID in the review data
	reviewData.RoomID = roomID
	
	// Debug output
	fmt.Printf("Found %d reviews, has more: %v\n", len(reviewData.Reviews), reviewData.HasMoreReviews)
	
	return reviewData, nil
}

// GetAllReviewsPaginated fetches all reviews for a room using pagination
func GetAllReviewsPaginated(roomID int64, proxyURL *url.URL) (ReviewData, error) {
	// Start with an empty result
	result := ReviewData{
		RoomID:   roomID,
		Reviews:  []Review{},
	}
	
	offset := 0
	pageSize := 50 // Number of reviews per request
	
	for {
		// Fetch a page of reviews
		params := PaginationParams{
			Offset: offset,
			Count:  pageSize,
		}
		
		pageData, err := GetReviewsPage(roomID, params, proxyURL)
		if err != nil {
			return result, trace.NewOrAdd(1, "reviews", "GetAllReviewsPaginated", err, "")
		}
		
		// Update metadata from the first page
		if offset == 0 {
			result.TotalReviews = pageData.TotalReviews
			result.Rating = pageData.Rating
		}
		
		// Add reviews to our collection
		result.Reviews = append(result.Reviews, pageData.Reviews...)
		
		// If there are no more reviews or we got fewer than requested, we're done
		if !pageData.HasMoreReviews || len(pageData.Reviews) < pageSize {
			break
		}
		
		// Update offset for next page
		offset += len(pageData.Reviews)
	}
	
	// No more reviews to fetch
	result.HasMoreReviews = false
	
	return result, nil
}
