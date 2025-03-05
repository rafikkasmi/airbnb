package reviews

import (
	"net/url"

	"gobnb/trace"
)

// GetFromRoomID fetches reviews for a specific room ID
func GetFromRoomID(roomID int64, proxyURL *url.URL) (ReviewData, error) {
	reviewData, _, err := getFromRoomID(roomID, proxyURL)
	if err != nil {
		return ReviewData{}, trace.NewOrAdd(1, "reviews", "GetFromRoomID", err, "")
	}
	
	return reviewData, nil
}

// GetFromRoomURL fetches reviews for a room using its URL
func GetFromRoomURL(roomURL string, proxyURL *url.URL) (ReviewData, error) {
	reviewData, _, err := getFromRoomURL(roomURL, proxyURL)
	if err != nil {
		return ReviewData{}, trace.NewOrAdd(1, "reviews", "GetFromRoomURL", err, "")
	}
	
	return reviewData, nil
}

// GetAllReviewsFromRoomID fetches all reviews for a room by making multiple paginated requests
func GetAllReviewsFromRoomID(roomID int64, proxyURL *url.URL) ([]Review, error) {
	var allReviews []Review
	var cursor string
	limit := 24 // Number of reviews per request, matching Airbnb's default
	
	for {
		// Fetch a batch of reviews
		params := PaginationParams{
			Offset: 0,
			Count:  limit,
			Cursor: cursor,
		}
		
		reviewData, err := GetReviewsPage(roomID, params, proxyURL)
		if err != nil {
			return nil, trace.NewOrAdd(1, "reviews", "GetAllReviewsFromRoomID", err, "")
		}
		
		// Add reviews to the collection
		allReviews = append(allReviews, reviewData.Reviews...)
		
		// Check if we've received all reviews
		if len(reviewData.Reviews) < limit || !reviewData.HasMoreReviews {
			break
		}
		
		// For now, we'll just use offset-based pagination since we don't have direct access to the cursor
		// This is a simpler approach that should work for most cases
		params.Offset += len(reviewData.Reviews)
		cursor = "" // Clear cursor to use offset-based pagination
	}
	
	return allReviews, nil
}

// GetReviewsWithPagination fetches a specific page of reviews
func GetReviewsWithPagination(roomID int64, offset, count int, proxyURL *url.URL) (ReviewData, error) {
	// Implementation would be similar to getFromRoomID but with pagination parameters
	// This is a placeholder - you would need to implement the actual API call
	
	// For now, we'll just call the regular getFromRoomID
	reviewData, _, err := getFromRoomID(roomID, proxyURL)
	if err != nil {
		return ReviewData{}, trace.NewOrAdd(1, "reviews", "GetReviewsWithPagination", err, "")
	}
	
	return reviewData, nil
}
