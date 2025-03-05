package reviews

import (
	"encoding/json"
	"fmt"
	"time"

	"gobnb/trace"
)

// ParseReviewResponse parses the raw API response into a structured ReviewData object
func ParseReviewResponse(body []byte) (ReviewData, error) {
	// Print the raw response for debugging
	fmt.Printf("Raw API response: %s\n", string(body))
	
	var rawResponse RawReviewResponse
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		fmt.Printf("Failed to parse response: %s\n", string(body))
		return ReviewData{}, trace.NewOrAdd(1, "reviews", "ParseReviewResponse", err, "")
	}
	
	// Print the parsed structure
	fmt.Printf("Parsed response structure: %+v\n", rawResponse)

	// Extract the reviews data from the new structure
	pdpReviews := rawResponse.Data.PdpReviews.Reviews
	
	// Create the standardized review data
	reviewData := ReviewData{
		TotalReviews:   pdpReviews.Metadata.ReviewsCount,
		Rating:         pdpReviews.Metadata.RatingValue,
		HasMoreReviews: pdpReviews.PaginationInfo.HasNextPage,
		Reviews:        make([]Review, 0, len(pdpReviews.Reviews)),
	}

	// Convert each raw review to our standardized format
	for _, rawReview := range pdpReviews.Reviews {
		review := Review{
			ID:             rawReview.ID,
			AuthorID:       rawReview.Reviewer.ID,
			AuthorName:     rawReview.Reviewer.FirstName,
			AuthorImageURL: rawReview.Reviewer.PictureURL,
			Comments:       rawReview.Comments,
			Rating:         rawReview.Rating,
		}

		// Parse the created date
		createdAt, err := parseDate(rawReview.CreatedAt)
		if err != nil {
			// Log the error but continue processing
			fmt.Printf("Error parsing date %s: %v\n", rawReview.CreatedAt, err)
		} else {
			review.CreatedAt = createdAt
		}

		// Handle collection tag if present
		if rawReview.CollectionTag.Tag != "" {
			review.CollectionTag = rawReview.CollectionTag.Tag
		}

		// Handle response if present
		if rawReview.Response.ID != "" {
			review.ResponseID = rawReview.Response.ID
			review.ResponseComment = rawReview.Response.Comments
			
			// Parse response date if available
			if rawReview.Response.CreatedAt != "" {
				responseDate, err := parseDate(rawReview.Response.CreatedAt)
				if err != nil {
					fmt.Printf("Error parsing response date %s: %v\n", rawReview.Response.CreatedAt, err)
				} else {
					review.ResponseDate = responseDate
				}
			}
		}

		reviewData.Reviews = append(reviewData.Reviews, review)
	}

	return reviewData, nil
}

// parseDate converts a date string from the API to a time.Time object
func parseDate(dateStr string) (time.Time, error) {
	// The API might return dates in different formats
	// Try a few common formats
	formats := []string{
		"2006-01-02T15:04:05Z",      // ISO 8601
		"2006-01-02T15:04:05-07:00", // ISO 8601 with timezone
		"2006-01-02",                // Simple date
	}

	var parseErr error
	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return t, nil
		}
		parseErr = err
	}

	return time.Time{}, fmt.Errorf("could not parse date '%s': %w", dateStr, parseErr)
}
