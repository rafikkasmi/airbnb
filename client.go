package gobnb

import (
	"net/url"

	"gobnb/details"
	"gobnb/reviews"
	"gobnb/search"
)

type Client struct {
	Currency string //ISO currency, example: USD, EUR
	ProxyURL *url.URL
}

func DefaultClient() Client {
	client := Client{
		Currency: "USD",
		ProxyURL: nil,
	}
	return client
}
func NewClient(currency string, proxyURL *url.URL) Client {
	client := Client{
		Currency: currency,
		ProxyURL: proxyURL,
	}
	return client
}

func (cl Client) DetailsFromRoomURL(roomURL string) (details.Data, error) {
	return details.GetFromRoomURL(roomURL, cl.Currency, cl.ProxyURL)
}
func (cl Client) DetailsFromRoomID(roomID int64) (details.Data, error) {
	return details.GetFromRoomID(roomID, cl.Currency, cl.ProxyURL)
}

func (cl Client) DetailsFromRoomIDAndDomain(roomID int64, domain string) (details.Data, error) {
	return details.GetFromRoomIDAndDomain(roomID, domain, cl.Currency, cl.ProxyURL)
}

func (cl Client) DetailsMainRoomIds(mailURL string) ([]int64, error) {
	return details.GetMainRoomIds(mailURL, cl.ProxyURL)
}

// ReviewsFromRoomID fetches reviews for a specific room ID
func (cl Client) ReviewsFromRoomID(roomID int64) (reviews.ReviewData, error) {
	return reviews.GetFromRoomID(roomID, cl.ProxyURL)
}

// ReviewsFromRoomURL fetches reviews for a room using its URL
func (cl Client) ReviewsFromRoomURL(roomURL string) (reviews.ReviewData, error) {
	return reviews.GetFromRoomURL(roomURL, cl.ProxyURL)
}

// AllReviewsFromRoomID fetches all reviews for a room by making multiple paginated requests
func (cl Client) AllReviewsFromRoomID(roomID int64) (reviews.ReviewData, error) {
	return reviews.GetAllReviewsPaginated(roomID, cl.ProxyURL)
}

// coordinates should 2 points one from south and one from north(if you think it like a square, this points represent the two points of the diagonal from this square)
// zoom value from 1 - 20, so from the "square" like I said on the coorinates, this represents how much zoom on this squere there is
func (cl Client) SearchAll(zoomValue int, coordinates search.CoordinatesInput, check search.Check) ([]search.Data, error) {
	input := search.InputData{
		ZoomValue:   zoomValue,
		Coordinates: coordinates,
		Check:       check,
	}
	return input.SearchAll(cl.Currency, cl.ProxyURL)
}

// coordinates should 2 points one from south and one from north(if you think it like a square, this points represent the two points of the diagonal from this square)
// zoom value from 1 - 20, so from the "square" like I said on the coorinates, this represents how much zoom on this squere there is
func (cl Client) SearchFirstPage(zoomValue int, coordinates search.CoordinatesInput, check search.Check) ([]search.Data, error) {
	input := search.InputData{
		ZoomValue:   zoomValue,
		Coordinates: coordinates,
		Check:       check,
	}
	return input.SearchFirstPage(cl.Currency, cl.ProxyURL)
}
