package reviews

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"gobnb/api"
	"gobnb/trace"
	"gobnb/utils"
)

// Endpoint is defined in variables.go

func (input InputData) GetAllReviewsOfRoom(roomId int64, currency string, proxyURL *url.URL) ([]Review, error) {
	apiKey, err := api.Get(proxyURL)
	if err != nil {
		return nil, trace.NewOrAdd(1, "search", "GetAllReviewsOfRoom", err, "")
	}
	var allResults []Review

	for i := 0; i < 8; i++ {
		resultsRaw, err := input.search(i, currency, apiKey, proxyURL)
		if err != nil {
			errData := trace.NewOrAdd(2, "search", "SearchAll", err, "")
			log.Println(errData)
			break
		}

		// Create a Data object with the room ID and reviews
		//append all resultRaw.Reviews to all results
		allResults = append(allResults, standardizeReviews(resultsRaw.Reviews)...)

		fmt.Printf("Reviews count: %d\n", len(allResults))

		// If we got fewer reviews than requested, we've reached the end
		if len(resultsRaw.Reviews) == 0 {
			break
		}
	}
	fmt.Println("Total reviews count: ", len(allResults))
	return allResults, nil
}

// standardizeReviews converts the raw review data into our Review struct format
func standardizeReviews(rawReviews []review) []Review {
	var reviews []Review

	for _, r := range rawReviews {
		translatedComments := r.Comments
		if r.LocalizedReview.Comments != "" {
			translatedComments = r.LocalizedReview.Comments
		}

		review := Review{
			ID:                 r.Id,
			Comments:           r.Comments,
			TranslatedComments: translatedComments,
			Language:           r.Language,
			CreatedAt:          r.CreatedAt,
			LocalizedDate:      r.LocalizedDate,
			Rating:             r.Rating,
			Highlight:          r.ReviewHighlight,
			Reviewer: User{
				ID:            r.Reviewer.Id,
				FirstName:     r.Reviewer.FirstName,
				FullName:      r.Reviewer.HostName,
				IsSuperhost:   r.Reviewer.IsSuperhost,
				ProfilePicURL: r.Reviewer.PictureUrl,
				ProfilePath:   r.Reviewer.ProfilePath,
			},
			Host: User{
				ID:            r.Reviewee.Id,
				FirstName:     r.Reviewee.FirstName,
				FullName:      r.Reviewee.HostName,
				IsSuperhost:   r.Reviewee.IsSuperhost,
				ProfilePicURL: r.Reviewee.PictureUrl,
				ProfilePath:   r.Reviewee.ProfilePath,
			},
		}

		reviews = append(reviews, review)
	}

	return reviews
}

func (input InputData) search(page int, currency, apiKey string, proxyURL *url.URL) (rootdatapresentationstayssearchResults, error) {
	// checkinS := getStringDate(input.Check.In)
	// checkoutS := getStringDate(input.Check.Out)
	// hours := input.Check.Out.Sub(input.Check.In).Hours()
	// days := int(hours / 24)
	urlParsed, err := url.Parse(ep)
	if err != nil {
		return rootdatapresentationstayssearchResults{}, trace.NewOrAdd(1, "search", "search", err, "")
	}
	query := url.Values{}

	extensions := extensions{
		PersistedQuery: persistedQuery{
			Version:    1,
			Sha256Hash: "dec1c8061483e78373602047450322fd474e79ba9afa8d3dbbc27f504030f91d",
		},
	}

	strRoomId := fmt.Sprintf("StayListing:%d", input.RoomId)
	roomData := utils.ToBase64(strRoomId)
	variables := variables{
		Id: roomData,
		PdpReviewsRequest: pdpReviewsRequest{
			FieldSelector:     "id",
			ForPreview:        false,
			Limit:             50,
			Offset:            fmt.Sprintf("%d", page*50),
			SortingPreference: "MOST_RECENT",
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

	req, err := http.NewRequest("GET", urlToUse, nil)
	if err != nil {
		return rootdatapresentationstayssearchResults{}, trace.NewOrAdd(3, "search", "search", err, "")
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
		return rootdatapresentationstayssearchResults{}, trace.NewOrAdd(4, "search", "search", err, "")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rootdatapresentationstayssearchResults{}, trace.NewOrAdd(5, "search", "search", err, "")
	}
	if resp.StatusCode != 200 {
		errData := fmt.Sprintf("status: %d headers: %+v", resp.StatusCode, resp.Header)
		return rootdatapresentationstayssearchResults{}, trace.NewOrAdd(6, "search", "search", trace.ErrStatusCode, errData)
	}
	body = utils.RemoveSpaceByte(body) //some values are returned with weird empty values
	var data root
	if err := json.Unmarshal(body, &data); err != nil {
		return rootdatapresentationstayssearchResults{}, trace.NewOrAdd(7, "search", "search", err, "")
	}
	return data.Data.Presentation.StayProductDetailPage.Reviews, nil
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
