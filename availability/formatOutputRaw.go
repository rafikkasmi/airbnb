package availability

type root struct {
	Data rootData `json:"data"`
}

type rootData struct {
	Merlin merlinQuery `json:"merlin"`
}

type merlinQuery struct {
	TypeName              string                            `json:"__typename"`
	PdpAvailabilityCalendar pdpAvailabilityCalendarResponse `json:"pdpAvailabilityCalendar"`
}

type pdpAvailabilityCalendarResponse struct {
	TypeName       string           `json:"__typename"`
	CalendarMonths []calendarMonth `json:"calendarMonths"`
}

type calendarMonth struct {
	TypeName string        `json:"__typename"`
	Month    int           `json:"month"`
	Year     int           `json:"year"`
	Days     []calendarDay `json:"days"`
}

type calendarDay struct {
	TypeName            string             `json:"__typename"`
	CalendarDate        string             `json:"calendarDate"`
	Available           bool               `json:"available"`
	MaxNights           int                `json:"maxNights"`
	MinNights           int                `json:"minNights"`
	AvailableForCheckin bool               `json:"availableForCheckin"`
	AvailableForCheckout bool              `json:"availableForCheckout"`
	Bookable            interface{}        `json:"bookable"`
	Price               calendarDayPrice   `json:"price"`
}

type calendarDayPrice struct {
	TypeName           string      `json:"__typename"`
	LocalPriceFormatted interface{} `json:"localPriceFormatted"`
}

type paginationInfo struct {
	PageCursors        []string `json:"pageCursors"`
	PreviousPageCursor string   `json:"previousPageCursor"`
	NextPageCursor     string   `json:"nextPageCursor"`
}
type review struct {
	TypeName                  string          `json:"__typename"`
	CollectionTag             interface{}     `json:"collectionTag"`
	Comments                  string          `json:"comments"`
	Id                        string          `json:"id"`
	Language                  string          `json:"language"`
	CreatedAt                 string          `json:"createdAt"`
	Reviewee                  reviewUser      `json:"reviewee"`
	Reviewer                  reviewUser      `json:"reviewer"`
	ReviewHighlight           string          `json:"reviewHighlight"`
	HighlightType             string          `json:"highlightType"`
	LocalizedDate             string          `json:"localizedDate"`
	LocalizedRespondedDate    interface{}     `json:"localizedRespondedDate"`
	LocalizedReviewerLocation string          `json:"localizedReviewerLocation"`
	LocalizedReview           localizedReview `json:"localizedReview"`
	Rating                    int             `json:"rating"`
	RatingAccessibilityLabel  string          `json:"ratingAccessibilityLabel"`
	RecommendedNumberOfLines  interface{}     `json:"recommendedNumberOfLines"`
	Response                  interface{}     `json:"response"`
	RoomTypeListingTitle      interface{}     `json:"roomTypeListingTitle"`
	HighlightedReviewSentence []interface{}   `json:"highlightedReviewSentence"`
	HighlightReviewMentioned  interface{}     `json:"highlightReviewMentioned"`
	ShowMoreButton            showMoreButton  `json:"showMoreButton"`
	SubtitleItems             []interface{}   `json:"subtitleItems"`
	Channel                   interface{}     `json:"channel"`
	ReviewMediaItems          []interface{}   `json:"reviewMediaItems"`
	IsHostHighlightedReview   interface{}     `json:"isHostHighlightedReview"`
	ReviewPhotoUrls           []interface{}   `json:"reviewPhotoUrls"`
}

type reviewUser struct {
	TypeName           string             `json:"__typename"`
	Deleted            bool               `json:"deleted"`
	FirstName          string             `json:"firstName"`
	HostName           string             `json:"hostName"`
	Id                 string             `json:"id"`
	PictureUrl         string             `json:"pictureUrl"`
	ProfilePath        string             `json:"profilePath"`
	IsSuperhost        bool               `json:"isSuperhost"`
	UserProfilePicture userProfilePicture `json:"userProfilePicture"`
}

type userProfilePicture struct {
	TypeName      string                `json:"__typename"`
	BaseUrl       string                `json:"baseUrl"`
	OnPressAction navigateToUserProfile `json:"onPressAction"`
}

type navigateToUserProfile struct {
	TypeName string `json:"__typename"`
	Url      string `json:"url"`
}

type localizedReview struct {
	TypeName           string      `json:"__typename"`
	Comments           string      `json:"comments"`
	CommentsLanguage   string      `json:"commentsLanguage"`
	Disclaimer         string      `json:"disclaimer"`
	NeedsTranslation   bool        `json:"needsTranslation"`
	Response           interface{} `json:"response"`
	ResponseDisclaimer interface{} `json:"responseDisclaimer"`
}

type showMoreButton struct {
	TypeName         string           `json:"__typename"`
	Title            string           `json:"title"`
	LoggingEventData loggingEventData `json:"loggingEventData"`
}

type loggingEventData struct {
	TypeName            string        `json:"__typename"`
	LoggingId           string        `json:"loggingId"`
	Experiments         []interface{} `json:"experiments"`
	EventData           interface{}   `json:"eventData"`
	EventDataSchemaName interface{}   `json:"eventDataSchemaName"`
	Section             interface{}   `json:"section"`
	Component           interface{}   `json:"component"`
}

type pricingWrapper1 struct {
	StructuredStayDisplayPrice pricingWrapper2 `json:"structuredStayDisplayPrice"`
}

type pricingWrapper2 struct {
	PrimaryLine     priceData    `json:"primaryLine"`
	SecondaryLine   priceData    `json:"secondaryLine"`
	ExplanationData priceDetails `json:"explanationData"`
}

type listingData struct {
	AvgRatingA11yLabel string                  `json:"avgRatingA11yLabel"`
	AvgRatingLocalized string                  `json:"avgRatingLocalized"`
	City               string                  `json:"city"`
	ContextualPictures []pricture              `json:"contextualPictures"`
	Coordinate         coordinate              `json:"coordinate"`
	FormattedBadges    []formatedbadgeWrapper1 `json:"formattedBadges"`
	Id                 int64                   `json:"id,string"`
	ListingObjType     string                  `json:"listingObjType"`
	LocalizedCityName  string                  `json:"localizedCityName"`
	Name               string                  `json:"name"`
	PdpUrlType         string                  `json:"pdpUrlType"`
	RoomTypeCategory   string                  `json:"roomTypeCategory"`
	TierId             int                     `json:"tierId"`
	Title              string                  `json:"title"`
	TitleLocale        string                  `json:"titleLocale"`
}
type pricture struct {
	Picture string `json:"picture"`
}

type coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type formatedbadgeWrapper1 struct {
	LoggingContext formatedbadgeWrapper2 `json:"loggingContext"`
}

type formatedbadgeWrapper2 struct {
	BadgeType string `json:"badgeType"`
}

type priceData struct {
	DisplayComponentType string `json:"displayComponentType"`
	AccessibilityLabel   string `json:"accessibilityLabel"`
	Price                string `json:"price"`
	OriginalPrice        string `json:"originalPrice"`
	DiscountedPrice      string `json:"discountedPrice"`
	Qualifier            string `json:"qualifier"`
	ShortQualifier       string `json:"shortQualifier"`
	ConcatQualifierLeft  bool   `json:"concatQualifierLeft"`
}

type priceDetails struct {
	PriceDetails []items `json:"priceDetails"`
}

type items struct {
	Items []itemData `json:"items"`
}

type itemData struct {
	DisplayComponentType string `json:"displayComponentType"`
	Description          string `json:"description"`
	PriceString          string `json:"priceString"`
}
