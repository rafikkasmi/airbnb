package availability

import "regexp"

// Endpoint for the availability calendar API
const Endpoint = "https://www.airbnb.com/api/v3/PdpAvailabilityCalendar"

var regexNumber = regexp.MustCompile(`\d+`)

// InputData represents the input parameters for the availability calendar search
type InputData struct {
	RoomId     int64
	StartMonth int
	StartYear  int
}

// Request structure for the API request
// type request struct {
// 	Count     int    `json:"count"`
// 	ListingId string `json:"listingId"`
// 	Month     int    `json:"month"`
// 	Year      int    `json:"year"`
// }

// // Variables structure for the API request
// type variables struct {
// 	Request request `json:"request"`
// }

// // PersistedQuery structure for the API request
// type persistedQuery struct {
// 	Version    int    `json:"version"`
// 	Sha256Hash string `json:"sha256Hash"`
// }

// // Extensions structure for the API request
// type extensions struct {
// 	PersistedQuery persistedQuery `json:"persistedQuery"`
// }
