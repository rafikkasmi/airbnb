package reviews

import (
	"time"
)

// ReviewData represents the overall review data for a room
type ReviewData struct {
	RoomID         int64    `json:"room_id"`
	TotalReviews   int      `json:"total_reviews"`
	Rating         float64  `json:"rating"`
	Reviews        []Review `json:"reviews"`
	HasMoreReviews bool     `json:"has_more_reviews"`
}

// Review represents a single review
type Review struct {
	ID              string    `json:"id"`
	AuthorID        string    `json:"author_id"`
	AuthorName      string    `json:"author_name"`
	AuthorImageURL  string    `json:"author_image_url"`
	Comments        string    `json:"comments"`
	CreatedAt       time.Time `json:"created_at"`
	CollectionTag   string    `json:"collection_tag,omitempty"`
	Rating          int       `json:"rating"`
	ResponseID      string    `json:"response_id,omitempty"`
	ResponseComment string    `json:"response_comment,omitempty"`
	ResponseDate    time.Time `json:"response_date,omitempty"`
}

// RawReviewResponse represents the raw API response structure for the new endpoint
type RawReviewResponse struct {
	Data struct {
		PdpReviews struct {
			Reviews struct {
				Metadata struct {
					ReviewsCount int     `json:"reviewsCount"`
					RatingValue  float64 `json:"ratingValue"`
				} `json:"metadata"`
				Reviews []struct {
					ID         string `json:"id"`
					Comments   string `json:"comments"`
					CreatedAt  string `json:"createdAt"`
					LocalizedDate string `json:"localizedDate"`
					Rating     int    `json:"rating"`
					Reviewer   struct {
						ID        string `json:"id"`
						FirstName string `json:"firstName"`
						PictureURL string `json:"pictureUrl"`
					} `json:"reviewer"`
					Response struct {
						ID         string `json:"id,omitempty"`
						Comments   string `json:"comments,omitempty"`
						CreatedAt  string `json:"createdAt,omitempty"`
					} `json:"response,omitempty"`
					CollectionTag struct {
						Tag string `json:"tag,omitempty"`
					} `json:"collectionTag,omitempty"`
				} `json:"reviews"`
				PaginationInfo struct {
					HasNextPage bool `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"paginationInfo"`
			} `json:"reviews"`
		} `json:"pdpReviews"`
	} `json:"data"`
}
