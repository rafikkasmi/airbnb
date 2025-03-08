package availability

type Data struct {
	RoomID           int64
	Badges           []string
	Name             string
	Title            string
	Type             string
	Kind             string
	Category         string
	Rating           Rating
	Coordinates      Coordinates
	Fee              Fee
	Price            PriceData
	LongStayDiscount Price
	Images           []Img
	Reviews          []Review
	Availability     AvailabilityData
}

type AvailabilityData struct {
	CalendarMonths []CalendarMonth
}

type CalendarMonth struct {
	Month int
	Year  int
	Days  []CalendarDay
}

type CalendarDay struct {
	Date                 string   `csv:"Date"`
	Available            bool     `csv:"Available"`
	MaxNights            int      `csv:"MaxNights"`
	MinNights            int      `csv:"MinNights"`
	AvailableForCheckin  bool     `csv:"AvailableForCheckin"`
	AvailableForCheckout bool     `csv:"AvailableForCheckout"`
	Price                DayPrice `csv:"Price"`
}

type DayPrice struct {
	Formatted string
	Amount    float32
	Currency  string
}

type Review struct {
	ID                 string
	Comments           string
	TranslatedComments string
	Language           string
	CreatedAt          string
	LocalizedDate      string
	Rating             int
	Reviewer           User
	Host               User
	Highlight          string
}

type User struct {
	ID            string
	FirstName     string
	FullName      string
	IsSuperhost   bool
	ProfilePicURL string
	ProfilePath   string
}
type Fee struct {
	Cleaning Price
	Airbn    Price
}
type PriceData struct {
	Total     Price
	Unit      UnitPrice
	BreakDown []PriceBreakDown
}
type PriceBreakDown struct {
	Description    string
	Amount         float32
	CurrencySymbol string
}
type Rating struct {
	Value       float32
	ReviewCount int
}
type Coordinates struct {
	Latitude float64
	Longitud float64
}

type Img struct {
	URL         string
	ContentType string
	Extension   string
	Content     []byte `json:"-"`
}
type UnitPrice struct {
	Qualifier string
	Discount  float32
	Price
}
type Price struct {
	Amount         float32
	CurrencySymbol string
}
