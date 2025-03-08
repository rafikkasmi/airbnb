package reviews

import "regexp"

const ep = "https://www.airbnb.com/api/v3/StaysPdpReviewsQuery/dec1c8061483e78373602047450322fd474e79ba9afa8d3dbbc27f504030f91d"

var regexNumber = regexp.MustCompile(`\d+`)
var treament = []string{
	"feed_map_decouple_m11_treatment",
	"stays_search_rehydration_treatment_desktop",
	"stays_search_rehydration_treatment_moweb",
	"selective_query_feed_map_homepage_desktop_treatment",
	"selective_query_feed_map_homepage_moweb_treatment",
}
